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
	github.com/sirupsen/logrus v1.9.4 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

replace github.com/kimnattanan/graph-rag-service/internal/common => ../common/
