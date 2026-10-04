// Audit the added reader path and the owner-directed capability route.
package main

import (
 "encoding/json"
 "fmt"
 "os"
 "path/filepath"
 "strings"
)

func read(path string) string { b,e:=os.ReadFile(path); if e!=nil { panic(e) }; return string(b) }
func require(ok bool, message string) { if !ok { panic(message) } }
func main() {
 root:=os.Args[1]
 textFile:="docs/book-join-theory-20261004/rendered.txt"
 if len(os.Args)>2 {textFile=os.Args[2]}
 text:=read(filepath.Join(root,textFile))
 theory,e:=filepath.Glob(filepath.Join(root,"book/theory/*.typ")); if e!=nil {panic(e)}
 count:=0; for _,f:=range theory {count+=strings.Count(read(f),"#strategy(")}
 require(len(theory)==6 && count==22,"wrong theory coverage")
 cases,e:=filepath.Glob(filepath.Join(root,"book/cases/*.typ")); if e!=nil {panic(e)}
 require(len(cases)==5,"wrong case coverage")
 for _,bad:=range []string{"not converted yet","awaiting conversion","else [","unverified development build"} {
  require(!strings.Contains(strings.ToLower(text),bad),"rendered forbidden token: "+bad)
 }
 for _,label:=range []string{"1. Explicit join", "2. Filter inside ON", "3. Filtered derived table", "4. CTE naming", "5. CTE naming", "6. Comma join", "Preserve every customer", "Keep only qualifying", "CTE folding and materialization", "Planned variation: join donations"} {
  require(strings.Contains(text,label),"missing rendered example: "+label)
 }
 parts:=strings.Split(text,"Appendix · Evidence and reproduction")
 require(len(parts)>=3,"missing appendix navigation")
 teaching:=strings.Join(parts[:len(parts)-1],"Appendix · Evidence and reproduction")
 var registry struct {Claims []struct {ID string `json:"claim_id"`} `json:"claims"`}
 if e:=json.Unmarshal([]byte(read(filepath.Join(root,"book/evidence/v6/claims.json"))),&registry);e!=nil {panic(e)}
 for _,c:=range registry.Claims {require(strings.Contains(teaching,c.ID),"missing teaching claim: "+c.ID)}
 var task struct {Min string `json:"minimum_capability"`; Dependencies []string `json:"dependencies"`; Supersedes []string `json:"supersedes"`}
 if e:=json.Unmarshal([]byte(read(filepath.Join(root,"docs/ai-work/tasks/2026/10/task-20261004T142623Z-select-book-high/task.json"))),&task);e!=nil {panic(e)}
 require(task.Min=="HIGH","publication permits LOW")
 require(strings.Contains(strings.Join(task.Dependencies,","),"select-synthesis"),"publication bypasses synthesis")
 require(len(task.Supersedes)==1 && task.Supersedes[0]=="task-20261002T131729Z-select-book","missing predecessor link")
 for _,taskID:=range []string{"select-join-cte","select-synthesis","select-book"} {
  amendments,e:=filepath.Glob(filepath.Join(root,"docs/ai-work/tasks/2026/10/task-20261002T131729Z-"+taskID+"/amendments/20261004*.md"));if e!=nil {panic(e)}
  require(len(amendments)==1,"missing task amendment: "+taskID)
 }
 fmt.Printf("PASS: %d theory groups, %d strategies, %d cases; six rendered SQL forms and both LEFT JOIN contracts; %d/%d teaching claims; LOW study -> HIGH synthesis -> HIGH publication.\n",len(theory),count,len(cases),len(registry.Claims),len(registry.Claims))
}
