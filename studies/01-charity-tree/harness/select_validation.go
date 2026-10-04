package main

import (
	"context"
	"fmt"

	"charitytree/selectbench"
)

type selectValidation struct {
	Case        string             `json:"case"`
	Fingerprint string             `json:"fixture_fingerprint"`
	Cells       []selectbench.Cell `json:"cells"`
}

func validateSelectFixtures(ctx context.Context, q selectbench.Querier, catalog *selectbench.Catalog) ([]selectValidation, []selectbench.ControlResult, error) {
	baseArgs := selectbench.Args{Pred: selectbench.AmountGreaterThan(100), PredSQL: "amount > 100", N: 1, Limit: 2, ValidateOnly: true}
	teaching := selectbench.Teaching()
	tie := selectbench.Page(teaching, 2, 0)[1]
	nullFixture := selectbench.Teaching()
	nullFixture.Orders = append(nullFixture.Orders, selectbench.Order{ID: 17, CustomerID: &nullFixture.Customers[1].ID})
	nullCursor := selectbench.PageRow{OrderID: 13}
	cases := []struct {
		name       string
		fx         selectbench.Fixture
		args       selectbench.Args
		statements []string
	}{
		{"teaching", teaching, baseArgs, nil},
		{"empty", selectbench.Fixture{}, baseArgs, nil},
		{"all-miss", teaching, baseArgs, nil},
		{"keyset-tie", teaching, baseArgs, []string{"s17_page_keyset"}},
		{"keyset-null", nullFixture, baseArgs, []string{"s17_page_keyset"}},
	}
	cases[2].args.Pred, cases[2].args.PredSQL = selectbench.AmountGreaterThan(1000), "amount > 1000"
	cases[3].args.Last, cases[4].args.Last = &tie, &nullCursor
	var validation []selectValidation
	var controls []selectbench.ControlResult
	var gateErr error
	for _, item := range cases {
		if err := selectbench.Prepare(ctx, q, item.fx); err != nil {
			return validation, controls, fmt.Errorf("validation %s preparation: %w", item.name, err)
		}
		cells := selectbench.RunCells(ctx, q, catalog, item.fx, item.args, item.statements, "fixture-"+item.name, "validation")
		validation = append(validation, selectValidation{item.name, selectbench.Fingerprint(item.fx), cells})
		for _, cell := range cells {
			if cell.Status != selectbench.StatusValid {
				gateErr = fmt.Errorf("live fixture %s: %s is %s: %s", item.name, cell.Statement, cell.Status, cell.Note)
			}
		}
		if item.name == "teaching" {
			controls = selectbench.RunControls(ctx, q, catalog, item.fx, item.args)
			if len(controls) != 7 {
				gateErr = fmt.Errorf("expected seven live negative controls, got %d", len(controls))
			}
			for _, control := range controls {
				if !control.Detected {
					gateErr = fmt.Errorf("control %s failed detection: %s", control.Statement, control.Error)
				}
			}
		}
	}
	return validation, controls, gateErr
}
