package gen

import (
	"regexp"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"gopkg.in/yaml.v3"
)

type SchemaWrapper struct {
	any
}

func (a SchemaWrapper) Key(key string) SchemaWrapper {
	return SchemaWrapper{a.any.(map[string]any)[key]}
}

func (a SchemaWrapper) Index(index int) SchemaWrapper {
	return SchemaWrapper{a.any.([]any)[index]}
}

func (a SchemaWrapper) Value() any {
	return a.any
}

func addFile(fd *desc.FileDescriptor, dst []*descriptorpb.FileDescriptorProto) []*descriptorpb.FileDescriptorProto {
	for _, dep := range fd.GetDependencies() {
		dst = addFile(dep, dst)
	}
	dst = append(dst, fd.AsFileDescriptorProto())
	return dst
}

func parseProto(t *testing.T, path string, opts ...PluginOption) *pluginpb.CodeGeneratorResponse {
	p := &Plugin{}
	for _, o := range opts {
		o(p)
	}

	t.Helper()
	parser := protoparse.Parser{
		ImportPaths: []string{
			".",
			// repository root
			"../../../../../../",
		},
		IncludeSourceCodeInfo: true,
		InferImportPaths:      false,
	}
	fds, err := parser.ParseFiles(path)
	if err != nil {
		return nil
	}

	var fdps []*descriptorpb.FileDescriptorProto
	for _, fd := range fds {
		fdps = addFile(fd, fdps)
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{path},
		ProtoFile:      fdps,
	}

	plugin, err := protogen.Options{}.New(req)
	if !assert.NoError(t, err) {
		return nil
	}

	err = p.Run(plugin)
	if err != nil {
		return nil
	}

	return plugin.Response()
}

func TestNoCRD(t *testing.T) {
	response := parseProto(t, "testdata/no_crd.proto")
	if response == nil {
		return
	}

	assert.Nil(t, response.Error)
	assert.Empty(t, response.File)
}

func TestTwoCRDs(t *testing.T) {
	response := parseProto(t, "testdata/two_crds.proto")
	if response == nil {
		return
	}

	// FIXME (torkve) in case of >1 crd we should report error,
	//                but we just skip such files for now
	//assert.NotEmpty(t, response.Error)
	assert.Empty(t, response.File)
}

func TestSimpleSpec(t *testing.T) {
	response := parseProto(t, "testdata/simple_spec.proto")
	if response == nil {
		return
	}
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)
	apiSpec := SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))
	versions := apiSpec.Key("spec").Key("versions")
	assert.Len(t, versions.Value(), 1)
	schema := versions.Index(0).Key("schema").Key("openAPIV3Schema").Key("properties")
	assert.Contains(t, schema.Value(), "spec")
	assert.Contains(t, schema.Value(), "status")

	spec := schema.Key("spec").Key("properties")
	assert.Equal(t, map[string]any{"type": "string"}, spec.Key("f1").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "int32"}, spec.Key("f2").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "uint32"}, spec.Key("f3").Value())
	assert.Equal(t, map[string]any{"x-kubernetes-int-or-string": true, "format": "int64"}, spec.Key("f4").Value())
	assert.Equal(t, map[string]any{"x-kubernetes-int-or-string": true, "format": "uint64"}, spec.Key("f5").Value())
	assert.Equal(t, map[string]any{"type": "string"}, spec.Key("f6").Key("properties").Key("f1").Value())

	nested := spec.Key("f6").Key("properties")
	assert.Equal(t, map[string]any{"type": "object", "x-kubernetes-preserve-unknown-fields": true}, nested.Key("f2").Value())
	assert.Equal(t, map[string]any{"type": "object", "nullable": true, "properties": map[string]any{}}, nested.Key("f3").Value())

	assert.Equal(t, map[string]any{"x-kubernetes-preserve-unknown-fields": true}, spec.Key("f7").Value())
}

func TestWrappers(t *testing.T) {
	response := parseProto(t, "testdata/wrappers_spec.proto")
	if response == nil {
		return
	}
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)

	apiSpec := SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))
	spec := apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties").Key("spec").Key("properties")
	assert.Equal(t, map[string]any{"type": "boolean", "nullable": true}, spec.Key("bool1").Value())
	assert.Equal(t, map[string]any{"type": "number", "format": "float", "nullable": true}, spec.Key("float1").Value())
	assert.Equal(t, map[string]any{"type": "number", "format": "double", "nullable": true}, spec.Key("double1").Value())
	assert.Equal(t, map[string]any{"type": "string", "nullable": true}, spec.Key("string1").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "int32", "nullable": true}, spec.Key("int321").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "uint32", "nullable": true}, spec.Key("uint321").Value())
	assert.Equal(t, map[string]any{"format": "int64", "nullable": true, "x-kubernetes-int-or-string": true}, spec.Key("int641").Value())
	assert.Equal(t, map[string]any{"format": "uint64", "nullable": true, "x-kubernetes-int-or-string": true}, spec.Key("uint641").Value())
	assert.Equal(t, map[string]any{"type": "array", "items": map[string]any{"x-kubernetes-preserve-unknown-fields": true}}, spec.Key("listvalue").Value())
}

