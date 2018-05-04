.PHONY: all

run:
	go-assets-builder --package=main public/ > assets.go
	go build -o ggg.bin && ./ggg.bin

assets:
	go-assets-builder --package=main public/ > assets.go

provisioning:
	bundle exec itamae ssh -h ggg01 -j provision/nodes/webapp_production.json -u a-know provision/provisioning.rb

deploy:
	go-assets-builder --package=main public/ > assets.go
	GOOS=linux GOARCH=amd64 go build -o ggg.bin
	rsync -a --backup-dir=./.rsync_backup/$(LANG=C date +%Y%m%d%H%M%S) -e ssh ./* ggg01:/var/www/grass-graph/app
	ssh ggg01 sudo systemctl daemon-reload
	ssh ggg01 sudo systemctl restart grass-graph