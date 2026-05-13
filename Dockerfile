FROM golang:alpine as dist
COPY . /build
WORKDIR /build
RUN "cd /build && go build -ldflags=\"-s -w\" -trimpath -o ./build/yapfa ./cmd/main.go"

FROM alpine as production
WORKDIR /app
COPY --from=dist "/build/cmd/build/yapfa" /app/yapfa
CMD ./app/yapfa

EXPOSE 8080