package main

import (
	"bytes"
	"html/template" // TODO: 源码
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

/* 实现插件 protoc-gen-go-netrpc 为标准库的 rpc 框架生成代码 */

// > go mod init protoc-gen-go-netrpc, 定义包名(插件名称)
// > ...
// > go intall, 安装该插件到 $GOROOT/bin 中
// > protoc --go-netrpc_out=plugins=netrpc:. hello.proto
// 生成的 go 文件为 ./main/hello.pb.go
// --go-netrpc_out 参数告知编译器加载名为 protoc-gen-go-netrpc 的插件,
// 插件中的 plugins=netrpc 指示启用内部唯一的名为 netrpc 的 netrpcPlugin 插件,
// 在新生成的 hello.pb.go 中将包含新增加的注释代码

// 使用该插件需要先通过 generator.RegisterPlugin() 函数注册插件
func init() {
	generator.RegisterPlugin(new(netrpcPlugin))
}

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
	n.P(`import "net/rpc"`)
}

func (n *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
	// 基于 buildServiceSpec() 方法构造的服务的元信息生成服务的代码
	spec := n.buildServiceSpec(svc)

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmplService))
	err := t.Execute(&buf, spec)
	if err != nil {
		log.Fatal(err)
	}

	n.P(buf.String())
}

// 要在自定义的 genServiceCode() 方法中为每个服务生成相关的代码,
// 每个服务最重要的是服务的名字, 每个服务有一组方法, 而对于服务定义
// 的方法, 最重要的是方法的名字, 还有输入参数和输出参数类型的名字
type ServiceSpec struct {
	ServiceName string
	MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}

// 新建 buildServiceSpec() 方法用来解析每个服务的 ServiceSpec 元信息
func (n *netrpcPlugin) buildServiceSpec(
	svc *descriptor.ServiceDescriptorProto,
	// 此参数完整描述了一个服务的所有信息
) *ServiceSpec {
	spec := &ServiceSpec{
		ServiceName: generator.CamelCase(svc.GetName()),
		// svc.GetName() 获取 protobuf 文件中定义的服务的名字
		// protobuf 文件中的名字转为 go 语言的名字后, 需要通过 CamelCase() 函数
		// 进行一次转换
	}

	for _, m := range svc.Method {
		spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
			MethodName: generator.CamelCase(m.GetName()),
			// 对输入参数和输出参数的解析:
			// 先通过 GetInputType() 获取输入参数的类型, 再通过 ObjectNamed()
			// 获取类型对应的类对象信息, 最后获取类对象的名字
			InputTypeName:  n.TypeName(n.ObjectNamed(m.GetInputType())),
			OutputTypeName: n.TypeName(n.ObjectNamed(m.GetOutputType())),
		})
	}

	return spec
}

// TODO:
// 当 protobuf 插件定制工作完成后, 每次 ./hello.proto 文件中的 rpc 服务的变化
// 都可以自动生成代码, 也可以通过更新插件的模板, 调整或增加生成代码的内容
// 服务的模板(基于go语言的模板生成代码)
const tmplService = `
	{{$root := .}}
	
	type {{.ServiceName}}Interface interface {
		{{- range $_, $m := .MethodList}}
		{{$m.MethodName}}(*{{$m.InputTypeName}}, *{{$m.OutputTypeName}}) error
		{{- end}}
	}


	func Register{{.ServiceName}} (
		srv *rpc.Server, x {{.ServiceName}}Interface,
	) error {
		if err := srv.RegisterName("{{.ServiceName}}", x); err != nil {
			return err
		}
		return nil
	}

	type {{.ServiceName}}Client struct {
		*rpc.Client
	}

	var _ {{.ServiceName}}Interface = (*{{.ServiceName}}Client)(nil)

	func Dial{{.ServiceName}} (network, address string) (
		*{{.ServiceName}}Client, error,
	) {
		c, err := rpc.Dial(network, address)
		if err != nil {
			return nil, err
		}

		return &{{.ServiceName}}Client{Client: c}, nil
	}

	{{range $_, $m := .MethodList}}
	func (n *{{$root.ServiceName}}Client) {{$m.MethodName}} (
		in *{{$m.InputTypeName}}, out *{{$m.OutputTypeName}},
	) error {
		return n.Client.Call("{{$root.ServiceName}}.{{$m.MethodName}}", in, out)
	}

	{{end}}
`
