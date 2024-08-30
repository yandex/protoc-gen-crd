package gen

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/gnostic/compiler"
	v3 "github.com/google/gnostic/openapiv3"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v3"

	crd "github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/proto"
)

const (
	patchMergeKeyField      = "x-kubernetes-patch-merge-key"
	patchMergeStrategyField = "x-kubernetes-patch-strategy"

	intOrStringField           = "x-kubernetes-int-or-string"
	preserveUnknownFieldsField = "x-kubernetes-preserve-unknown-fields"
)

var (
	columnTypeName = map[crd.ColumnType]string{
		crd.ColumnType_CT_INTEGER: "integer",
		crd.ColumnType_CT_NUMBER:  "number",
		crd.ColumnType_CT_STRING:  "string",
		crd.ColumnType_CT_BOOLEAN: "boolean",
		crd.ColumnType_CT_DATE:    "date",
	}
	columnFormatName = map[crd.ColumnFormat]string{
		crd.ColumnFormat_CF_INT32:    "int32",
		crd.ColumnFormat_CF_INT64:    "int64",
		crd.ColumnFormat_CF_FLOAT:    "float",
		crd.ColumnFormat_CF_DOUBLE:   "double",
		crd.ColumnFormat_CF_BYTE:     "byte",
		crd.ColumnFormat_CF_DATE:     "date",
		crd.ColumnFormat_CF_DATETIME: "date-time",
		crd.ColumnFormat_CF_PASSWORD: "password",
	}
)

var intOrStringExtension = &v3.NamedAny{
	Name:  intOrStringField,
	Value: &v3.Any{Yaml: "true"},
}
var preserveUnknownFieldsExtension = &v3.NamedAny{
	Name:  preserveUnknownFieldsField,
	Value: &v3.Any{Yaml: "true"},
}

var opaqueSchema = &v3.SchemaOrReference{
	Oneof: &v3.SchemaOrReference_Schema{Schema: &v3.Schema{
		Type:                   "object",
		SpecificationExtension: []*v3.NamedAny{preserveUnknownFieldsExtension},
	}},
}

func makeStringAny(s string) *v3.Any {
	return &v3.Any{Yaml: fmt.Sprintf("%q", s)}
}

func makeIntAny(i int32) *v3.Any {
	return &v3.Any{Yaml: fmt.Sprint(i)}
}

func getMessageExtension[T protoreflect.ExtensionType](message *protogen.Field, xt T) any {
	if message == nil {
		return nil
	}
	extension := proto.GetExtension(message.Desc.Options(), xt)
	if extension == nil || extension == xt.InterfaceOf(xt.Zero()) {
		return nil
	}
	return extension
}

type Schema struct {
	visitedSchemas map[string]struct{}
	typesStack     map[string]bool
	schemas        *v3.SchemasOrReferences
	metadata       *crd.K8SCRD
	isClientSchema bool

	linterRulePattern *regexp.Regexp
	typePatchRules    map[string]*crd.K8SPatch
	fieldPatchRules   *RadixTree[*crd.K8SPatch]
}

func fullMessageTypeName(message protoreflect.MessageDescriptor) string {
	return "." + string(message.ParentFile().Package()) + "." + string(message.FullName())
}

func (s *Schema) formatFieldName(field *protogen.Field) string {
	return string(field.Desc.Name())
}

func (s *Schema) formatMessageName(message *protogen.Message) string {
	return string(message.Desc.Name())
}

func (s *Schema) formatMessageRef(name string) string {
	return name
}

func (s *Schema) OneCrd() bool {
	return len(s.schemas.AdditionalProperties) == 1
}

func (s *Schema) getPatchAnnotation(field *protogen.Field, fieldTree *RadixTree[*crd.K8SPatch]) *crd.K8SPatch {
	if rule, ok := fieldTree.Value(); ok {
		return *rule
	}

	if field.Message != nil {
		message := string(field.Message.Desc.FullName())
		if rule, ok := s.typePatchRules[message]; ok {
			return rule
		}
	} else if field.Enum != nil {
		message := string(field.Enum.Desc.FullName())
		if rule, ok := s.typePatchRules[message]; ok {
			return rule
		}
	}

	ext := getMessageExtension(field, crd.E_K8SPatch)
	if ext == nil {
		return nil
	}
	return ext.(*crd.K8SPatch)

}

func (s *Schema) buildPatchRules() {
	typeMap := make(map[string]*crd.K8SPatch)
	pathMap := NewRadixTree[*crd.K8SPatch]()
	for _, rule := range s.metadata.GetFieldPatchStrategies() {
		switch target := rule.GetTarget().(type) {
		case *crd.K8SPatchSelector_ProtobufType:
			typeMap[target.ProtobufType] = rule.K8SPatch
		case *crd.K8SPatchSelector_FieldPath:
			pathMap.Add(strings.Split(target.FieldPath, "."), rule.K8SPatch)
		default:
			log.Printf("(TODO) Patch strategy has unknown target type: %T", target)
		}
	}

	s.typePatchRules = typeMap
	s.fieldPatchRules = pathMap
}

