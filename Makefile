.PHONY: all

VERSION = 100

deps:
	env GO111MODULE=on go mod download

deploy: 
	gcloud app deploy --project a-know-home --version ${VERSION}
