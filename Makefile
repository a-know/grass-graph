.PHONY: all

run:
	go-assets-builder --package=main public/ > assets.go
	go build -o ggg.bin && ./ggg.bin

assets:
	go-assets-builder --package=main public/ > assets.go
	