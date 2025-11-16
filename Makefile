GO := GOAMD64=v2 GOARM64=v8.0 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
ABC2SVG := vendorjs/abc2svg-compiled.js
JS_INPUT_FILES = $(wildcard js/* vendorjs/*.js)
JS_OUTPUT_FILES := build/editor.js build/sheet.js build/livereload.js
TMPL_FILES = $(wildcard internal/views/*.html)
NASHEETS := ./nasheets
MAIN := main.go
WASM_MAIN := cmd/wasm/main.go
ENTR:= entr

TAG ?= $(shell date +'%Y%m%d.%H%M')
IMAGE_NAME := lesheets

.PHONY: watch-css
watch-css:
	@echo css
	mkdir -p build
	NODE_ENV=production yarn tailwindcss --input css/styles.css --output build/compiled.css --minify --watch

.PHONY: css
css:
	mkdir -p build
	NODE_ENV=production yarn tailwindcss --input css/styles.css --output build/compiled.css --minify

build/compiled.css: css/styles.css
	NODE_ENV=production yarn tailwindcss --input css/styles.css --output build/compiled.css --minify

.PHONY: editor
editor:
	yarn tailwindcss --input css/styles.css --output build/compiled.css
	make js wasm nasheets && ./nasheets editor

.PHONY: dev
dev:
	yarn run concurrently --kill-others-on-fail \
		"make watch-css" \
		"make watch-wasm" \
		"make watch-js" \
		"make watch-build" \
		"make watch-run"

.PHONY: prod
prod: css wasm js build nasheets html compress

.PHONY: test
test:
	$(GO) test ./...

IN ?= examples/*.nns

.PHONY: watch-build
watch-build:
	@echo watch-build
	ls $(TMPL_FILES) build/compiled.css $(JS_OUTPUT_FILES) $(GO_FILES) build/wasm.wasm | \
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
	node build.mjs --dev --watch

.PHONY: js
js: $(JS_OUTPUT_FILES)
$(JS_OUTPUT_FILES): $(JS_INPUT_FILES)
	@echo build-js
	cp fonts/*.woff2 fonts/*.ttf build/
	node build.mjs

nasheets: $(GO_FILES) $(TMPL_FILES) build/compiled.css $(JS_OUTPUT_FILES) build/wasm.wasm
	@echo build-exec
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)


.PHONY: html
html:
	$(NASHEETS) html examples/*.nns

output/%.html: examples/%.nns
	$(NASHEETS) html $<


watch-wasm:
	@echo watch-wasm
	ls $(TMPL_FILES) build/compiled.css $(GO_FILES) | \
			$(ENTR) -a -s "make build/wasm.wasm"

.PHONY: wasm
wasm: build/wasm.wasm
build/wasm.wasm: $(GO_FILES) $(TMPL_FILES)
	#GOOS=js GOARCH=wasm GOTRACEBACK=all TG_CACHE=~/.tinygo-cache tinygo build -no-debug -opt=1 -o build/unoptimized.wasm $(WASM_MAIN)
	#cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js $@
	@echo build-wasm
	GOOS=js GOARCH=wasm GOTRACEBACK=all go build -ldflags="-s -w" -o build/unoptimized.wasm $(WASM_MAIN)
	#GOOS=js GOARCH=wasm GOTRACEBACK=all go build -o build/unoptimized.wasm $(WASM_MAIN)
	wasm-opt build/unoptimized.wasm -Oz --enable-bulk-memory-opt -o $@
	#cp build/unoptimized.wasm $@


.PHONY: compress
compress:
	gzip -k -9 output/*.{js,css,wasm,html}

.PHONY: clean
clean:
	rm -rf output/* build/*

.PHONY: deploy
deploy: prod docker
	 bin/deploy.sh $(TAG)


.PHONY: docker
docker: build/lesheets.tgz
build/lesheets.tgz: Dockerfile $(JS_OUTPUT_FILES) output/*
	docker buildx build --platform linux/amd64 -t $(IMAGE_NAME):$(TAG) .
	docker save $(IMAGE_NAME):$(TAG) | gzip > $@
	docker load -i $@
