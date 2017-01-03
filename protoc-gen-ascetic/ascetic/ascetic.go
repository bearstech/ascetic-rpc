package ascetic

import (
	"fmt"
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"log"
	"os"
)

var (
	l          *log.Logger
	asceticPkg string
)

func init() {
	generator.RegisterPlugin(new(ascetic))
	l = log.New(os.Stderr, "", 0)
	l.Println("init")
}

type ascetic struct {
	gen *generator.Generator
}

func (a *ascetic) Name() string {
	return "ascetic"
}

func (a *ascetic) Init(g *generator.Generator) {
	l.Println("Init")
	a.gen = g
	asceticPkg = generator.RegisterUniquePackageName("ascetic", nil)
}

func (a *ascetic) Generate(file *generator.FileDescriptor) {
	l.Println("Generate")
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	for i, service := range file.FileDescriptorProto.Service {
		a.generateService(file, service, i)
	}
}

func (a *ascetic) GenerateImports(file *generator.FileDescriptor) {
	l.Println("Generate Import")
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
}

func (a *ascetic) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) // 6 means service.
	l.Println("path", path)

	origServName := service.GetName()
	fullServName := origServName
	if pkg := file.GetPackage(); pkg != "" {
		fullServName = pkg + "." + fullServName
	}
	servName := generator.CamelCase(origServName)

	a.P()
	a.P("// Client API for ", servName, " service")
	a.P()

	a.P("// Server API for ", servName, " service")
	a.P()
	a.P("type ", origServName, "Server interface {")
	for i, service := range file.FileDescriptorProto.Service {
		a.generateService(file, service, i)
	}

	a.P("}")

}

// P forwards to g.gen.P.
func (a *ascetic) P(args ...interface{}) { a.gen.P(args...) }
