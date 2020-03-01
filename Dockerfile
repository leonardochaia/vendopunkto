FROM golang:alpine AS go-builder

# install dependencies
RUN apk add --no-cache alpine-sdk

WORKDIR /build

# install go packages in separate layer
COPY go.mod go.sum /build/
RUN go mod download

COPY . /build/
# build vendopunkto and dev plugins
RUN make
RUN go build -o ./vendopunkto-dev ./plugins/development/main.go

# new base layer for nodejs builds
FROM node:stretch AS node-builder

WORKDIR /build/spa

# install node packages on separate layer
COPY ./spa/package.json ./spa/package-lock.json /build/spa/
RUN npm install

COPY ./spa /build/spa

# build angular apps
RUN npm run build shared --prod \
    && npm run build vendopunkto --prod \
    && npm run build admin --prod


# vendopunkto runtime base layer
FROM alpine as vendopunkto

WORKDIR /vendopunkto
RUN adduser -S -D -H -h /vendopunkto vp-user \
    && chown vp-user /vendopunkto \
    && chmod 755 /vendopunkto

COPY --from=go-builder /build/vendopunkto-server /vendopunkto/
COPY --from=node-builder /build/spa/dist /vendopunkto/spa/dist/

USER vp-user
EXPOSE 9080 8080
ENTRYPOINT [ "./vendopunkto-server" ]
CMD ["api"]
