package generate

//go:generate sh -c "mkdir -p ../../gens/server && go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.generate.yaml -o ../../gens/server/gen.go ../openapi.yaml"
