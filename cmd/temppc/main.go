package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Dionizio8/go-temppc-dist/configs"
	"github.com/Dionizio8/go-temppc-dist/internal/infra/client"
	"github.com/Dionizio8/go-temppc-dist/internal/infra/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := configs.LoadConfig(".")
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

	addressRepository := client.NewAddressRepository(cfg.ViaCEPClientURL)
	temperatureRepository := client.NewTemperatureRepository(cfg.WeatherAPIClientURL, cfg.WeatherAPIClientAPIKey)
	temperatureHandler := web.NewWebTemperatureHandler(addressRepository, temperatureRepository, tracer)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/metrics", promhttp.Handler())
	r.Route("/temppc", func(r chi.Router) {
		r.Get("/temperature/{zipCode}", temperatureHandler.GetTemperature)
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
