module example.com/m

go 1.22.1

replace github.com/yandex/protoc-gen-crd v1.0.0 => ..

require github.com/yandex/protoc-gen-crd v1.0.0

require google.golang.org/protobuf v1.34.1 // indirect
