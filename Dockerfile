FROM golang:alpine AS dist
WORKDIR /build
COPY . .
# 增加 CGO_ENABLED=0 确保静态编译，去除 /build/build 这种绕圈子的相对路径
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o yapfa ./cmd/main.go

FROM alpine AS production
WORKDIR /app
# 从编译阶段拷贝文件
COPY --from=dist /build/yapfa /app/yapfa

EXPOSE 8080
# 【关键修改】因为已经在 /app 下了，直接执行 ./yapfa 即可
CMD ["./yapfa"]
