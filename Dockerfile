FROM golang:alpine AS dist
COPY . /build
WORKDIR /build
RUN cd /build && go build -ldflags="-s -w" -trimpath -o ./build/yapfa ./cmd/main.go

FROM alpine AS production
WORKDIR /app
COPY --from=dist "/build/build/yapfa" /app/yapfa
CMD ["./app/yapfa"]

EXPOSE 8080