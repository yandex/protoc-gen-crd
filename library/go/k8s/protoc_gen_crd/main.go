package main

import (
	"flag"
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/pkg/gen"
)

var flags flag.FlagSet

func setBoolVar(dst **bool) func(val string) error {
	return func(s string) error {
		val, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		*dst = &val
		return nil
	}
}

func main() {
	plugin := gen.Plugin{}

	flags.BoolFunc("client-schema", "", setBoolVar(&plugin.IsClientSchema))
	flags.BoolFunc("schemaless", "", setBoolVar(&plugin.IsSchemaless))
	flags.BoolFunc("strict-schema", "", setBoolVar(&plugin.IsStrictSchema))
	flags.BoolFunc("generate-merge-keys", "", setBoolVar(&plugin.IsGeneratingMergeKeysEnabled))
	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(plugin.Run)
}
