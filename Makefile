.PHONY: start
start:
	hugo serve --bind="0.0.0.0" --baseUrl=$(shell hostname)
