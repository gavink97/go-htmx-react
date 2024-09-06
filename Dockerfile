# https://docs.docker.com/build/building/best-practices/#build-and-test-your-images-in-ci
# https://programmingwithwolfgang.com/fixing-broken-unit-test-execution-in-dockerfile/

ARG VERSION="0.1-beta"
ARG NODE_VERSION=22-alpine
ARG GO_VERSION=1.22-alpine

FROM node:$NODE_VERSION AS node
WORKDIR /app

COPY --link package.json ./
COPY --link ./public/css/globals.css ./public/css/globals.css
COPY --link ./build/esbuild.mjs ./build/esbuild.mjs
COPY --link ./build/sharp.js ./build/sharp.js
COPY --link ./public/images ./public/images
COPY --link ./Makefile ./
COPY --link ./src ./src
COPY --link ./tsconfig.json ./
COPY --link ./tailwind.config.js ./
COPY --link ./postcss.config.js ./
COPY --link ./babel.config.js ./
COPY --link ./internal ./internal

RUN : \
&& apk add --no-cache make \
&& npm install --verbose \
&& :


FROM node AS jest-test
RUN make jest


FROM node AS prebuild
WORKDIR /app

RUN : \
&& npx tailwindcss -i ./public/css/globals.css -o ./public/css/style.min.css --minify \
&& make eb \
&& make sharp \
&& :


FROM golang:$GO_VERSION AS go
WORKDIR /app
COPY go.mod go.sum ./

RUN : \
&& go mod download \
# && apk add --no-cache make \
&& apk add --no-cache make build-base \
&& :


FROM ghcr.io/a-h/templ:latest AS templ
COPY --chown=65532:65532 . /app
COPY --from=go /app /app
COPY --from=prebuild --chown=65532:65532 /app /app

WORKDIR /app


RUN ["templ", "generate"]


FROM go AS build
WORKDIR /app

COPY --from=templ /app /app

RUN : \
# && CGO_ENABLED=0 GOOS=linux go build -o /gavin-site ./cmd/gavin-site/main.go \
&& CGO_ENABLED=1 GOOS=linux go build -o /gavin-site ./cmd/gavin-site/main.go \
&& chmod +x /gavin-site \
# && adduser --disabled-password -u 10001 nonroot \
&& :


FROM build AS go-test
RUN make gotest


#FROM scratch AS deploy
FROM alpine AS deploy
WORKDIR /

COPY --from=build /gavin-site ./gavin-site
COPY --from=build /app/assets ./assets
COPY --from=build /app/public/css ./public/css
COPY --from=build /app/public/favicon.ico ./public/favicon.ico
COPY --from=build /app/bin/images ./public/images

#COPY --link --from=build /etc/passwd /etc/passwd
#COPY --chown=nonroot --from=build /app/bin/gavin-site .
#COPY --chown=nonroot --from=build /app/assets ./assets
#COPY --chown=nonroot --from=build /app/public/css ./public/css
#COPY --chown=nonroot --from=build /app/public/favicon.ico ./public/favicon.ico
#COPY --chown=nonroot --from=build /app/bin/images ./public/images

# USER nonroot
ENV env=prod

EXPOSE 8080

ENTRYPOINT ["./gavin-site"]

LABEL vendor=gavink \
      ink.gav.is-beta=True\
      ink.gav.is-production=True \
      ink.gav.version=$VERSION \
      ink.gav.release-date="2024-08-30"
