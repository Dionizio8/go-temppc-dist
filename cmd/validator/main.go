package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Dionizio8/go-temppc-dist/configs"
	"github.com/Dionizio8/go-temppc-dist/internal/infra/client"
	"github.com/Dionizio8/go-temppc-dist/internal/infra/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := configs.LoadConfigValidator(".")
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdown, err := configs.InitProvider(ctx, cfg.OtelServiceName, cfg.OtelCollectorURL)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")

	temppcRepository := client.NewTemppcRepository(cfg.TemppcAPICleintURL)
	temperatureByClientHandler := web.NewWebTemperatureByClientHandler(temppcRepository, tracer)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/validator", func(r chi.Router) {
		r.Get("/temperature/{zipCode}", temperatureByClientHandler.GetTemperature)
	})

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- http.ListenAndServe(cfg.WebServerPort, r)
	}()

	select {
	case err = <-serverErr:
		log.Fatal(err)
		return
	case <-ctx.Done():
		stop()
	}
}
