package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/chi-net/yapfa/internal/handler"
	"github.com/chi-net/yapfa/internal/service"
)

// Pod file agent is used for get file info, list, upload or download files.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	h := handler.New("/data")
	client := service.NewClient(h.Instance())
	client.Start(ctx)
	<-ctx.Done()
}
