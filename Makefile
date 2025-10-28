
GO := GOAMD64=v2 GOARM64=v8.0 CGO_ENABLED=1 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
TMPL_FILES := $(shell find . -name "tmpl.html")
NASHEETS := ./nasheets -d output

.PHONY: livereload
livereload:
	@echo livereload
	node ./node_modules/.bin/live-server --host=mini.home --port=3232 --open=output

.PHONY: tailwind
tailwind:
	@echo tailwind
	yarn tailwindcss --input views/styles.css --output output/compiled.css --watch

.PHONY: dev
dev:
	yarn run concurrently "make tailwind" "make livereload" "make watch" "make watch-html"

.PHONY: watch
watch:
	@echo watch
	ls views/*.html views/styles.css internal/*.go cmd/nasheets/*.go | \
			entr -n -s "make nasheets ; ls -1 *.nns |xargs -I @ bin/update.sh @ "

.PHONY: watch-html
watch-html:
	@echo watch-html
	ls *.nns | entr -a -n $(NASHEETS) /_

nasheets: $(GO_FILES) $(TMPL_FILES)
	$(GO) build $(GO_FLAGS) -o nasheets cmd/nasheets/main.go

%.html: %.nns
	$(NASHEETS) $< $@