func (s *Schema) needAddToSchema(message *protogen.Field) bool {
	if message == nil {
		return false
	}
	schemaSide := getMessageExtension(message, crd.E_Schema)
	if schemaSide == nil || schemaSide.(crd.SchemaSide) != crd.SchemaSide_SS_CLIENT {
		return true
	}
	return s.isClientSchema
}

func (s *Schema) makeSpecificationExtension(patchAnnotation *crd.K8SPatch) []*v3.NamedAny {
	var out []*v3.NamedAny
	if patchAnnotation == nil {
		return out
	}
	if len(patchAnnotation.MergeKey) > 0 && s.isClientSchema {
		out = append(out, &v3.NamedAny{
			Name: patchMergeKeyField,
			Value: &v3.Any{
				Yaml: patchAnnotation.MergeKey,
			},
		})
	}
	if len(patchAnnotation.MergeStrategy) > 0 && s.isClientSchema {
		out = append(out, &v3.NamedAny{
			Name: patchMergeStrategyField,
			Value: &v3.Any{
				Yaml: patchAnnotation.MergeStrategy,
			},
		})
	}
	return out
}

func (s *Schema) shouldVisitSchema(typeName string) bool {
	_, ok := s.visitedSchemas[typeName]
	if ok {
		return false
	}
	s.visitedSchemas[typeName] = struct{}{}
	return true
}

func (s *Schema) filterCommentString(c protogen.Comments, removeNewLines bool) string {
	comment := string(c)
	if removeNewLines {
		comment = strings.Replace(comment, "\n", "", -1)
	}
	comment = s.linterRulePattern.ReplaceAllString(comment, "")
	return strings.TrimSpace(comment)
}

func (s *Schema) schemaReferenceForTypeName(typeName string) string {
	parts := strings.Split(typeName, ".")
	lastPart := parts[len(parts)-1]
	return "#/components/schemas/" + s.formatMessageRef(lastPart)
}

func (s *Schema) schemaOrReferenceForTypeOrMessage(typeName string, message *protogen.Message, fieldTree *RadixTree[*crd.K8SPatch]) *v3.SchemaOrReference {
	switch typeName {

	// TODO (torkve) Create oneof here: we probably should allow user to provide either formatted string (RFC3339 etc)
	//	             or proto-compatible struct: to support direct passing objects from dctl.
	//               But gnostic currently doesn't support Type to be an array.

	case "google.protobuf.Int32Value":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "integer", Format: "int32", Nullable: true}}}
	case "google.protobuf.UInt32Value":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "integer", Format: "uint32", Nullable: true}}}

	case "google.protobuf.Int64Value":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Format: "int64", Nullable: true, SpecificationExtension: []*v3.NamedAny{intOrStringExtension}}}}

	case "google.protobuf.UInt64Value":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Format: "uint64", Nullable: true, SpecificationExtension: []*v3.NamedAny{intOrStringExtension}}}}

	case "google.protobuf.FloatValue":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "number", Format: "float", Nullable: true}}}

	case "google.protobuf.DoubleValue":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "number", Format: "double", Nullable: true}}}

	case "google.protobuf.StringValue":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Nullable: true}}}

	case "google.protobuf.BoolValue":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "boolean", Nullable: true}}}

	case "google.protobuf.Timestamp":
		// Timestamps are serialized as strings
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Format: "RFC3339"}}}

	case "google.type.Date":
		// Dates are serialized as strings
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Format: "date"}}}

	case "google.type.DateTime":
		// DateTimes are serialized as strings
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Format: "date-time"}}}

	case "google.type.Duration":
		// Duration are serialized as strings
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Format: "duration"}}}

	case "google.protobuf.Value":
		// Value is equivalent to a JSON value with any type
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					SpecificationExtension: []*v3.NamedAny{preserveUnknownFieldsExtension},
				},
			},
		}
	case "google.protobuf.ListValue":
		// ListValue is equivalent to an array of Values
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Type: "array",
					Items: &v3.ItemsItem{
						SchemaOrReference: []*v3.SchemaOrReference{
							{
								Oneof: &v3.SchemaOrReference_Schema{
									Schema: &v3.Schema{
										SpecificationExtension: []*v3.NamedAny{preserveUnknownFieldsExtension},
									},
								},
							},
						},
					},
				},
			},
		}
	case "google.protobuf.Struct":
		// Struct is equivalent to a JSON object
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Type:                   "object",
					SpecificationExtension: []*v3.NamedAny{preserveUnknownFieldsExtension},
				},
			},
		}

	case "google.protobuf.Empty":
		// Empty is close to JSON undefined than null, so ignore this field
		return nil //&v3.SchemaOrReference{Oneof: &v3.SchemaOrReference_Schema{Schema: &v3.Schema{Type: "null"}}}

	case "google.protobuf.Any":
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Type:     "object",
					Nullable: true,
					Properties: &v3.Properties{AdditionalProperties: []*v3.NamedSchemaOrReference{
						{
							Name: "@type",
							Value: &v3.SchemaOrReference{Oneof: &v3.SchemaOrReference_Schema{
								Schema: &v3.Schema{Type: "string"},
							}},
						},
					}},
					Required: []string{"@type"},
					AdditionalProperties: &v3.AdditionalPropertiesItem{Oneof: &v3.AdditionalPropertiesItem_Boolean{
						Boolean: true,
					}},
					SpecificationExtension: []*v3.NamedAny{preserveUnknownFieldsExtension},
				},
			},
		}
	default:
		return s.schemaForMessage(message, false, fieldTree)
	}
}

