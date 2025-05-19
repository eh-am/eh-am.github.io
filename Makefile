.PHONY: serve
serve:
	hugo serve --bind="0.0.0.0"

.PHONY: serve-external
serve-external:
	# TODO: only works on mac
	hugo serve --bind="0.0.0.0" --baseURL=$(shell ifconfig | grep "inet " | grep -v 127.0.0.1 | cut -d\  -f2)

.PHONY: build
build:
	cd ./scripts/weekly-digest && go build .

.PHONY: digest
digest: build
	./scripts/weekly-digest/weekly-digest -out=./data/digest last week
	
.PHONY: digest-all
digest-all: build
	./scripts/weekly-digest/weekly-digest -out=./data/digest all
