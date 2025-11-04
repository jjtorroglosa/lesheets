GO := GOAMD64=v2 GOARM64=v8.0 CGO_ENABLED=1 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
TMPL_FILES := $(shell find internal/views -name "*.html")
NASHEETS := ./nasheets -d output
MAIN := main.go
WASM_MAIN := cmd/wasm/main.go
ENTR:= entr -n

.PHONY: livereload
livereload:
	@echo livereload
	node ./node_modules/.bin/live-server --host=mini.home --port=3232 --open=output

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
		"make livereload" \
		"make watch" \
		"make watch-wasm" \
		"make watch-nns" \
		"make watch-js"


.PHONY: test
test:
	$(GO) test ./...


.PHONY: watch
watch:
	@echo watch
	ls internal/views/*.html css/styles.css $(GO_FILES) | \
				$(ENTR) -s "make nasheets ; ls examples/*.nns |xargs -I @ bin/update.sh @ "

watch-wasm:
	@echo watch-wasm
	ls $(TMPL_FILES) css/styles.css internal/*.go cmd/wasm/*.go | \
			$(ENTR) -s "make wasm"

.PHONY: watch-js
watch-js:
	mkdir -p build
	ls js/*.js vendorjs/*.js | $(ENTR) -a cp fonts/*.woff2 fonts/*.ttf js/*.js vendorjs/*.js build/

.PHONY: watch-nns
watch-nns:
	@echo watch-html
	ls examples/*.nns | $(ENTR) $(NASHEETS) html /_

nasheets: $(GO_FILES) $(TMPL_FILES) build/compiled.css
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)

output/%.html: examples/%.nns
	$(NASHEETS) html $<

.PHONY: wasm
wasm: output/wasm.wasm
output/wasm.wasm: output/wasm_exec.js $(GO_FILES) $(TMPL_FILES)
	GOOS=js GOARCH=wasm TG_CACHE=~/.tinygo-cache tinygo build -no-debug -opt=1 -o $@ $(WASM_MAIN)
	cp js/index.js output/index.js
	cp $@ build/wasm.wasm

output/wasm_exec.js:
	cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js $@

.PHONY: clean
clean:
	rm output/* build/*
