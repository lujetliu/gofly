package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

/* 实现插件 proto-gen-go-netrpc 为标准库的 rpc 框架生成代码 */

// 因为 go 语言的包只能静态导入, 所有无法向已经安装的 protoc-gen-go 添加新
// 编写的插件, 所以这里克隆 protoc-gen-go 对应的 main() 函数

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	// send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}

type netrpcPlugin struct {
	*generator.Generator
}

func (n *netrpcPlugin) Name() string                { return "netrpc" }
func (n *netrpcPlugin) Init(g *generator.Generator) { n.Generator = g }

func (n *netrpcPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		n.genImportCode(file)
	}
}

func (n *netrpcPlugin) Generate(file *generator.FileDescriptor) {
	for _, svc := range file.Service {
		n.genServiceCode(svc)
	}
}

func (n *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
	n.P("//TODO:import code")
}

func (n *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
	n.P("//TODO:service code, Name =" + svc.GetName())
}
