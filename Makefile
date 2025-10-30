GO := GOAMD64=v2 GOARM64=v8.0 CGO_ENABLED=1 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
TMPL_FILES := $(shell find . -name "tmpl.html")
NASHEETS := ./nasheets -d output
MAIN := cmd/nasheets/main.go
WASM_MAIN := cmd/wasm/main.go
ENTR:= entr -n

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
	yarn run concurrently \
		"make tailwind" \
		"make livereload" \
		"make watch" \
		"make watch-wasm" \
		"make watch-nns" \
		"make watch-js"

.PHONY: watch
watch:
	@echo watch
	ls views/*.html views/styles.css internal/*.go cmd/nasheets/*.go | \
			$(ENTR) -s "make nasheets ; ls *.nns |xargs -I @ bin/update.sh @ "

watch-wasm:
	@echo watch-wasm
	ls views/*.html views/styles.css internal/*.go cmd/wasm/*.go | \
			$(ENTR) -s "make wasm"

.PHONY: watch-js
watch-js:
	ls views/*.js vendorjs/*.js | $(ENTR) -a cp views/*.ttf views/*.js vendorjs/*.js output/

.PHONY: watch-nns
watch-nns:
	@echo watch-html
	ls *.nns | $(ENTR) $(NASHEETS) -i /_ html

nasheets: $(GO_FILES) $(TMPL_FILES)
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)

output/%.html: %.nns
	$(NASHEETS) -i $< html

.PHONY: wasm
wasm: output/wasm.wasm
output/wasm.wasm: output/wasm_exec.js $(GO_FILES) $(TMPL_FILES)
	GOOS=js GOARCH=wasm TG_CACHE=~/.tinygo-cache tinygo build -no-debug -opt=1 -o $@ $(WASM_MAIN)
	cp views/index.js output/index.js

output/wasm_exec.js:
	cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js $@

.PHONY: clean
clean:
	rm output/*
