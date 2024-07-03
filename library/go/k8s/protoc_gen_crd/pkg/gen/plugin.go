package gen

import (
	"fmt"
	"regexp"

	"github.com/google/gnostic/compiler"
	v3 "github.com/google/gnostic/openapiv3"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	"gopkg.in/yaml.v3"

	crd "github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/proto"
)

// Plugin transforms input files into OpenAPIv3 schema.
type Plugin struct {
	// IsClientSchema controls output mode.
	// If IsClientSchema set to true, then spec will be generated with kustomize
	// validation specifics in mind, and kustomize merge annotations set.
	// In all other cases valid CRD would be produced and all kustomize annotations
	// would be omitted.
	IsClientSchema *bool
}

type PluginOption func(p *Plugin)

func WithClientSchema(isClientSchema bool) PluginOption {
	return func(p *Plugin) {
		p.IsClientSchema = &isClientSchema
	}
}

func (p Plugin) Run(plugin *protogen.Plugin) error {
	isClientSchema := p.IsClientSchema != nil && *p.IsClientSchema

	plugin.SupportedFeatures |= uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	for _, file := range plugin.Files {
		schema := &Schema{
			visitedSchemas:    make(map[string]struct{}),
			typesStack:        map[string]bool{},
			schemas:           &v3.SchemasOrReferences{AdditionalProperties: make([]*v3.NamedSchemaOrReference, 0)},
			linterRulePattern: regexp.MustCompile(`\(-- .* --\)`),
			isClientSchema:    isClientSchema,
		}
		schema.addSchemas(file.Messages)

		if !schema.OneCrd() {
			continue
		}

		// NOTE (torkve) we manually construct YAML nodes to ensure nice looking order of fields

		header := compiler.NewMappingNode()
		header.Content = append(header.Content, compiler.NewScalarNodeForString("apiVersion"))
		header.Content = append(header.Content, compiler.NewScalarNodeForString("apiextensions.k8s.io/v1"))

		header.Content = append(header.Content, compiler.NewScalarNodeForString("kind"))
		header.Content = append(header.Content, compiler.NewScalarNodeForString("CustomResourceDefinition"))

		metadata := compiler.NewMappingNode()

		metadata.Content = append(metadata.Content, compiler.NewScalarNodeForString("name"))
		metadata.Content = append(metadata.Content, compiler.NewScalarNodeForString(schema.metadata.Plural+"."+schema.metadata.ApiGroup))

		metadata.Content = append(metadata.Content, compiler.NewScalarNodeForString("annotations"))
		metadata.Content = append(metadata.Content, compiler.NewMappingNode())

		header.Content = append(header.Content, compiler.NewScalarNodeForString("metadata"))
		header.Content = append(header.Content, metadata)

		spec := compiler.NewMappingNode()
		spec.Content = append(spec.Content, compiler.NewScalarNodeForString("group"))
		spec.Content = append(spec.Content, compiler.NewScalarNodeForString(schema.metadata.ApiGroup))

		spec.Content = append(spec.Content, compiler.NewScalarNodeForString("scope"))
		var scope string
		if schema.metadata.Scope == crd.ScopeType_ST_CLUSTER {
			scope = "Cluster"
		} else {
			scope = "Namespaced"
		}
		spec.Content = append(spec.Content, compiler.NewScalarNodeForString(scope))

		names := compiler.NewMappingNode()

		names.Content = append(names.Content, compiler.NewScalarNodeForString("kind"))
		names.Content = append(names.Content, compiler.NewScalarNodeForString(schema.metadata.Kind))

		names.Content = append(names.Content, compiler.NewScalarNodeForString("listKind"))
		names.Content = append(names.Content, compiler.NewScalarNodeForString(schema.metadata.Kind+"List"))

		names.Content = append(names.Content, compiler.NewScalarNodeForString("plural"))
		names.Content = append(names.Content, compiler.NewScalarNodeForString(schema.metadata.Plural))

		names.Content = append(names.Content, compiler.NewScalarNodeForString("singular"))
		names.Content = append(names.Content, compiler.NewScalarNodeForString(schema.metadata.Singular))

		names.Content = append(names.Content, compiler.NewScalarNodeForString("shortNames"))
		names.Content = append(names.Content, compiler.NewSequenceNodeForStringArray(schema.metadata.ShortNames))

		names.Content = append(names.Content, compiler.NewScalarNodeForString("categories"))
		names.Content = append(names.Content, compiler.NewSequenceNodeForStringArray(schema.metadata.Categories))

		spec.Content = append(spec.Content, compiler.NewScalarNodeForString("names"))
		spec.Content = append(spec.Content, names)

		versions := compiler.NewSequenceNode()

		versionV1 := compiler.NewMappingNode()
		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("name"))
		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("v1")) // TODO (torkve) support multiple versions and custom naming

		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("served"))
		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForBool(true))

		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("storage"))
		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForBool(true))

		subresources := compiler.NewMappingNode()
		// FIXME (torkve) currently we are enforcing exactly one subresource named "status". It must be
		//                configurable via field annotations
		subresources.Content = append(subresources.Content, compiler.NewScalarNodeForString("status"))
		subresources.Content = append(subresources.Content, compiler.NewMappingNode())

		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("subresources"))
		versionV1.Content = append(versionV1.Content, subresources)

		versionV1Schema := compiler.NewMappingNode()
		versionV1Schema.Content = append(versionV1Schema.Content, compiler.NewScalarNodeForString("openAPIV3Schema"))
		versionV1Schema.Content = append(versionV1Schema.Content, schema.schemas.ToRawInfo().Content[1])

		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("schema"))
		versionV1.Content = append(versionV1.Content, versionV1Schema)

		additionalColumns := compiler.NewSequenceNode()
		for _, column := range schema.metadata.AdditionalColumns {
			additionalColumns.Content = append(additionalColumns.Content, renderAdditionalColumn(column))
		}

		versionV1.Content = append(versionV1.Content, compiler.NewScalarNodeForString("additionalPrinterColumns"))
		versionV1.Content = append(versionV1.Content, additionalColumns)

		versions.Content = append(versions.Content, versionV1)

		spec.Content = append(spec.Content, compiler.NewScalarNodeForString("versions"))
		spec.Content = append(spec.Content, versions)

		header.Content = append(header.Content, compiler.NewScalarNodeForString("spec"))
		header.Content = append(header.Content, spec)

		rawInfo := &yaml.Node{
			Kind:        yaml.DocumentNode,
			Style:       0,
			Content:     []*yaml.Node{header},
			HeadComment: "Generated by protoc-gen-crd from " + file.GeneratedFilenamePrefix + ".proto",
		}

		suffix := ".crd.yaml"
		if schema.isClientSchema {
			suffix = ".kustomize.yaml"
		}
		outputFileName := file.GeneratedFilenamePrefix + suffix
		outputFile := plugin.NewGeneratedFile(outputFileName, "")
		e := yaml.NewEncoder(outputFile)
		e.SetIndent(2)
		if err := e.Encode(rawInfo); err != nil {
			return fmt.Errorf("failed to marshal yaml: %s", err.Error())
		}
	}
	return nil
}
