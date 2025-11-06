GO := GOAMD64=v2 GOARM64=v8.0 CGO_ENABLED=1 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
TMPL_FILES := $(shell find internal/views -name "*.html")
NASHEETS := ./nasheets
MAIN := main.go
WASM_MAIN := cmd/wasm/main.go
ENTR:= entr -n

.PHONY: tailwind
tailwind:
	@echo tailwind
	mkdir -p build
	yarn tailwindcss --input css/styles.css --output build/compiled.css --watch

build/compiled.css: css/styles.css
	yarn tailwindcss --input css/styles.css --output build/compiled.css

.PHONY: dev
dev:
	yarn run concurrently \
		"make tailwind" \
		"make watch" \
		"make watch-wasm" \
		"make watch-js"


.PHONY: test
test:
	$(GO) test ./...

IN ?= examples/*.nns

.PHONY: watch
watch:
	@echo watch
	ls internal/views/*.html css/styles.css build/*.js $(GO_FILES) | \
				$(ENTR) -r -s "make nasheets ; $(NASHEETS) watch $(IN)"
.PHONY: run

run:
	$(NASHEETS) watch *.nns

watch-wasm:
	@echo watch-wasm
	ls $(TMPL_FILES) css/styles.css internal/*.go cmd/wasm/*.go | \
			$(ENTR) -s "make wasm"

.PHONY: watch-js
watch-js:
	mkdir -p build
	ls js/*.js vendorjs/*.js | $(ENTR) -a cp fonts/*.woff2 fonts/*.ttf js/*.js vendorjs/*.js build/

nasheets: $(GO_FILES) $(TMPL_FILES) build/compiled.css build/*.js
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)

output/%.html: examples/%.nns
	$(NASHEETS) html $<

.PHONY: wasm
wasm: output/wasm.wasm
build/wasm.wasm: build/wasm_exec.js $(GO_FILES) $(TMPL_FILES)
	GOOS=js GOARCH=wasm TG_CACHE=~/.tinygo-cache tinygo build -no-debug -opt=1 -o $@ $(WASM_MAIN)
	cp js/index.js build/index.js

output/wasm_exec.js:
	cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js $@

.PHONY: clean
clean:
	rm output/* build/*
