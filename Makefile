GO := GOAMD64=v2 GOARM64=v8.0 GOPROXY=https://athens.jjtorroglosa.com,direct go
GO_FLAGS = -ldflags -w

GO_FILES := $(shell find . -name "*.go")
ABC2SVG := vendorjs/abc2svg-compiled.js
JS_INPUT_FILES = $(wildcard js/* vendorjs/*.js)
JS_OUTPUT_FILES := build/editor.js build/sheet.js build/livereload.js build/
TMPL_FILES = $(wildcard internal/views/*.templ)
LESHEETS := ./build/lesheets
MAIN := main.go
WASM_MAIN := cmd/wasm/main.go
ENTR:= entr
IN ?= examples/*.nns

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

build/compiled.css: css/styles.css node_modules/.bin/tailwindcss
	NODE_ENV=production yarn tailwindcss --input css/styles.css --output build/compiled.css --minify

node_modules/.bin/tailwindcss:
	yarn install --immutable --check-cache

.PHONY: editor
editor:
	yarn tailwindcss --input css/styles.css --output build/compiled.css
	make js wasm $(LESHEETS) && $(LESHEETS) editor

.PHONY: dev
dev:
	make css templ wasm js build $(LESHEETS)
	yarn run concurrently \
		"make watch-css" \
		"make watch-templ" \
		"make watch-wasm" \
		"make watch-js" \
		"make watch-build" \
		"make watch-run"

.PHONY: prod
prod: css templ wasm js build $(LESHEETS) html compress

.PHONY: test
test:
	$(GO) test ./...


.PHONY: watch-templ
watch-templ:
	@echo watch-templ
	templ generate -watch

.PHONY: templ
templ:
	templ generate


.PHONY: watch-build
watch-build:
	@echo watch-build
	ls build/wasm.wasm | \
				entr -a make $(LESHEETS)


.PHONY: watch-run
watch-run:
	ls $(LESHEETS) | entr -r -s "$(LESHEETS) watch $(IN)"

.PHONY: run
run: $(LESHEETS)
	$(LESHEETS) watch *.nns

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
	touch $@

$(LESHEETS): templ $(GO_FILES) $(TMPL_FILES) build/compiled.css $(JS_OUTPUT_FILES) build/wasm.wasm
	@echo build-exec
	$(GO) build $(GO_FLAGS) -o $@ $(MAIN)


.PHONY: html
html:
	$(LESHEETS) html examples/*.nns

output/%.html: examples/%.nns
	$(LESHEETS) html $<


watch-wasm:
	@echo watch-wasm
	ls $(TMPL_FILES) build/compiled.css $(GO_FILES) | \
			$(ENTR) -a -s "make build/wasm.wasm"

.PHONY: wasm
wasm: build/wasm.wasm
build/wasm.wasm: $(GO_FILES) $(TMPL_FILES)
	@echo build-wasm
	GOOS=js GOARCH=wasm tinygo build -no-debug -opt=1 -o $@ $(WASM_MAIN)
	@#GOOS=js GOARCH=wasm GOTRACEBACK=all tinygo build -o $@ $(WASM_MAIN)
	@#cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js build/wasm_exec.js
	@#GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o build/unoptimized.wasm $(WASM_MAIN)
	@#GOOS=js GOARCH=wasm GOTRACEBACK=all go build -o build/unoptimized.wasm $(WASM_MAIN)
	@#wasm-opt build/unoptimized.wasm -Oz -o $@
	@#cp build/unoptimized.wasm $@


.PHONY: compress
compress:
	gzip -f -k -9 output/*.{js,css,wasm,html}

.PHONY: clean
clean:
	rm -rf output/* build/* internal/views/*_templ.go

.PHONY: deploy
deploy: prod docker
	 bin/deploy.sh $(TAG)


.PHONY: docker
docker: build/lesheets.tgz
build/lesheets.tgz: Dockerfile $(JS_OUTPUT_FILES) output/*
	docker buildx build --platform linux/amd64 -t $(IMAGE_NAME):$(TAG) .
	docker save $(IMAGE_NAME):$(TAG) | gzip > $@
	docker load -i $@

define confirm
	@read -p "$(1)? (y/N) " ans; \
	[ "$$ans" = "y" ]
endef

.PHONY: readme
readme:
	mdsh README.template.md README.md

