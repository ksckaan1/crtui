run:
	@go run ./cmd/crtui/

install:
	@go install ./cmd/crtui

build:
	@goreleaser build --snapshot --clean

release:
	@goreleaser release --clean
