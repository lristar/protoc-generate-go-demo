package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

// 将文件模板写入到定义的变量中

//go:embed demo.go.tpl
var tpl string

//go:embed demo2.go.tpl
var tpl2 string

func NewTemplate() {
	u := struct {
		Name   string
		Gender string
		Age    int
	}{
		Name:   "zhangsan",
		Gender: "男",
		Age:    18,
	}
	tplStr := template.Must(template.New("demo").Parse(tpl))
	tplStr = template.Must(tplStr.New("demo2").Parse(tpl2))
	var buff bytes.Buffer
	// u 对应的数据集合  ”demo2“是模板对应的namespace  buff是输出的字节缓冲流变量
	if err := tplStr.ExecuteTemplate(&buff,"demo2", u); err != nil {
		panic(err)
	}
	fmt.Println( buff.String())

	// 清掉内容
	buff.Reset()

	if err := tplStr.ExecuteTemplate(&buff,"demo", u); err != nil {
		panic(err)
	}
	fmt.Println( buff.String())
}

func main() {
	NewTemplate()
}
