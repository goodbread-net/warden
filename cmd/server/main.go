package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/matthiase/warden"
	"github.com/matthiase/warden/config"
	"github.com/matthiase/warden/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	cfg := config.ReadEnv()

	app, err := warden.NewApplication(cfg)
	if err != nil {
		log.Fatalf("Could not set up application: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	//listener, err := net.Listen("tcp", addr)
	//if err != nil {
	//	log.Fatalf("Error occurred: %s", err.Error())
	//}

	//httpHandler := routes.NewHandler(app)
	//server := &http.Server{
	//	Handler: httpHandler,
	//}
	//go func() {
	//	server.Serve(listener)
	//}()

	//defer Stop(server)

	log.Printf("Starting server on %s", addr)
	http.ListenAndServe(addr, routes.NewHandler(app))

	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	//log.Println(fmt.Sprint(<-ch))
	//log.Println("Stopping API server")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
