run:
	@go run ./cmd/crtui/

install:
	@go install ./cmd/crtui

build:
	@goreleaser release --snapshot --clean

release:
	@goreleaser release --clean
