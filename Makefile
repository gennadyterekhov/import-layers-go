
build:
	go build -o=analyzers/ilg cmd/main.go

build_all:
	GOOS=linux GOARCH=amd64 go build  -o ./dist/import_layers_linux_amd64  ./cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./dist/import_layers_darwin_amd64 ./cmd/main.go
	GOOS=darwin GOARCH=arm64 go build -o ./dist/import_layers_darwin_arm64 ./cmd/main.go
	GOOS=windows GOARCH=amd64 go build -o ./dist/import_layers_windows_amd64.exe ./cmd/main.go
