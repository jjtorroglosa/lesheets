GO := GOAMD64=v2 GOARM64=v8.0 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
JS_FILES := $(wildcard js/*.js vendorjs/*.js)
TMPL_FILES := $(wildcard internal/views/*.html)
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
		"make watch-wasm" \
		"make watch" \
		"make watch-js"


.PHONY: test
test:
	$(GO) test ./...

IN ?= examples/*.nns

.PHONY: watch
watch:
	@echo watch
	ls internal/views/*.html build/compiled.css build/*.js $(GO_FILES) build/*.wasm | \
				$(ENTR) -r -s "make nasheets ; $(NASHEETS) watch $(IN)"
.PHONY: run

run:
	$(NASHEETS) watch *.nns

.PHONY: watch-js
watch-js:
	mkdir -p build
	ls $(JS_FILES) | $(ENTR) -a -s "make js"

.PHONY: js
js:
	@echo build-js
	cp fonts/*.woff2 fonts/*.ttf js/*.js vendorjs/*.js build/

nasheets: $(GO_FILES) $(TMPL_FILES) build/compiled.css build/*.js
	@echo build-exec
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)

output/%.html: examples/%.nns
	$(NASHEETS) html $<


watch-wasm:
	@echo watch-wasm
	ls $(TMPL_FILES) build/compiled.css $(GO_FILES) $(wildcard build/*.js) | \
			$(ENTR) -a -s "make wasm"

.PHONY: wasm
wasm: build/wasm.wasm
build/wasm.wasm: $(GO_FILES) $(TMPL_FILES) $(wildcard build/*.js)
	#GOOS=js GOARCH=wasm GOTRACEBACK=all TG_CACHE=~/.tinygo-cache tinygo build -no-debug -opt=1 -o $@ $(WASM_MAIN)
	#cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js $@
	@echo build-wasm
	GOOS=js GOARCH=wasm GOTRACEBACK=all go build -ldflags="-s -w" -o $@ $(WASM_MAIN)
	cp vendorjs/wasm_exec_go.js build/wasm_exec.js
	cp js/index.js build/index.js

.PHONY: clean
clean:
	rm output/* build/*