func TestClientSchema(t *testing.T) {
	response := parseProto(t, "testdata/client_spec.proto")
	if response == nil {
		return
	}
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)

	apiSpec := SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))

	schema := apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties")
	spec := schema.Key("spec").Key("properties")
	// test only common fields
	assert.Equal(t, map[string]any{
		"f1": map[string]any{"type": "string"},
		"f3": map[string]any{"x-kubernetes-preserve-unknown-fields": true},
	}, spec.Value())
	assert.Nil(t, schema.Key("f2").Value())

	response = parseProto(t, "testdata/client_spec.proto", WithClientSchema(true))
	if response == nil {
		return
	}
	apiSpecData = response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)
	apiSpec = SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))
	schema = apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties")
	spec = schema.Key("spec").Key("properties")
	assert.Equal(t, map[string]any{
		"f1": map[string]any{"type": "string"},
		"f2": map[string]any{"type": "string"},
		"f3": map[string]any{"x-kubernetes-preserve-unknown-fields": true, "deprecated": true},
	}, spec.Value())
	assert.Equal(t, map[string]any{"type": "string"}, schema.Key("f2").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "int32"}, schema.Key("f3").Value())
	assert.Equal(t, map[string]any{"format": "int64", "x-kubernetes-int-or-string": true}, schema.Key("f4").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "uint32"}, schema.Key("f5").Value())
	assert.Equal(t, map[string]any{"type": "string"}, schema.Key("f6").Key("properties").Key("f1").Value())
	assert.Equal(t, map[string]any{
		"format":                     "enum",
		"enum":                       []any{0, "F1", 1, "F2"},
		"x-kubernetes-int-or-string": true,
	}, schema.Key("f7").Value())
	assert.Equal(t, map[string]any{"type": "string", "nullable": true}, schema.Key("f8").Value())
	assert.Equal(t, map[string]any{"type": "string", "nullable": true}, schema.Key("f9").Value())
	assert.Equal(t, map[string]any{"type": "array", "items": map[string]interface{}{"type": "string"}}, schema.Key("f10").Value())
	assert.Equal(t, map[string]any{"type": "object", "additionalProperties": map[string]interface{}{"type": "string"}}, schema.Key("f11").Value())
}

func TestSchemaless(t *testing.T) {
	response := parseProto(t, "testdata/client_spec.proto", WithSchemalessCrd(true))
	if response == nil {
		return
	}
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)

	apiSpec := SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))

	schema := apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties")
	assert.Equal(t, map[string]any{
		"spec": map[string]any{
			"x-kubernetes-preserve-unknown-fields": true,
			"type":                                 "object",
		},
		"status": map[string]any{
			"x-kubernetes-preserve-unknown-fields": true,
			"type":                                 "object",
		},
	}, schema.Value())

	assert.Nil(t, parseProto(t, "testdata/client_spec.proto", WithSchemalessCrd(true), WithClientSchema(true)))
}

func TestPatchExternals(t *testing.T) {
	response := parseProto(t, "testdata/patch_externals.proto", WithClientSchema(true))
	if response == nil {
		return
	}
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)

	apiSpec := SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))

	schema := apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties")
	container := schema.Key("spec").Key("properties").Key("container").Key("properties")

	assert.Equal(t,
		map[string]any{
			"type":                  "object",
			"nullable":              true,
			"additionalProperties":  false,
			"properties":            map[string]any{},
			patchMergeKeyField:      "key",
			patchMergeStrategyField: "merge",
		}, container.Key("nested").Value())
	assert.Equal(t,
		map[string]any{
			"type":                  "object",
			"nullable":              true,
			"additionalProperties":  false,
			"properties":            map[string]any{},
			patchMergeKeyField:      "key",
			patchMergeStrategyField: "merge",
		}, container.Key("nested_with_annotation").Value())
	assert.Equal(t,
		map[string]any{
			"type":                  "object",
			"nullable":              true,
			"additionalProperties":  false,
			"properties":            map[string]any{},
			patchMergeStrategyField: "drop",
		}, container.Key("nested_with_annotation2").Value())
}

func TestWellKnown(t *testing.T) {
	response := parseProto(t, "testdata/well_known.proto")
	if response == nil {
		return
	}
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)

	apiSpec := SchemaWrapper{any: map[string]any{}}
	assert.NoError(t, yaml.Unmarshal([]byte(apiSpecData), &apiSpec.any))
	spec := apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties").Key("spec").Key("properties")
	assert.Equal(t, map[string]any{"type": "object", "nullable": true, "properties": map[string]any{}}, spec.Key("empty_value").Value())
}

func TestYamlSpecialNames(t *testing.T) {
	response := parseProto(t, "testdata/yaml_conversions.proto")
	assert.Empty(t, response.Error)
	assert.Len(t, response.File, 1)

	apiSpecData := response.File[0].GetContent()
	assert.NotEmpty(t, apiSpecData)

	// NOTE (torkve) we cannot use yaml.v3 to test unmarshalling (it does it correct), and we do not want to depend on yaml.v2, so let's just regexp it.
	yamlUnquoted := regexp.MustCompile(`\n\s+y:\n\s+type: string\n`)
	yamlQuoted := regexp.MustCompile(`\n\s+"y":\n\s+type: string\n`)
	assert.Nil(t, yamlUnquoted.FindStringIndex(apiSpecData), apiSpecData)
	assert.NotNil(t, yamlQuoted.FindStringIndex(apiSpecData), apiSpecData)
}
