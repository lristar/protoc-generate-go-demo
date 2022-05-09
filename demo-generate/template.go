package demo_generate

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

var tpl string

type service struct {
	Name     string //Demo
	FullName string //hello.Demo
	FilePath string //api/hello/hello.proto

	Method    []*method
	MethodSet map[string]*method
}

func (s *service) execute() string {
	if s.MethodSet == nil {
		s.MethodSet = map[string]*method{}
		for _, m := range s.Method {
			m1 := m
			s.MethodSet[m.Name] = m1
		}
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return buf.String()

}

func (s *service) InterfaceName() string {
	return s.Name + "_" + "HTTPServer"
}

type method struct {
	Name    string // SayHello
	Num     int    // 一个 rpc 方法可以对应多个 http 请求
	Request string // SayHelloReq
	Reply   string // SayHelloResp
	// http_rule
	Path         string // 路由
	Method       string // HTTP Method
	Body         string
	ResponseBody string
}

func (m *method) HandlerName() string {
	return fmt.Sprintf("%s_%s", m.Name, MathecEnglish(m.Num))
}

func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}
	return false
}

// 转换参数路由 {type} --> :type
func (m *method) EPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, path := range paths {
		if len(path) > 0 && (path[0] == '{' && path[len(path)-1] == '}' || path[0] == ':') {
			paths[i] = ":" + path[1:len(path)-1]
		}
	}
	m.Path = strings.Join(paths, "/")
}

//数字转英文
func MathecEnglish(i int) string {
	switch i {
	case 1:
		return "ONE"
	case 2:
		return "TWO"
	case 3:
		return "THREE"
	case 4:
		return "FOUR"
	case 5:
		return "FIVE"
	case 6:
		return "SIX"
	case 7:
		return "SEVEN"
	case 8:
		return "EIGHT"
	case 9:
		return "NINE"
	case 10:
		return "TEN"
	default:
		return "ALL"
	}
}
