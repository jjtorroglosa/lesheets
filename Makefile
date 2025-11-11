GO := GOAMD64=v2 GOARM64=v8.0 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
JS_FILES := $(wildcard js/*.js)
TMPL_FILES := $(wildcard internal/views/*.html)
NASHEETS := ./nasheets
MAIN := main.go
WASM_MAIN := cmd/wasm/main.go
ENTR:= entr

.PHONY: tailwind
tailwind:
	@echo tailwind
	mkdir -p build
	yarn tailwindcss --input css/styles.css --output build/compiled.css --watch

build/compiled.css: css/styles.css
	yarn tailwindcss --input css/styles.css --output build/compiled.css

editor:
	yarn tailwindcss --input css/styles.css --output build/compiled.css
	make js wasm nasheets && ./nasheets editor
.PHONY: dev
dev:
	yarn run concurrently \
		"make tailwind" \
		"make watch-wasm" \
		"make watch-js" \
		"make watch-build" \
		"make watch-run"


.PHONY: test
test:
	$(GO) test ./...

IN ?= examples/*.nns

.PHONY: watch-build
watch-build:
	@echo watch-build
	ls $(TMPL_FILES) build/compiled.css build/*.js $(GO_FILES) build/*.wasm | \
				entr -a make nasheets


.PHONY: watch-run
watch-run:
	ls nasheets | entr -r -s "$(NASHEETS) watch $(IN)"

.PHONY: run
run:
	$(NASHEETS) watch *.nns

.PHONY: watch-js
watch-js:
	mkdir -p build
	ls $(JS_FILES) | $(ENTR) -a -s "make js"

ABC2SVG :=vendorjs/abc2svg-compiled.js
.PHONY: js
js: build/bundle.js
build/bundle.js: $(ABC2SVG)
	@echo build-js
	cp fonts/*.woff2 fonts/*.ttf build/
	node build.mjs

$(ABC2SVG): vendorjs/abc2svg-1.js vendorjs/abc2svg-caller.js
	# Concat the vendored abc2svg code with the -caller, to export the RenderFunction
	cat vendorjs/abc2svg-1.js vendorjs/abc2svg-caller.js > $@


nasheets: $(GO_FILES) $(TMPL_FILES) build/compiled.css $(wildcard build/*.js)
	@echo build-exec
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)

output/%.html: examples/%.nns
	$(NASHEETS) html $<


watch-wasm:
	@echo watch-wasm
	ls $(TMPL_FILES) build/compiled.css $(GO_FILES) $(wildcard build/*.js) | \
			$(ENTR) -a -s "make build/wasm.wasm"

.PHONY: wasm
wasm: build/wasm.wasm
build/wasm.wasm: $(GO_FILES) $(TMPL_FILES) $(wildcard build/*.js)
	#GOOS=js GOARCH=wasm GOTRACEBACK=all TG_CACHE=~/.tinygo-cache tinygo build -no-debug -opt=1 -o $@ $(WASM_MAIN)
	#cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js $@
	@echo build-wasm
	GOOS=js GOARCH=wasm GOTRACEBACK=all go build -ldflags="-s -w" -o $@ $(WASM_MAIN)
	#GOOS=js GOARCH=wasm GOTRACEBACK=all go build -o $@ $(WASM_MAIN)
	cp vendorjs/wasm_exec_go.js build/wasm_exec.js
	cp js/index.js build/index.js

.PHONY: clean
clean:
	rm -rf output/* build/*
