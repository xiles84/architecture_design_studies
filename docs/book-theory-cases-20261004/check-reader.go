// Acceptance audit of the assembled reader path and rendered teaching citations.
package main

import (
 "encoding/json"
 "fmt"
 "os"
 "path/filepath"
 "strings"
)

func must(err error) { if err != nil { panic(err) } }
func read(path string) string { b,e:=os.ReadFile(path);must(e);return string(b) }
func require(ok bool, message string) { if !ok { panic(message) } }
func main() {
 root:=os.Args[1]
 text:=read(os.Args[2])
 mainSource:=read(filepath.Join(root,"book/main.typ"))
 require(!strings.Contains(mainSource,"include \"chapters/"),"legacy chapter still in reading path")
 require(!strings.Contains(mainSource,"include \"concepts/"),"legacy concept still in reading path")
 for _,bad:=range []string{"not converted yet","awaiting conversion","else [","unverified development build"} {
  require(!strings.Contains(strings.ToLower(text),bad),"rendered forbidden token: "+bad)
 }
 theoryFiles,e:=filepath.Glob(filepath.Join(root,"book/theory/*.typ"));must(e)
 count:=0;for _,f:=range theoryFiles { count+=strings.Count(read(f),"#strategy(") }
 require(len(theoryFiles)==6 && count==20,"wrong theory group/strategy coverage")
 cases,e:=filepath.Glob(filepath.Join(root,"book/cases/*.typ"));must(e)
 require(len(cases)==5,"wrong applied-case coverage")
 parts:=strings.Split(text,"Appendix · Evidence and reproduction")
 require(len(parts)>=3,"appendix heading/contents not found")
 // Last split is the actual appendix, earlier occurrence is its contents entry.
 teaching:=strings.Join(parts[:len(parts)-1],"Appendix · Evidence and reproduction")
 var registry struct { Claims []struct { ClaimID string `json:"claim_id"` } `json:"claims"` }
 must(json.Unmarshal([]byte(read(filepath.Join(root,"book/evidence/v6/claims.json"))),&registry))
 for _,c:=range registry.Claims { require(strings.Contains(teaching,c.ClaimID),"claim missing from teaching: "+c.ClaimID) }
 fmt.Printf("PASS: %d topic groups, %d strategies, %d cases; %d/%d claims in rendered teaching; legacy/literal-template checks clear.\n",len(theoryFiles),count,len(cases),len(registry.Claims),len(registry.Claims))
}
