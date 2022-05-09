package runningTime

import "fmt"

func {{$.FunctionName}}(str string)(){
    {{if eq .Method true}}
    fmt.Println("method is true")
    {{end}}
    fmt.Println("{{.Name}} str",str)
}

