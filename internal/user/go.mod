module github.com/kimnattanan/graph-rag-service/internal/user

go 1.25.0

require (
	github.com/go-chi/chi/v5 v5.3.1
	github.com/go-chi/render v1.0.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.11.0
	github.com/kimnattanan/graph-rag-service/internal/common v0.0.0-00010101000000-000000000000
	github.com/oapi-codegen/runtime v1.6.0
	github.com/sirupsen/logrus v1.9.4
	golang.org/x/crypto v0.53.0
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/go-chi/cors v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/grpc v1.82.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/kimnattanan/graph-rag-service/internal/common => ../common/
