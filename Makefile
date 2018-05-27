.PHONY: all

run:
	go-assets-builder --package=main public/ > assets.go
	go build -o ggg.bin && ./ggg.bin

assets:
	go-assets-builder --package=main public/ > assets.go

provisioning:
	bundle exec itamae ssh -h ${TARGET} -j provision/nodes/${TARGET}.json -u a-know provision/provisioning.rb

deploy:
	go-assets-builder --package=main public/ > assets.go
	GOOS=linux GOARCH=amd64 go build -o ggg.bin
	rsync -a --backup-dir=./.rsync_backup/$(LANG=C date +%Y%m%d%H%M%S) -e ssh ./* ${TARGET}:/var/www/grass-graph/app
	ssh ${TARGET} sudo systemctl daemon-reload
	ssh ${TARGET} sudo systemctl restart grass-graph