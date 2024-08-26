ARG NODE_VERSION=22-alpine
ARG GO_VERSION=1.22-alpine


FROM node:$NODE_VERSION AS prebuild

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

RUN apk add --no-cache make

RUN : \
&& npm install --verbose\
&& npx tailwindcss -i ./public/css/globals.css -o ./public/css/style.min.css --minify \
&& make eb \
&& make sharp \
&& :


FROM golang:$GO_VERSION AS fetch

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


FROM ghcr.io/a-h/templ:latest AS generate
COPY --chown=65532:65532 . /app
COPY --from=fetch /app /app
COPY --from=prebuild /app /app

WORKDIR /app
RUN ["templ", "generate"]


FROM golang:$GO_VERSION AS build
COPY --from=generate /app /app
ENV env=prod
ENV HOST=0.0.0.0

RUN apk add --no-cache make build-base

WORKDIR /app
RUN : \
# && CGO_ENABLED=0 GOOS=linux go build -o /gavin-site ./cmd/gavin-site/main.go \
&& CGO_ENABLED=1 GOOS=linux go build -o /gavin-site ./cmd/gavin-site/main.go \
&& chmod +x /gavin-site \
&& adduser --disabled-password -u 10001 nonroot \
&& :


# go not running tests
FROM build AS test
RUN make test


#FROM scratch AS deploy
FROM alpine
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
