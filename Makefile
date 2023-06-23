.PHONY: install-dev dev

install-dev:
	go install github.com/cespare/reflex@latest

dev:
	env $$(cat .env|xargs) reflex -d none -c reflex.conf