func (s *Schema) schemaOrReferenceForField(field *protogen.Field, fieldTree *RadixTree[*crd.K8SPatch]) *v3.SchemaOrReference {
	if !s.needAddToSchema(field) {
		return nil
	}
	patchAnnotation := s.getPatchAnnotation(field, fieldTree)

	if field.Desc.IsMap() {
		mapMessage := field.Message.Fields[1]
		return &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "object",
					AdditionalProperties: &v3.AdditionalPropertiesItem{
						Oneof: &v3.AdditionalPropertiesItem_SchemaOrReference{
							SchemaOrReference: s.schemaOrReferenceForField(mapMessage, fieldTree),
						},
					},
					SpecificationExtension: s.makeSpecificationExtension(patchAnnotation),
				},
			},
		}
	}

	var kindSchema *v3.SchemaOrReference

	fieldDescription := s.filterCommentString(field.Comments.Leading, true)
	primitiveFieldIsNullable := field.Desc.HasPresence()

	kind := field.Desc.Kind()

	switch kind {

	case protoreflect.MessageKind:
		typeName := string(field.Desc.Message().FullName())
		kindSchema = s.schemaOrReferenceForTypeOrMessage(typeName, field.Message, fieldTree)
		if kindSchema == nil {
			return nil
		}

	case protoreflect.StringKind:
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Description: fieldDescription, Nullable: primitiveFieldIsNullable}}}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Fixed64Kind:
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "integer", Format: kind.String(), Description: fieldDescription, Nullable: primitiveFieldIsNullable}}}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind:
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Format:                 kind.String(),
					SpecificationExtension: []*v3.NamedAny{intOrStringExtension},
					Description:            fieldDescription,
					Nullable:               primitiveFieldIsNullable,
				},
			},
		}
	case protoreflect.EnumKind:
		enumValues := make([]*v3.Any, 0, len(field.Enum.Values)*2)
		for _, enumValue := range field.Enum.Values {
			enumValues = append(enumValues, makeIntAny(int32(enumValue.Desc.Number())), makeStringAny(string(enumValue.Desc.Name())))
		}
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Format:                 "enum",
					SpecificationExtension: []*v3.NamedAny{intOrStringExtension},
					Description:            fieldDescription,
					Enum:                   enumValues,
					Nullable:               primitiveFieldIsNullable,
				},
			},
		}

	case protoreflect.BoolKind:
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Type:        "boolean",
					Description: fieldDescription,
					Nullable:    primitiveFieldIsNullable,
				}}}

	case protoreflect.FloatKind, protoreflect.DoubleKind:
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "number", Format: kind.String(), Description: fieldDescription, Nullable: primitiveFieldIsNullable}}}

	case protoreflect.BytesKind:
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{Type: "string", Format: "bytes", Description: fieldDescription, Nullable: primitiveFieldIsNullable}}}

	default:
		log.Printf("(TODO) Unsupported field type: %+v", fullMessageTypeName(field.Desc.Message()))
		return nil
	}

	if field.Desc.IsList() {
		kindSchema = &v3.SchemaOrReference{
			Oneof: &v3.SchemaOrReference_Schema{
				Schema: &v3.Schema{
					Type:  "array",
					Items: &v3.ItemsItem{SchemaOrReference: []*v3.SchemaOrReference{kindSchema}},
				},
			},
		}
	}

	if schema := kindSchema.GetSchema(); patchAnnotation != nil && schema != nil {
		schema.SpecificationExtension = append(schema.SpecificationExtension, s.makeSpecificationExtension(patchAnnotation)...)
	}

	return kindSchema
}

