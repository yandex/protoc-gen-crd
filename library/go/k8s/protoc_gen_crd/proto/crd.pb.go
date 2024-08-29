// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.5
// source: library/go/k8s/protoc_gen_crd/proto/crd.proto

package crd

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ColumnType int32

const (
	// Unspecified column type
	ColumnType_CT_NONE ColumnType = 0
	// Non-floating-point numbers
	ColumnType_CT_INTEGER ColumnType = 1
	// Floating point numbers
	ColumnType_CT_NUMBER ColumnType = 2
	// Strings
	ColumnType_CT_STRING ColumnType = 3
	// true or false
	ColumnType_CT_BOOLEAN ColumnType = 4
	// Rendered differentially as time since this timestamp
	ColumnType_CT_DATE ColumnType = 5
)

// Enum value maps for ColumnType.
var (
	ColumnType_name = map[int32]string{
		0: "CT_NONE",
		1: "CT_INTEGER",
		2: "CT_NUMBER",
		3: "CT_STRING",
		4: "CT_BOOLEAN",
		5: "CT_DATE",
	}
	ColumnType_value = map[string]int32{
		"CT_NONE":    0,
		"CT_INTEGER": 1,
		"CT_NUMBER":  2,
		"CT_STRING":  3,
		"CT_BOOLEAN": 4,
		"CT_DATE":    5,
	}
)

func (x ColumnType) Enum() *ColumnType {
	p := new(ColumnType)
	*p = x
	return p
}

func (x ColumnType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ColumnType) Descriptor() protoreflect.EnumDescriptor {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[0].Descriptor()
}

func (ColumnType) Type() protoreflect.EnumType {
	return &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[0]
}

func (x ColumnType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ColumnType.Descriptor instead.
func (ColumnType) EnumDescriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{0}
}

type ColumnFormat int32

const (
	// Unspecified column format, the default one for the ColumnType will be used.
	ColumnFormat_CF_NONE     ColumnFormat = 0
	ColumnFormat_CF_INT32    ColumnFormat = 1
	ColumnFormat_CF_INT64    ColumnFormat = 2
	ColumnFormat_CF_FLOAT    ColumnFormat = 3
	ColumnFormat_CF_DOUBLE   ColumnFormat = 4
	ColumnFormat_CF_BYTE     ColumnFormat = 5
	ColumnFormat_CF_DATE     ColumnFormat = 6
	ColumnFormat_CF_DATETIME ColumnFormat = 7
	ColumnFormat_CF_PASSWORD ColumnFormat = 8
)

// Enum value maps for ColumnFormat.
var (
	ColumnFormat_name = map[int32]string{
		0: "CF_NONE",
		1: "CF_INT32",
		2: "CF_INT64",
		3: "CF_FLOAT",
		4: "CF_DOUBLE",
		5: "CF_BYTE",
		6: "CF_DATE",
		7: "CF_DATETIME",
		8: "CF_PASSWORD",
	}
	ColumnFormat_value = map[string]int32{
		"CF_NONE":     0,
		"CF_INT32":    1,
		"CF_INT64":    2,
		"CF_FLOAT":    3,
		"CF_DOUBLE":   4,
		"CF_BYTE":     5,
		"CF_DATE":     6,
		"CF_DATETIME": 7,
		"CF_PASSWORD": 8,
	}
)

func (x ColumnFormat) Enum() *ColumnFormat {
	p := new(ColumnFormat)
	*p = x
	return p
}

func (x ColumnFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ColumnFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[1].Descriptor()
}

func (ColumnFormat) Type() protoreflect.EnumType {
	return &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[1]
}

func (x ColumnFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ColumnFormat.Descriptor instead.
func (ColumnFormat) EnumDescriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{1}
}

type ScopeType int32

const (
	// Object must always be attached to some namespace, and subjects to namespace-scoped RBAC roles
	ScopeType_ST_NAMESPACED ScopeType = 0
	// Object would exist in global scope, and only ClusterRole RBAC would be applied to id
	ScopeType_ST_CLUSTER ScopeType = 1
)

// Enum value maps for ScopeType.
var (
	ScopeType_name = map[int32]string{
		0: "ST_NAMESPACED",
		1: "ST_CLUSTER",
	}
	ScopeType_value = map[string]int32{
		"ST_NAMESPACED": 0,
		"ST_CLUSTER":    1,
	}
)

func (x ScopeType) Enum() *ScopeType {
	p := new(ScopeType)
	*p = x
	return p
}

