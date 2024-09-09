package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/pkg/gen"
)

var flags flag.FlagSet

func main() {
	isClientSchema := flags.Bool("client-schema", false, "")
	isSchemalessCrd := flags.Bool("schemaless", false, "")
	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	plugin := gen.Plugin{IsClientSchema: isClientSchema, IsSchemaless: isSchemalessCrd}
	opts.Run(plugin.Run)
}