func (s *Schema) schemaForMessage(message *protogen.Message, isRoot bool, fieldTree *RadixTree[*crd.K8SPatch]) *v3.SchemaOrReference {
	typename := fullMessageTypeName(message.Desc)

	if s.typesStack[typename] {
		return opaqueSchema
	}

	messageDescription := s.filterCommentString(message.Comments.Leading, true)
	definitionProperties := &v3.Properties{
		AdditionalProperties: make([]*v3.NamedSchemaOrReference, 0),
	}

	s.typesStack[typename] = true
	defer delete(s.typesStack, typename)

	// TODO (torkve) process oneof's separately

	for _, field := range message.Fields {
		// The field is either described by a reference or a schema.
		fieldSchema := s.schemaOrReferenceForField(field, fieldTree.Child(s.formatFieldName(field)))
		if fieldSchema == nil {
			continue
		}

		if schema, ok := fieldSchema.Oneof.(*v3.SchemaOrReference_Schema); ok {
			// Get the field description from the comments.
			schema.Schema.Description = s.filterCommentString(field.Comments.Leading, true)
		}

		definitionProperties.AdditionalProperties = append(
			definitionProperties.AdditionalProperties,
			&v3.NamedSchemaOrReference{
				Name:  s.formatFieldName(field),
				Value: fieldSchema,
			},
		)
	}

	reservedNames := message.Desc.ReservedNames()
	for idx := 0; idx < reservedNames.Len(); idx++ {
		reservedName := string(reservedNames.Get(idx))
		definitionProperties.AdditionalProperties = append(
			definitionProperties.AdditionalProperties,
			&v3.NamedSchemaOrReference{
				Name: reservedName,
				Value: &v3.SchemaOrReference{
					Oneof: &v3.SchemaOrReference_Schema{Schema: &v3.Schema{
						Deprecated:             s.isClientSchema,
						SpecificationExtension: []*v3.NamedAny{preserveUnknownFieldsExtension},
					}},
				},
			},
		)
	}

	var additionalProperties *v3.AdditionalPropertiesItem
	if s.isClientSchema && !isRoot {
		additionalProperties = &v3.AdditionalPropertiesItem{
			Oneof: &v3.AdditionalPropertiesItem_Boolean{
				Boolean: false,
			},
		}
	}
	return &v3.SchemaOrReference{
		Oneof: &v3.SchemaOrReference_Schema{
			Schema: &v3.Schema{
				Type:                 "object",
				Nullable:             !isRoot,
				Description:          messageDescription,
				Properties:           definitionProperties,
				AdditionalProperties: additionalProperties,
			},
		},
	}
}

func (s *Schema) addSchemas(messages []*protogen.Message) {
	for _, message := range messages {
		if message.Messages != nil {
			s.addSchemas(message.Messages)
		}

		typeName := fullMessageTypeName(message.Desc)
		if !s.shouldVisitSchema(typeName) {
			continue
		}

		xt := crd.E_K8SCrd
		extension := proto.GetExtension(message.Desc.Options(), xt)
		if extension == nil || extension == xt.InterfaceOf(xt.Zero()) {
			continue
		}
		s.metadata = extension.(*crd.K8SCRD)
		s.buildPatchRules()

		s.schemas.AdditionalProperties = append(s.schemas.AdditionalProperties,
			&v3.NamedSchemaOrReference{
				Name:  s.formatMessageName(message),
				Value: s.schemaForMessage(message, true, s.fieldPatchRules),
			},
		)
	}
}

func renderAdditionalColumn(column *crd.PrinterColumn) *yaml.Node {
	node := compiler.NewMappingNode()

	node.Content = append(node.Content, compiler.NewScalarNodeForString("name"))
	node.Content = append(node.Content, compiler.NewScalarNodeForString(column.Name))

	if column.Type != crd.ColumnType_CT_NONE {
		node.Content = append(node.Content, compiler.NewScalarNodeForString("type"))
		node.Content = append(node.Content, compiler.NewScalarNodeForString(columnTypeName[column.Type]))
	}

	if column.Format != crd.ColumnFormat_CF_NONE {
		node.Content = append(node.Content, compiler.NewScalarNodeForString("format"))
		node.Content = append(node.Content, compiler.NewScalarNodeForString(columnFormatName[column.Format]))
	}

	if column.Description != "" {
		node.Content = append(node.Content, compiler.NewScalarNodeForString("description"))
		node.Content = append(node.Content, compiler.NewScalarNodeForString(column.Description))
	}

	node.Content = append(node.Content, compiler.NewScalarNodeForString("jsonPath"))
	node.Content = append(node.Content, compiler.NewScalarNodeForString(column.JsonPath))

	if column.Priority != 0 {
		node.Content = append(node.Content, compiler.NewScalarNodeForString("priority"))
		node.Content = append(node.Content, compiler.NewScalarNodeForInt(int64(column.Priority)))
	}

	return node
}
