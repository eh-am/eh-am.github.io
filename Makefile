.PHONY: serve
serve:
	hugo serve --bind="0.0.0.0"

.PHONY: serve-external
serve-external:
	# TODO: only works on mac
	hugo serve --bind="0.0.0.0" --baseURL=$(shell ifconfig | grep "inet " | grep -v 127.0.0.1 | cut -d\  -f2)

.PHONY: digest
digest:
	cd ./scripts/weekly-digest && go run . last week
	cp scripts/weekly-digest/output/* data/digest/
	
