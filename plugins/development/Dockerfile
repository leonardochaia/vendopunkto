FROM golang:alpine AS go-builder

# install dependencies
RUN apk add --no-cache alpine-sdk

WORKDIR /build

# install go packages in separate layer
COPY go.mod go.sum /build/
RUN go mod download

COPY . /build/
# build dev plugins
RUN go build -o ./vendopunkto-dev ./plugins/development/main.go

# dev plugins runtime base layer
FROM alpine as dev-plugins

WORKDIR /vendopunkto
RUN adduser -S -D -H -h /vendopunkto vp-user \
    && chown vp-user /vendopunkto \
    && chmod 755 /vendopunkto

COPY --from=go-builder /build/vendopunkto-dev /vendopunkto/

USER vp-user
EXPOSE 4000
ENTRYPOINT [ "./vendopunkto-dev" ]
CMD [""]