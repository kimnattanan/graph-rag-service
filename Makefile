include .env
export

.PHONY: openapi
openapi: openapi_http openapi_js

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh user internal/user/ports ports
	@./scripts/openapi-http.sh knowledge internal/knowledge/ports ports
	@./scripts/openapi-http.sh conversation internal/conversation/ports ports

.PHONY: openapi_js
openapi_js:
	@./scripts/openapi-js.sh user
	@./scripts/openapi-js.sh knowledge
	@./scripts/openapi-js.sh conversation

.PHONY: proto
proto:
	@./scripts/proto.sh user
	@./scripts/proto.sh knowledge
	@./scripts/proto.sh conversation
