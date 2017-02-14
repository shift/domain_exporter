package = github.com/shift/domain_exporter

.PHONY: release

release:
	go get
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/domain_exporter-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/domain_exporter-linux-386 $(package)
	GOOS=linux GOARCH=arm go build -o release/domain_exporter-linux-arm $(package)
	GOOS=linux GOARCH=arm64 go build -o release/domain_exporter-linux-arm64 $(package)
	GOOS=linux GOARCH=mips64 go build -o release/domain_exporter-linux-mips64 $(package)
	GOOS=linux GOARCH=mips64le go build -o release/domain_exporter-linux-mip64le $(package)
	GOOS=darwin GOARCH=amd64 go build -o release/domain_exporter-darwin-amd64 $(package)
	GOOS=freebsd GOARCH=amd64 go build -o release/domain_exporter-freebsd-amd64 $(package)
	GOOS=netbsd GOARCH=amd64 go build -o release/domain_exporter-netbsd-amd64 $(package)
	GOOS=openbsd GOARCH=amd64 go build -o release/domain_exporter-openbsd-amd64 $(package)
	GOOS=windows GOARCH=amd64 go build -o release/domain_exporter-windows-amd64 $(package)
