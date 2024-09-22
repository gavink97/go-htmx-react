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
	npx playwright test tests

.PHONY: codegen
codegen:
	npx playwright codegen

.PHONY: report
report:
	npx playwright show-report

.PHONY: jest
jest:
	npx jest src

.PHONY: gotest
gotest:
	go test -race -v -timeout 30s ./...

.PHONY: dev
dev:
	go build -o .tmp/gavin-site ./cmd/gavin-site/main.go \
	&& air

.PHONY: test
test:
	make jest \
	&& make gotest \
	&& make playwright

.PHONY: build
build:
	docker build . -f Dockerfile

.PHONY: prod
prod:
	docker-compose -f docker-compose.yml up -d --build


.PHONY: goupdate
goupdate:
	go get -u ./...
