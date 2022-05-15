package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"

	"github.com/fananchong/test_protobuf_options"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	contextPkgPath = "context"
)

const (
	errPkg = protogen.GoImportPath("errors")
)

func init() {
	generator.RegisterPlugin(new(gogogo))
}

// test is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for test support.
type gogogo struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "test".
func (g *gogogo) Name() string {
	return "gogogo"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	contextPkg string
)

// Init initializes the plugin.
func (g *gogogo) Init(gen *generator.Generator) {
	g.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *gogogo) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *gogogo) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *gogogo) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *gogogo) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	contextPkg = string(g.gen.AddImport(contextPkgPath))

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}
}

// GenerateImports generates the import declaration for this file.
func (g *gogogo) GenerateImports(file *generator.FileDescriptor) {
	// import定义
	g.P("//", errPkg.Ident(""), errPkg.Ident(""))
}

func unexport(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// generateService generates all the code for the named service.
func (g *gogogo) generateService(file *generator.FileDescriptor, srv *pb.ServiceDescriptorProto, index int) {
	origServName := srv.GetName()
	servName := generator.CamelCase(origServName)
	servAlias := servName + "Service"

	// strip suffix
	if strings.HasSuffix(servAlias, "ServiceService") {
		servAlias = strings.TrimSuffix(servAlias, "Service")
	}

	g.P()
	g.P("// For example")
	g.P()
	sd := &service{
		Name: srv.GetName(),
	}

	//// Client interface.
	//g.P("type ", servAlias, " interface {")
	for _, method := range srv.Method {
		//g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
		//g.P(g.generateClientSignature(servName, method))
		sd.Methods = append(sd.Methods, GenMethod(method)...)
	}
	g.P(sd.execute())
}

// 生成method方法
func GenMethod(m *pb.MethodDescriptorProto) (methods []*method) {
	methods = make([]*method, 0)
	method := method{
		Name:    *m.Name,
		Num:     0,
		Request: *m.InputType,
		Reply:   *m.OutputType,
	}
	methods = append(methods, &method)
	return
}

// generateClientSignature returns the client-side signature for a method.
func (g *gogogo) generateClientSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)

	isBroadcast := false
	if v, err := proto.GetExtension(method.GetOptions(), mypack.E_Broadcast); err == nil {
		isBroadcast = *(v.(*bool))
	}
	if isBroadcast {
		return fmt.Sprintf("Broadcast%s(ctx %s.Context) error", methName, contextPkg)
	}
	return fmt.Sprintf("%s(ctx %s.Context) error", methName, contextPkg)
}

// AddPluginToParams Simplify the protoc call statement by adding 'plugins=test' directly to the command line arguments.
func AddPluginToParams(p string) string {
	params := p
	if strings.Contains(params, "plugins=") {
		params = strings.Replace(params, "plugins=", "plugins=gogogo+", -1)
	} else {
		if len(params) > 0 {
			params += ","
		}
		params += "plugins=gogogo"
	}
	return params
}
