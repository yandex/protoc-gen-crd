package gen

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// parseProto runs the plugin over a proto file and requires it to succeed.
func parseProto(t *testing.T, path string, opts ...PluginOption) *pluginpb.CodeGeneratorResponse {
	t.Helper()
	response, err := runProto(t, path, nil, opts...)
	require.NoError(t, err)
	require.NotNil(t, response)
	return response
}

// runProtoSource runs the plugin over an inline proto source, which may import the repository protos.
func runProtoSource(t *testing.T, source string, opts ...PluginOption) (*pluginpb.CodeGeneratorResponse, error) {
	t.Helper()
	const path = "inline.proto"
	return runProto(t, path, map[string]string{path: source}, opts...)
}

// runProto runs the plugin over a proto file, reading the files in sources instead of the disk.
func runProto(t *testing.T, path string, sources map[string]string, opts ...PluginOption) (*pluginpb.CodeGeneratorResponse, error) {
	t.Helper()
	p := &Plugin{}
	for _, o := range opts {
		o(p)
	}

	parser := protoparse.Parser{
		ImportPaths: []string{
			".",
			// repository root
			"../../../../../../",
		},
		IncludeSourceCodeInfo: true,
		InferImportPaths:      false,
		Accessor: func(filename string) (io.ReadCloser, error) {
			if source, ok := sources[filename]; ok {
				return io.NopCloser(strings.NewReader(source)), nil
			}
			return os.Open(filename)
		},
	}
	fds, err := parser.ParseFiles(path)
	require.NoError(t, err)

	var fdps []*descriptorpb.FileDescriptorProto
	for _, fd := range fds {
		fdps = addFile(fd, fdps)
	}

	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{path},
		ProtoFile:      fdps,
	}

	plugin, err := protogen.Options{}.New(req)
	require.NoError(t, err)

	if err := p.Run(plugin); err != nil {
		return nil, err
	}
	return plugin.Response(), nil
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

func TestStrictsSchema(t *testing.T) {
	response := parseProto(t, "testdata/simple_spec.proto", WithScrictSchema(true))
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
	assert.Equal(t, false, schema.Key("spec").Key("additionalProperties").Value())
	assert.Equal(t, map[string]any{"type": "string"}, spec.Key("f1").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "int32"}, spec.Key("f2").Value())
	assert.Equal(t, map[string]any{"type": "integer", "format": "uint32"}, spec.Key("f3").Value())
	assert.Equal(t, map[string]any{"x-kubernetes-int-or-string": true, "format": "int64"}, spec.Key("f4").Value())
	assert.Equal(t, map[string]any{"x-kubernetes-int-or-string": true, "format": "uint64"}, spec.Key("f5").Value())
	assert.Equal(t, map[string]any{"type": "string"}, spec.Key("f6").Key("properties").Key("f1").Value())

	nested := spec.Key("f6").Key("properties")
	assert.Equal(t, map[string]any{"type": "object", "x-kubernetes-preserve-unknown-fields": true}, nested.Key("f2").Value())
	assert.Equal(t, map[string]any{"additionalProperties": false, "type": "object", "nullable": true, "properties": map[string]any{}}, nested.Key("f3").Value())

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

	_, err := runProto(t, "testdata/client_spec.proto", nil, WithSchemalessCrd(true), WithClientSchema(true))
	assert.Error(t, err)
}

// crdProperties runs the plugin over a proto and returns the properties of the CRD root schema.
func crdProperties(t *testing.T, path string, opts ...PluginOption) SchemaWrapper {
	t.Helper()
	response := parseProto(t, path, opts...)
	require.Len(t, response.File, 1)

	apiSpec := SchemaWrapper{any: map[string]any{}}
	require.NoError(t, yaml.Unmarshal([]byte(response.File[0].GetContent()), &apiSpec.any))
	return apiSpec.Key("spec").Key("versions").Index(0).Key("schema").Key("openAPIV3Schema").Key("properties")
}

// extensionsOf returns the x-kubernetes map and list type markers of a schema.
func extensionsOf(schema SchemaWrapper) map[string]any {
	markers := map[string]any{}
	for key, value := range schema.Value().(map[string]any) {
		if key == mapTypeField || key == listTypeField || key == listMapKeysField {
			markers[key] = value
		}
	}
	return markers
}

