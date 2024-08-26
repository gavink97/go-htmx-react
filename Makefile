.PHONY: tw
tw:
	tailwindcss -i ./public/css/globals.css -o ./public/css/styles.css --watch

.PHONY: tb
tb:
	tailwindcss -i ./public/css/globals.css -o ./public/css/style.min.css --minify

.PHONY: tg
tg:
	templ generate

.PHONY: eb
eb:
	node build/esbuild.mjs -build

.PHONY: ew
ew:
	node build/esbuild.mjs

.PHONY: sharp
sharp:
	node build/sharp.js

.PHONY: playwright
playwright:
	npx playwright test tests/playwright

.PHONY: codegen
codegen:
	npx playwright codegen

.PHONY: report
report:
	npx playwright show-report

.PHONY: jest
jest:
	npx jest tests/jest

.PHONY: go
go:
	go test -race -v -timeout 30s ./...

.PHONY: dev
dev:
	go build -o .tmp/gavin-site ./cmd/gavin-site/main.go \
	&& air

.PHONY: test
test:
	make playwright \
	# make jest \
	make go

.PHONY: build
build:
	docker build . -f Dockerfile

.PHONY: prod
prod:
	docker-compose -f docker-compose.yml up -d --build
