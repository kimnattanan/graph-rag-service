module github.com/kimnattanan/graph-rag-service/internal/knowledge

go 1.25.0

require (
	github.com/go-chi/chi/v5 v5.3.1
	github.com/kimnattanan/graph-rag-service/internal/common v0.0.0-00010101000000-000000000000
	github.com/oapi-codegen/runtime v1.6.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/go-chi/cors v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.42.1 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/term v0.38.0 // indirect
)

replace github.com/kimnattanan/graph-rag-service/internal/common => ../common/