func TestTopology(t *testing.T) {
	for name, opts := range map[string][]PluginOption{
		"default": nil,
		"strict":  {WithScrictSchema(true)},
		"client":  {WithClientSchema(true), WithGeneratingMergeKeys(true)},
	} {
		t.Run(name, func(t *testing.T) {
			properties := crdProperties(t, "testdata/topology.proto", opts...)
			spec := properties.Key("spec").Key("properties")

			for field, want := range map[string]map[string]any{
				"atomic_message":          {mapTypeField: "atomic"},
				"granular_message":        {mapTypeField: "granular"},
				"atomic_map":              {mapTypeField: "atomic"},
				"atomic_list":             {listTypeField: "atomic"},
				"string_set":              {listTypeField: "set"},
				"item_map":                {listTypeField: "map", listMapKeysField: []any{"name"}},
				"atomic_struct":           {mapTypeField: "atomic"},
				"by_type":                 {mapTypeField: "atomic"},
				"by_type_over_annotation": {mapTypeField: "atomic"},
				"by_path":                 {mapTypeField: "atomic"},
				"reset_by_path":           {},
				"plain":                   {},
				"int_set":                 {listTypeField: "set"},
				"holders":                 {mapTypeField: "atomic"},
			} {
				assert.Equal(t, want, extensionsOf(spec.Key(field)), field)
			}

			holderValue := spec.Key("holders").Key("additionalProperties")
			assert.Empty(t, extensionsOf(holderValue))
			assert.Equal(t, map[string]any{mapTypeField: "atomic"}, extensionsOf(holderValue.Key("properties").Key("inner")))

			items := spec.Key("item_map").Key("items").Value().(map[string]any)
			assert.Equal(t, []any{"name"}, items["required"])
			assert.NotContains(t, items, "nullable")

			assert.Equal(t, map[string]any{mapTypeField: "granular"}, extensionsOf(properties.Key("spec")))
			assert.Equal(t, map[string]any{mapTypeField: "atomic"}, extensionsOf(properties.Key("status")))
		})
	}
}

func TestTopologyKeepsPatchMergeKey(t *testing.T) {
	properties := crdProperties(t, "testdata/topology.proto", WithClientSchema(true), WithGeneratingMergeKeys(true))
	itemMap := properties.Key("spec").Key("properties").Key("item_map").Value().(map[string]any)
	assert.Equal(t, "name", itemMap[patchMergeKeyField])
	assert.Equal(t, "merge", itemMap[patchMergeStrategyField])
}

func TestTopologySchemaless(t *testing.T) {
	properties := crdProperties(t, "testdata/topology.proto", WithSchemalessCrd(true))
	assert.Equal(t, map[string]any{
		"spec": map[string]any{
			"x-kubernetes-preserve-unknown-fields": true,
			"x-kubernetes-map-type":                "granular",
			"type":                                 "object",
			"description":                          "Desired state.",
		},
		"status": map[string]any{
			"x-kubernetes-preserve-unknown-fields": true,
			"x-kubernetes-map-type":                "atomic",
			"type":                                 "object",
			"description":                          "Observed state.",
		},
	}, properties.Value())
}

