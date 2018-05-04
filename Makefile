.PHONY: all

run:
	go-assets-builder --package=main public/ > assets.go
	go build -o ggg.bin && ./ggg.bin

assets:
	go-assets-builder --package=main public/ > assets.go

provisioning:
	bundle exec itamae ssh -h ggg01 -j provision/nodes/webapp_production.json -u a-know provision/provisioning.rb
