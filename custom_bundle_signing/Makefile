build:
	@go build \
		-o bin/opa \
		./cmd/opa
	@./bin/opa build \
		--bundle ./policy \
		--output ./policy/bundle.tar.gz \
		--signing-key my-key \
		--signing-plugin custom
	@./bin/opa sign \
		--bundle \
		--signing-key my-key \
		--signing-plugin custom \
		./policy
	@diff -q .signatures.json .signatures.json.expected

fmt:
	@go fmt ./cmd/opa
	@go fmt ./internal
	@go fmt ./setup

run:
	@./bin/opa run \
		--bundle \
		--verification-key foo \
		--verification-key-id custom \
		--exclude-files-verify data.json \
		--exclude-files-verify policy/awesome.rego \
		--exclude-files-verify .manifest \
		--exclude-files-verify .signatures.json \
		./policy/bundle.tar.gz

