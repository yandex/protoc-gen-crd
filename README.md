protoc-gen-crd
==============

protoc-gen-crd is a protobuf compiler plugin to generate Kubernetes YAML spec for CRD from protobuf definition.

It is used to simplify creating and maintaing Kubernetes CRD schemas. Key features:

- Human-readable and human-writable. Protobuf is user friendly, K8s flavored OpenAPI is definitely not.
- Composable. K8s flavored OpenAPI denies even `$ref` between document parts, let alone references between multiple schemas.
  Protobuf allows you to split schema into multiple files, import external schema parts, and compose them as you wish.
- Schema versioning replaced with protobuf backward compatibility guidelines. Forget about dealing with multiple structures
  and conversion functions, use single protobuf message.
- Two modes: permissive server one and strict client one with support of client-only fields.
- The same protobuf files may be used to generate code in any language.

Just annotate single message in a proto-file with our special annotation, and compile it using this plugin.

Requirements
------------

- [Go](https://golang.org/doc/install) 1.22 (to build the plugin)
- [Protobuf compiler](https://github.com/protocolbuffers/protobuf/releases/latest)

Installation
------------

To install plugin, run the following command:

```sh
$ go install github.com/yandex/protoc-gen-crd/cmd/protoc-gen-crd@latest
```

Usage
-----

First, you need to add protobuf file  `proto/mycrd.proto` with the CRD annotations into your project:

```protobuf
import "github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/proto/crd.proto";

option go_package = "example.com/m/proto";

message Spec {
    // your fields
}

message Status {
    // your fields
}

message MyCrdKind {
    option (protoc_gen_crd.k8s_crd) = {
           api_group: "my-api.my-company.org",
           kind: "MyCrdKind",
           plural: "mycrdkinds",
           singular: "mycrdkind",
           short_names: ["mck", "mycrd"],
           categories: ["coolstuff", "my-company"],
           additional_columns: [
               {
                   name: "Gen",
                   type: CT_INTEGER,
                   description: "Object generation",
                   json_path: ".metadata.generation",
               },
           ]
    };

    Spec spec = 1;
    Status status = 2;
}

```

Option fields `api_group`, `kind`, `plural`, and `singular` are mandatory. Other ones can be omitted.
For more information about supported fields see comments in `library/go/k8s/protoc_gen_crd/proto/crd.proto`.

Then compile proto file into YAML using installed protoc-gen-crd plugin:

```sh
$ protoc -I=proto -I=vendor --crd_out=paths=source_relative:./proto mycrd.proto
```

This will give you `proto/mycrd.crd.yaml`, which can be put with `kubectl apply -f proto/mycrd.crd.yaml`.

Also see full example with protobuf and Makefile in the `example/` subdirectory.

Kustomize patch hints
---------------------

Optionally you can specify kustomize patch parameters via special annotation:

```protobuf
message Spec {
    repeated MyField my_fields = 1 [(protoc_gen_crd.k8s_patch) = {
        merge_key: "some_key",
        merge_strategy: "merge",
    }];
}
```

The other way to specify patch parameters is to specify them as part of protoc_gen_crd.k8s_crd option:

```protobuf
message MyCrdKind {
    option (protoc_gen_crd.k8s_crd) = {
        api_group: "my-api.my-company.org",
        kind: "MyCrdKind",
        /* ... */
        field_patch_strategies: [
            {
                // select some single field inside the MyCrdKind message.
                field_path: "spec.my_fields",
                k8s_patch: {
                    merge_key: "some_key",
                    merge_strategy: "merge",
                }
            },
            {
                // apply patch to all fields of the specified type.
                protobuf_type: "MyField",
                k8s_patch: {
                    merge_key: "some_key",
                    merge_strategy: "merge",
                }
            }
        ]
    };
};
```

When some patch parameters conflict, precedence is the following: field-specific patch > type-specific patch > field annotation.

Second way is intended to use for external APIs that are imported into your own proto message, and cannot be modified.

To build client version of the schema, add option `--crd_opt=client-schema=true` into your `protoc` command invocation.

Also see full example with protobuf and Makefile in the `example/` subdirectory.

Client schema and fields
------------------------

By default CRDs are compiled in server mode, but you can enable client mode with `--crd_opt=client-schema=true` option.

Key differces between modes:

1. Server mode is permissive and allows unknown object fields (`additionalProperties: true` in json schema)
   to enable protobuf compatibility with future schema versions. Client mode is strict because client always operates
   with known protobuf version.
2. Client mode allows to mark some fields as client-only, so that they would be processed on the client, but should not exist on server.

Schemaless CRD
------------------------

If you need to validate a custom resource with your own tools, you can replace CRD schema on the server by opaque object. In this case, K8S will save all unknown fields inside "spec" and "status", but will still validate common fields required for correct k8s operations (metadata, kind, apiVersion).

To generate such a scheme, add the option `--crd_opt=schemaless=true`

Known caveats
-------------

### oneOf

Protobuf `oneof` and OpenAPI `oneOf` have different semantics.
In protobuf `oneof` denotes just a single group of alternatives, and you may have as many groups as you wish.
In OpenAPI it is a single property of the message that should list all possible combinations (Cartesian product).

Hence, the following protobuf message:
```protobuf
message M {
    oneof group1 {
        string s1 = 1;
        string s2 = 2;
    }
    oneof group2 {
        string s3 = 3;
        string s4 = 4;
    }
}
```

Would become the following set of alternatives in OpenAPI:
```yaml
oneOf:
  - properties:
      s1: ...
      s3: ...
  - properties:
      s1: ...
      s4: ...
  - properties:
      s2: ...
      s3: ...
  - properties:
      s2: ...
      s4: ...
  - properties:
      s1: ...
  - properties:
      s2: ...
  - properties:
      s3: ...
  - properties:
      s4: ...
  - properties: {}
```

And on top of that Kubernetes requires for all properties inside `oneOf` [to be repeated][structural schema] outside of it.

Since this would effectively lead to [combinatorial explosion][combinatorial explosion], protobuf `oneof`s are NOT marked
in CRD in any way, and just rendered like a set of optional fields. Currently one should validate `oneof`s with validation
webhook or something similar: just parse the message as protobuf object and catch errors if any.

### Recursive structures

Protobuf allows structures to reference themselves, effectively making recursive structure, e.g.:

```protobuf
message M {
    M inner = 1;
}
```

While this also is possible in vanilla OpenAPI, kubernetes implementation [forbids][schema validation] any definition `$ref`s.

To deal with this limitation at the points of recursion we mark nested objects as opaque with unknown structure:
```yaml
properties:
  inner:
    type: object
    x-kubernetes-preserve-unknown-fields: true
```

Again, you are recommended to use validation webhook which would parse your object as protobuf message to check its schema.

[structural schema]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#specifying-a-structural-schema
[schema validation]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation
[combinatorial explosion]: https://en.wikipedia.org/wiki/Combinatorial_explosion