func (x ScopeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ScopeType) Descriptor() protoreflect.EnumDescriptor {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[2].Descriptor()
}

func (ScopeType) Type() protoreflect.EnumType {
	return &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[2]
}

func (x ScopeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ScopeType.Descriptor instead.
func (ScopeType) EnumDescriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{2}
}

type SchemaSide int32

const (
	// Generate field in both CRD modes: client and server
	SchemaSide_SS_COMMON SchemaSide = 0
	// Generate field only for client CRD to apply kustomize and validate objects
	SchemaSide_SS_CLIENT SchemaSide = 1
)

// Enum value maps for SchemaSide.
var (
	SchemaSide_name = map[int32]string{
		0: "SS_COMMON",
		1: "SS_CLIENT",
	}
	SchemaSide_value = map[string]int32{
		"SS_COMMON": 0,
		"SS_CLIENT": 1,
	}
)

func (x SchemaSide) Enum() *SchemaSide {
	p := new(SchemaSide)
	*p = x
	return p
}

func (x SchemaSide) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SchemaSide) Descriptor() protoreflect.EnumDescriptor {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[3].Descriptor()
}

func (SchemaSide) Type() protoreflect.EnumType {
	return &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes[3]
}

func (x SchemaSide) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SchemaSide.Descriptor instead.
func (SchemaSide) EnumDescriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{3}
}

