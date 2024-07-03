package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/pkg/gen"
)

var flags flag.FlagSet

func main() {
	isClientSchema := flags.Bool("client-schema", false, "")
	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	plugin := gen.Plugin{IsClientSchema: isClientSchema}
	opts.Run(plugin.Run)
}