func TestTopologyErrors(t *testing.T) {
	const header = `
syntax = "proto3";
package testdata;
option go_package = "github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/pkg/gen/testdata";
import "google/protobuf/struct.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/wrappers.proto";
import "library/go/k8s/protoc_gen_crd/proto/crd.proto";
message Item {
    string name = 1;
    optional string optional_name = 2;
    Item nested = 3;
}
message Status {}
`
	crd := func(specFields string, crdOptions string) string {
		return header + `
message Spec {
` + specFields + `
}
message Kind {
    option (protoc_gen_crd.k8s_crd) = {
        api_group: "group", kind: "Kind", plural: "kinds", singular: "kind"` + crdOptions + `
    };
    Spec spec = 1;
    Status status = 2;
}
`
	}

	for name, tc := range map[string]struct {
		source  string
		wantErr string
	}{
		"set on scalar": {
			crd(`string f = 1 [(protoc_gen_crd.k8s_topology) = KT_SET];`, ""),
			`field testdata.Spec.f: k8s_topology KT_SET: requires a list, the field renders as "string"`,
		},
		"granular on list": {
			crd(`repeated string f = 1 [(protoc_gen_crd.k8s_topology) = KT_GRANULAR];`, ""),
			"field testdata.Spec.f: k8s_topology KT_GRANULAR: requires an object or a map, the field renders as a list",
		},
		"set on object": {
			crd(`Item f = 1 [(protoc_gen_crd.k8s_topology) = KT_SET];`, ""),
			"field testdata.Spec.f: k8s_topology KT_SET: requires a list, the field renders as an object",
		},
		"set of objects": {
			crd(`repeated Item f = 1 [(protoc_gen_crd.k8s_topology) = KT_SET];`, ""),
			"field testdata.Spec.f: k8s_topology KT_SET: requires scalar list items, the items render as an object",
		},
		"set of list values": {
			crd(`google.protobuf.ListValue f = 1 [(protoc_gen_crd.k8s_topology) = KT_SET];`, ""),
			"field testdata.Spec.f: k8s_topology KT_SET: requires scalar list items, the items render as an untyped value",
		},
		"set of json values": {
			crd(`repeated google.protobuf.Value f = 1 [(protoc_gen_crd.k8s_topology) = KT_SET];`, ""),
			"field testdata.Spec.f: k8s_topology KT_SET: requires scalar list items, the items render as an untyped value",
		},
		"set of nullable wrappers": {
			crd(`repeated google.protobuf.StringValue f = 1 [(protoc_gen_crd.k8s_topology) = KT_SET];`, ""),
			"field testdata.Spec.f: k8s_topology KT_SET: requires non-nullable list items",
		},
		"map without merge key": {
			crd(`repeated Item f = 1 [(protoc_gen_crd.k8s_topology) = KT_MAP];`, ""),
			"field testdata.Spec.f: k8s_topology KT_MAP: requires the k8s_patch merge_key of the field as the list key",
		},
		"map of scalars": {
			crd(`repeated string f = 1 [(protoc_gen_crd.k8s_topology) = KT_MAP, (protoc_gen_crd.k8s_patch) = {merge_key: "name"}];`, ""),
			"field testdata.Spec.f: k8s_topology KT_MAP: requires a list of objects",
		},
		"map with absent key": {
			crd(`repeated Item f = 1 [(protoc_gen_crd.k8s_topology) = KT_MAP, (protoc_gen_crd.k8s_patch) = {merge_key: "absent"}];`, ""),
			`field testdata.Spec.f: k8s_topology KT_MAP: list key "absent" is not a field of the list items`,
		},
		"map with object key": {
			crd(`repeated Item f = 1 [(protoc_gen_crd.k8s_topology) = KT_MAP, (protoc_gen_crd.k8s_patch) = {merge_key: "nested"}];`, ""),
			`field testdata.Spec.f: k8s_topology KT_MAP: list key "nested" must be a scalar field, it renders as an object`,
		},
		"map with optional key": {
			crd(`repeated Item f = 1 [(protoc_gen_crd.k8s_topology) = KT_MAP, (protoc_gen_crd.k8s_patch) = {merge_key: "optional_name"}];`, ""),
			`field testdata.Spec.f: k8s_topology KT_MAP: list key "optional_name" must not be optional`,
		},
		"atomic timestamp": {
			crd(`google.protobuf.Timestamp f = 1 [(protoc_gen_crd.k8s_topology) = KT_ATOMIC];`, ""),
			`field testdata.Spec.f: k8s_topology KT_ATOMIC: requires an object, a map or a list, the field renders as "string"`,
		},
		"map on root": {
			crd(``, `, field_topologies: [{field_path: "status", topology: KT_MAP}]`),
			"field testdata.Kind.status: k8s_topology KT_MAP: requires a list, the field renders as an object",
		},
		"selector without target": {
			crd(``, `, field_topologies: [{field_path: "spec", topology: KT_ATOMIC}, {topology: KT_ATOMIC}]`),
			"Kind: field_topologies[1]: selector without protobuf_type or field_path",
		},
		"selector with empty field_path": {
			crd(``, `, field_topologies: [{field_path: "", topology: KT_ATOMIC}]`),
			"Kind: field_topologies[0]: empty field_path",
		},
		"selector with empty protobuf_type": {
			crd(``, `, field_topologies: [{field_path: "spec", topology: KT_ATOMIC}, {protobuf_type: "testdata.Spec", topology: KT_ATOMIC}, {protobuf_type: "", topology: KT_ATOMIC}]`),
			"Kind: field_topologies[2]: empty protobuf_type",
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := runProtoSource(t, tc.source)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestTopologyErrorsInSchemaless(t *testing.T) {
	source := `
syntax = "proto3";
package testdata;
option go_package = "github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/pkg/gen/testdata";
import "library/go/k8s/protoc_gen_crd/proto/crd.proto";
message Spec {}
message Status {}
message Kind {
    option (protoc_gen_crd.k8s_crd) = {api_group: "group", kind: "Kind", plural: "kinds", singular: "kind"};
    Spec spec = 1;
    Status status = 2 [(protoc_gen_crd.k8s_topology) = KT_SET];
}
`
	_, err := runProtoSource(t, source, WithSchemalessCrd(true))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "field testdata.Kind.status: k8s_topology KT_SET: requires a list, the field renders as an object")
}

func TestPatchExternals(t *testing.T) {
	response := parseProto(t, "testdata/patch_externals.proto", WithClientSchema(true), WithGeneratingMergeKeys(true))
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
	assert.Equal(t, map[string]any{
		"type": "array",
		"items": map[string]any{
			"type":                 "object",
			"nullable":             true,
			"additionalProperties": false,
			"properties": map[string]any{
				fakeMergeKey: map[string]any{
					"description": fakeMergeKeyDescription,
					"type":        "string",
				},
			},
		},
		patchMergeKeyField:      fakeMergeKey,
		patchMergeStrategyField: "merge",
	}, container.Key("nested_without_merge_key").Value())
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
	assert.Equal(t, map[string]any{"type": "string"}, spec.Key("duration").Value())
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