type PrinterColumn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Human readable name of the column
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Column type, affects how column content is interpreted before printing.
	Type ColumnType `protobuf:"varint,2,opt,name=type,proto3,enum=protoc_gen_crd.ColumnType" json:"type,omitempty"`
	// Column format hint.
	// See https://github.com/OAI/OpenAPI-Specification/blob/7cc8f4c4e742a20687fa65ace54ed32fcb8c6df0/versions/2.0.md#data-types
	// for full list of types and their supported formats
	Format ColumnFormat `protobuf:"varint,3,opt,name=format,proto3,enum=protoc_gen_crd.ColumnFormat" json:"format,omitempty"`
	// Human readable description of the column
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// Very simple JSON Path expression (only path without any logic) that would be evaluated for object to retrieve column value
	JsonPath string `protobuf:"bytes,5,opt,name=json_path,json=jsonPath,proto3" json:"json_path,omitempty"`
	// Column priority. The less value, the higher priority.
	// Columns with priority > 0 can be omitted from output if no sufficient space on the screen for printing
	Priority int32 `protobuf:"varint,6,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *PrinterColumn) Reset() {
	*x = PrinterColumn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrinterColumn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrinterColumn) ProtoMessage() {}

func (x *PrinterColumn) ProtoReflect() protoreflect.Message {
	mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrinterColumn.ProtoReflect.Descriptor instead.
func (*PrinterColumn) Descriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{0}
}

func (x *PrinterColumn) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PrinterColumn) GetType() ColumnType {
	if x != nil {
		return x.Type
	}
	return ColumnType_CT_NONE
}

func (x *PrinterColumn) GetFormat() ColumnFormat {
	if x != nil {
		return x.Format
	}
	return ColumnFormat_CF_NONE
}

func (x *PrinterColumn) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PrinterColumn) GetJsonPath() string {
	if x != nil {
		return x.JsonPath
	}
	return ""
}

func (x *PrinterColumn) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type K8SCRD struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// CRD API group name. Used to group your CRDs, must conform to DNS domain-like format, like "apis.example.com"
	ApiGroup string `protobuf:"bytes,1,opt,name=api_group,json=apiGroup,proto3" json:"api_group,omitempty"`
	// Resource kind. Used to identify resource type in serialized specs. "CamelCasedSingular" should be used.
	Kind string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	// Singular name for the object type, used as an alias for kind in CLI, and for display.
	// Normally it is lowercased kind, e.g. "camelcasedsingular"
	Singular string `protobuf:"bytes,3,opt,name=singular,proto3" json:"singular,omitempty"`
	// Plural name for the object type, used in API URLs.
	// Normally it is singular name with plural suffix added, e.g. "camelcasedsingulars"
	Plural string `protobuf:"bytes,4,opt,name=plural,proto3" json:"plural,omitempty"`
	// Short aliases that can be used in CLI to refer to resource, if you want.
	// E.g. ["camelcased", "ccs"]
	ShortNames []string `protobuf:"bytes,5,rep,name=short_names,json=shortNames,proto3" json:"short_names,omitempty"`
	// Categories the custom resource belongs to.
	// Used in CLI to list resources from some category regardless of their kind.
	// E.g. ["all", "mycompany", "myservice"]
	Categories []string `protobuf:"bytes,6,rep,name=categories,proto3" json:"categories,omitempty"`
	// Additional columns available in kubectl get
	AdditionalColumns []*PrinterColumn `protobuf:"bytes,7,rep,name=additional_columns,json=additionalColumns,proto3" json:"additional_columns,omitempty"`
	// Object kind scope
	Scope ScopeType `protobuf:"varint,8,opt,name=scope,proto3,enum=protoc_gen_crd.ScopeType" json:"scope,omitempty"`
	// List of fully qualified protobuf message types (including proto package name)
	// or field names (e.g. "spec.field"), that should have custom kustomize patch strategy applied.
	// This option overrides K8sPatch option set on the field, and field targeted rules has precedence over type targeted ones.
	FieldPatchStrategies []*K8SPatchSelector `protobuf:"bytes,9,rep,name=field_patch_strategies,json=fieldPatchStrategies,proto3" json:"field_patch_strategies,omitempty"`
}

func (x *K8SCRD) Reset() {
	*x = K8SCRD{}
	if protoimpl.UnsafeEnabled {
		mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *K8SCRD) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SCRD) ProtoMessage() {}

func (x *K8SCRD) ProtoReflect() protoreflect.Message {
	mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SCRD.ProtoReflect.Descriptor instead.
func (*K8SCRD) Descriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{1}
}

func (x *K8SCRD) GetApiGroup() string {
	if x != nil {
		return x.ApiGroup
	}
	return ""
}

func (x *K8SCRD) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *K8SCRD) GetSingular() string {
	if x != nil {
		return x.Singular
	}
	return ""
}

func (x *K8SCRD) GetPlural() string {
	if x != nil {
		return x.Plural
	}
	return ""
}

func (x *K8SCRD) GetShortNames() []string {
	if x != nil {
		return x.ShortNames
	}
	return nil
}

func (x *K8SCRD) GetCategories() []string {
	if x != nil {
		return x.Categories
	}
	return nil
}

func (x *K8SCRD) GetAdditionalColumns() []*PrinterColumn {
	if x != nil {
		return x.AdditionalColumns
	}
	return nil
}

func (x *K8SCRD) GetScope() ScopeType {
	if x != nil {
		return x.Scope
	}
	return ScopeType_ST_NAMESPACED
}

func (x *K8SCRD) GetFieldPatchStrategies() []*K8SPatchSelector {
	if x != nil {
		return x.FieldPatchStrategies
	}
	return nil
}

type K8SPatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Field that is used as a key of the struct when "merge" strategy is used for the list. Items with the matching key would be merged into one.
	MergeKey string `protobuf:"bytes,1,opt,name=merge_key,json=mergeKey,proto3" json:"merge_key,omitempty"`
	// Merge strategy.
	// Normally there are only two options: "merge" and "replace", with the former being default for structs and the latter default for lists.
	MergeStrategy string `protobuf:"bytes,2,opt,name=merge_strategy,json=mergeStrategy,proto3" json:"merge_strategy,omitempty"`
}

func (x *K8SPatch) Reset() {
	*x = K8SPatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *K8SPatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SPatch) ProtoMessage() {}

func (x *K8SPatch) ProtoReflect() protoreflect.Message {
	mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SPatch.ProtoReflect.Descriptor instead.
func (*K8SPatch) Descriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{2}
}

func (x *K8SPatch) GetMergeKey() string {
	if x != nil {
		return x.MergeKey
	}
	return ""
}

func (x *K8SPatch) GetMergeStrategy() string {
	if x != nil {
		return x.MergeStrategy
	}
	return ""
}

type K8SPatchSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Target:
	//
	//	*K8SPatchSelector_ProtobufType
	//	*K8SPatchSelector_FieldPath
	Target isK8SPatchSelector_Target `protobuf_oneof:"target"`
	// Patch strategy to apply to the field.
	K8SPatch *K8SPatch `protobuf:"bytes,3,opt,name=k8s_patch,json=k8sPatch,proto3" json:"k8s_patch,omitempty"`
}

func (x *K8SPatchSelector) Reset() {
	*x = K8SPatchSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *K8SPatchSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SPatchSelector) ProtoMessage() {}

func (x *K8SPatchSelector) ProtoReflect() protoreflect.Message {
	mi := &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SPatchSelector.ProtoReflect.Descriptor instead.
func (*K8SPatchSelector) Descriptor() ([]byte, []int) {
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP(), []int{3}
}

func (m *K8SPatchSelector) GetTarget() isK8SPatchSelector_Target {
	if m != nil {
		return m.Target
	}
	return nil
}

func (x *K8SPatchSelector) GetProtobufType() string {
	if x, ok := x.GetTarget().(*K8SPatchSelector_ProtobufType); ok {
		return x.ProtobufType
	}
	return ""
}

func (x *K8SPatchSelector) GetFieldPath() string {
	if x, ok := x.GetTarget().(*K8SPatchSelector_FieldPath); ok {
		return x.FieldPath
	}
	return ""
}

func (x *K8SPatchSelector) GetK8SPatch() *K8SPatch {
	if x != nil {
		return x.K8SPatch
	}
	return nil
}

type isK8SPatchSelector_Target interface {
	isK8SPatchSelector_Target()
}

type K8SPatchSelector_ProtobufType struct {
	// Fully qualified name (with package) of protobuf type of the message to patch, e.g. "example.com.MyMessage"
	ProtobufType string `protobuf:"bytes,1,opt,name=protobuf_type,json=protobufType,proto3,oneof"`
}

type K8SPatchSelector_FieldPath struct {
	// Field path of the field to patch, e.g. "spec.field"
	FieldPath string `protobuf:"bytes,2,opt,name=field_path,json=fieldPath,proto3,oneof"`
}

func (*K8SPatchSelector_ProtobufType) isK8SPatchSelector_Target() {}

func (*K8SPatchSelector_FieldPath) isK8SPatchSelector_Target() {}

var file_library_go_k8s_protoc_gen_crd_proto_crd_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*K8SCRD)(nil),
		Field:         73394821,
		Name:          "protoc_gen_crd.k8s_crd",
		Tag:           "bytes,73394821,opt,name=k8s_crd",
		Filename:      "library/go/k8s/protoc_gen_crd/proto/crd.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*K8SPatch)(nil),
		Field:         73394822,
		Name:          "protoc_gen_crd.k8s_patch",
		Tag:           "bytes,73394822,opt,name=k8s_patch",
		Filename:      "library/go/k8s/protoc_gen_crd/proto/crd.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*SchemaSide)(nil),
		Field:         73394823,
		Name:          "protoc_gen_crd.schema",
		Tag:           "varint,73394823,opt,name=schema,enum=protoc_gen_crd.SchemaSide",
		Filename:      "library/go/k8s/protoc_gen_crd/proto/crd.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// This option is added to a message in proto file, to create CRD from this file.
	// Message type must contain "spec" and "status" fields to conform to k8s CRD rules.
	//
	// optional protoc_gen_crd.K8sCRD k8s_crd = 73394821;
	E_K8SCrd = &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// Add this option to a field to specify how it should be treated during strategic merge.
	// See https://github.com/kubernetes/community/blob/6690abcd6b833f46550f5eaba2ec17a9e39b35c4/contributors/devel/sig-api-machinery/strategic-merge-patch.md
	// about various patch strategies.
	// NOTE: no support for server-side apply yet.
	//
	// optional protoc_gen_crd.K8sPatch k8s_patch = 73394822;
	E_K8SPatch = &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_extTypes[1]
	//	Add this option to a field to specify in which CRD mode field should be present. By default all fields are generated both for server and client schemas
	//
	// optional protoc_gen_crd.SchemaSide schema = 73394823;
	E_Schema = &file_library_go_k8s_protoc_gen_crd_proto_crd_proto_extTypes[2]
)

var File_library_go_k8s_protoc_gen_crd_proto_crd_proto protoreflect.FileDescriptor

var file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x79, 0x2f, 0x67, 0x6f, 0x2f, 0x6b, 0x38, 0x73,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x1a,
	0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xe4, 0x01, 0x0a, 0x0d, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67,
	0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x46,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1b, 0x0a, 0x09, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6a, 0x73, 0x6f, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x22, 0x85, 0x03, 0x0a, 0x06, 0x4b, 0x38, 0x73,
	0x43, 0x52, 0x44, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x70, 0x69, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x70, 0x69, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x6e, 0x67, 0x75, 0x6c, 0x61, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x69, 0x6e, 0x67, 0x75, 0x6c, 0x61, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x4c, 0x0a, 0x12, 0x61, 0x64, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18,
	0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67,
	0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f,
	0x6c, 0x75, 0x6d, 0x6e, 0x52, 0x11, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c,
	0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x12, 0x2f, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f,
	0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x56, 0x0a, 0x16, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69,
	0x65, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x4b, 0x38, 0x73, 0x50, 0x61, 0x74,
	0x63, 0x68, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x14, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x50, 0x61, 0x74, 0x63, 0x68, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73,
	0x22, 0x4e, 0x0a, 0x08, 0x4b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1b, 0x0a, 0x09,
	0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x6d, 0x65, 0x72,
	0x67, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79,
	0x22, 0x9b, 0x01, 0x0a, 0x10, 0x4b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x6c,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x25, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0c,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0a,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x50, 0x61, 0x74, 0x68, 0x12, 0x35, 0x0a,
	0x09, 0x6b, 0x38, 0x73, 0x5f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72,
	0x64, 0x2e, 0x4b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x08, 0x6b, 0x38, 0x73, 0x50,
	0x61, 0x74, 0x63, 0x68, 0x42, 0x08, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2a, 0x64,
	0x0a, 0x0a, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x54, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x54, 0x5f,
	0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x54, 0x5f,
	0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x54, 0x5f, 0x53,
	0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x54, 0x5f, 0x42, 0x4f,
	0x4f, 0x4c, 0x45, 0x41, 0x4e, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x54, 0x5f, 0x44, 0x41,
	0x54, 0x45, 0x10, 0x05, 0x2a, 0x90, 0x01, 0x0a, 0x0c, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x46,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x46, 0x5f, 0x4e, 0x4f, 0x4e, 0x45,
	0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x46, 0x5f, 0x49, 0x4e, 0x54, 0x33, 0x32, 0x10, 0x01,
	0x12, 0x0c, 0x0a, 0x08, 0x43, 0x46, 0x5f, 0x49, 0x4e, 0x54, 0x36, 0x34, 0x10, 0x02, 0x12, 0x0c,
	0x0a, 0x08, 0x43, 0x46, 0x5f, 0x46, 0x4c, 0x4f, 0x41, 0x54, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09,
	0x43, 0x46, 0x5f, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x43,
	0x46, 0x5f, 0x42, 0x59, 0x54, 0x45, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x46, 0x5f, 0x44,
	0x41, 0x54, 0x45, 0x10, 0x06, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x46, 0x5f, 0x44, 0x41, 0x54, 0x45,
	0x54, 0x49, 0x4d, 0x45, 0x10, 0x07, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x46, 0x5f, 0x50, 0x41, 0x53,
	0x53, 0x57, 0x4f, 0x52, 0x44, 0x10, 0x08, 0x2a, 0x2e, 0x0a, 0x09, 0x53, 0x63, 0x6f, 0x70, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x53,
	0x50, 0x41, 0x43, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x5f, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x10, 0x01, 0x2a, 0x2a, 0x0a, 0x0a, 0x53, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x53, 0x69, 0x64, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x53, 0x5f, 0x43, 0x4f, 0x4d, 0x4d,
	0x4f, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x53, 0x5f, 0x43, 0x4c, 0x49, 0x45, 0x4e,
	0x54, 0x10, 0x01, 0x3a, 0x53, 0x0a, 0x07, 0x6b, 0x38, 0x73, 0x5f, 0x63, 0x72, 0x64, 0x12, 0x1f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x85, 0xd5, 0xff, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x4b, 0x38, 0x73, 0x43, 0x52, 0x44,
	0x52, 0x06, 0x6b, 0x38, 0x73, 0x43, 0x72, 0x64, 0x3a, 0x57, 0x0a, 0x09, 0x6b, 0x38, 0x73, 0x5f,
	0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x86, 0xd5, 0xff, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x4b,
	0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x08, 0x6b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63,
	0x68, 0x3a, 0x54, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x1d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x87, 0xd5, 0xff, 0x22, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e,
	0x5f, 0x63, 0x72, 0x64, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x53, 0x69, 0x64, 0x65, 0x52,
	0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x42, 0x66, 0x0a, 0x23, 0x6c, 0x69, 0x62, 0x72, 0x61,
	0x72, 0x79, 0x2e, 0x67, 0x6f, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x03,
	0x43, 0x52, 0x44, 0x50, 0x01, 0x5a, 0x38, 0x61, 0x2e, 0x79, 0x61, 0x6e, 0x64, 0x65, 0x78, 0x2d,
	0x74, 0x65, 0x61, 0x6d, 0x2e, 0x72, 0x75, 0x2f, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x79, 0x2f,
	0x67, 0x6f, 0x2f, 0x6b, 0x38, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65,
	0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x63, 0x72, 0x64, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescOnce sync.Once
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescData = file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDesc
)

func file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescGZIP() []byte {
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescOnce.Do(func() {
		file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescData = protoimpl.X.CompressGZIP(file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescData)
	})
	return file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDescData
}

var file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_library_go_k8s_protoc_gen_crd_proto_crd_proto_goTypes = []interface{}{
	(ColumnType)(0),                     // 0: protoc_gen_crd.ColumnType
	(ColumnFormat)(0),                   // 1: protoc_gen_crd.ColumnFormat
	(ScopeType)(0),                      // 2: protoc_gen_crd.ScopeType
	(SchemaSide)(0),                     // 3: protoc_gen_crd.SchemaSide
	(*PrinterColumn)(nil),               // 4: protoc_gen_crd.PrinterColumn
	(*K8SCRD)(nil),                      // 5: protoc_gen_crd.K8sCRD
	(*K8SPatch)(nil),                    // 6: protoc_gen_crd.K8sPatch
	(*K8SPatchSelector)(nil),            // 7: protoc_gen_crd.K8sPatchSelector
	(*descriptorpb.MessageOptions)(nil), // 8: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 9: google.protobuf.FieldOptions
}
var file_library_go_k8s_protoc_gen_crd_proto_crd_proto_depIdxs = []int32{
	0,  // 0: protoc_gen_crd.PrinterColumn.type:type_name -> protoc_gen_crd.ColumnType
	1,  // 1: protoc_gen_crd.PrinterColumn.format:type_name -> protoc_gen_crd.ColumnFormat
	4,  // 2: protoc_gen_crd.K8sCRD.additional_columns:type_name -> protoc_gen_crd.PrinterColumn
	2,  // 3: protoc_gen_crd.K8sCRD.scope:type_name -> protoc_gen_crd.ScopeType
	7,  // 4: protoc_gen_crd.K8sCRD.field_patch_strategies:type_name -> protoc_gen_crd.K8sPatchSelector
	6,  // 5: protoc_gen_crd.K8sPatchSelector.k8s_patch:type_name -> protoc_gen_crd.K8sPatch
	8,  // 6: protoc_gen_crd.k8s_crd:extendee -> google.protobuf.MessageOptions
	9,  // 7: protoc_gen_crd.k8s_patch:extendee -> google.protobuf.FieldOptions
	9,  // 8: protoc_gen_crd.schema:extendee -> google.protobuf.FieldOptions
	5,  // 9: protoc_gen_crd.k8s_crd:type_name -> protoc_gen_crd.K8sCRD
	6,  // 10: protoc_gen_crd.k8s_patch:type_name -> protoc_gen_crd.K8sPatch
	3,  // 11: protoc_gen_crd.schema:type_name -> protoc_gen_crd.SchemaSide
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	9,  // [9:12] is the sub-list for extension type_name
	6,  // [6:9] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_library_go_k8s_protoc_gen_crd_proto_crd_proto_init() }
func file_library_go_k8s_protoc_gen_crd_proto_crd_proto_init() {
	if File_library_go_k8s_protoc_gen_crd_proto_crd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrinterColumn); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*K8SCRD); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*K8SPatch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*K8SPatchSelector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*K8SPatchSelector_ProtobufType)(nil),
		(*K8SPatchSelector_FieldPath)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   4,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_library_go_k8s_protoc_gen_crd_proto_crd_proto_goTypes,
		DependencyIndexes: file_library_go_k8s_protoc_gen_crd_proto_crd_proto_depIdxs,
		EnumInfos:         file_library_go_k8s_protoc_gen_crd_proto_crd_proto_enumTypes,
		MessageInfos:      file_library_go_k8s_protoc_gen_crd_proto_crd_proto_msgTypes,
		ExtensionInfos:    file_library_go_k8s_protoc_gen_crd_proto_crd_proto_extTypes,
	}.Build()
	File_library_go_k8s_protoc_gen_crd_proto_crd_proto = out.File
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_rawDesc = nil
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_goTypes = nil
	file_library_go_k8s_protoc_gen_crd_proto_crd_proto_depIdxs = nil
}
