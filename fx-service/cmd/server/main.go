package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"fxservice/internal/config"
	"fxservice/internal/handler"
	"fxservice/internal/middleware"
	"fxservice/internal/service"
	pb "fxservice/rpc/fxservice"
)

func main() {
	fxService := service.NewFXService()
	twirpHandler := handler.NewTwirpHandler(fxService)
	twirpServer := pb.NewFXServiceServer(twirpHandler)

	// Create rate limiter: 10 requests per minute per IP
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)

	mux := http.NewServeMux()
	mux.Handle(pb.FXServicePathPrefix,
		middleware.RateLimitMiddleware(rateLimiter)(
			middleware.ChaosMiddleware(
				handler.LoggingMiddleware(twirpServer))))

	port := ":" + config.GetPort()

	printBanner(port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func printBanner(port string) {
	banner := `
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║   ███████╗██╗  ██╗    ███████╗███████╗██████╗ ██╗   ██╗     ║
║   ██╔════╝╚██╗██╔╝    ██╔════╝██╔════╝██╔══██╗██║   ██║     ║
║   █████╗   ╚███╔╝     ███████╗█████╗  ██████╔╝██║   ██║     ║
║   ██╔══╝   ██╔██╗     ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝     ║
║   ██║     ██╔╝ ██╗    ███████║███████╗██║  ██║ ╚████╔╝      ║
║   ╚═╝     ╚═╝  ╚═╝    ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝       ║
║                                                               ║
║              Foreign Exchange Rate Service                    ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝

Server running on http://localhost%s

`
	fmt.Printf(banner, port)
}
