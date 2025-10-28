
.PHONY: livereload
livereload:
	node ./node_modules/.bin/live-server --host=mini.home --open=views

.PHONY: tailwind
tailwind:
	yarn tailwindcss --input views/styles.css --output views/compiled.css --watch

.PHONY: dev
dev:
	ls *.nns views/tmpl.html views/styles.css internal/*.go cmd/nasheets/*.go | \
		entr go run cmd/nasheets/main.go index.nns
