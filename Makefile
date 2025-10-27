
.PHONY: livereload
livereload:
	node ./node_modules/.bin/live-server --host=mini.home --open=views

.PHONY: tailwind
tailwind:
	yarn tailwindcss --input views/styles.css --output compiled.css --watch

.PHONY: dev
dev:
	ls *.nns tmpl.html internal/*.go cmd/nasheets/*.go | entr go run cmd/nasheets/main.go
