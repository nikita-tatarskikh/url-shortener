package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/configuration"
	logger2 "url-shortener/internal/logger"
	"url-shortener/internal/metrics/prometheus"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
	"url-shortener/internal/transport"
	"url-shortener/internal/transport/handlers"
	"url-shortener/pkg/redis"
	"url-shortener/pkg/zap"

	"github.com/oklog/run"
	"github.com/pkg/errors"
)

var ErrOsSignal = errors.New("got os signal")

const (
	GolangEnv = "GOLANG_ENVIRONMENT"
)

func main() {
	cfg, errAppConf := configuration.NewAppConfiguration(
		os.Getenv(GolangEnv), true)
	if errAppConf != nil {
		log.Fatal(errors.WithMessage(errAppConf, "app configuration provider"))
	}

	zapLogger, cleanupZapLogger, errZapLogger := zap.New(zap.Mode(cfg.ZapLoggerMode))
	if errZapLogger != nil {
		log.Fatal(errors.WithMessage(errZapLogger, "zap logger provider"))
	}

	logger := logger2.NewLogger(zapLogger)

	metricsRecorder := prometheus.NewMetricsRecorder(prometheus.MetricsConfig{
		Namespace: cfg.Metrics.Namespace,
		Subsystem: cfg.Metrics.Subsystem,
	})

	prometheusServer, prometheusServerCleanUp := transport.NewPrometheusMetricsServer(
		transport.PrometheusMetricsServerConfig{Addr: cfg.Metrics.Addr},
		logger,
		metricsRecorder.Registry,
	)

	redisConn, redisConnCleanUp, err := redis.NewConnection(cfg.Redis, logger)
	if err != nil {
		log.Fatal("failed connection to redis")
	}

	redisRepo := repository.NewRedisRepository(redisConn)

	hashService := service.NewHashService(redisRepo, logger)
	redirectService := service.NewRedirectService(redisRepo, logger)
	urlShortenerService := service.NewURLShortenerService(hashService, redisRepo, logger, cfg.API.BaseURL)

	createHandler := handlers.NewCreateHandler(urlShortenerService, logger, metricsRecorder)
	redirectHandler := handlers.NewRedirectHandler(redirectService, logger, metricsRecorder)

	fastHTTPHandlers := transport.NewFastHTTPHandlers(createHandler, redirectHandler)
	router := transport.NewFastHTTPRouter(fastHTTPHandlers)

	server, serverCleanUp := transport.NewFastHTTPServer(router, logger)

	interruptionChannel := make(chan os.Signal, 1)
	var g run.Group

	g.Add(func() error {
		signal.Notify(interruptionChannel, syscall.SIGINT, syscall.SIGTERM)
		osSignal := <-interruptionChannel

		return fmt.Errorf("%w: %s", ErrOsSignal, osSignal)
	}, func(error) {
		interruptionChannel <- syscall.SIGINT
	})

	g.Add(func() error {
		logger.LogInfo("fast http server",
			fmt.Sprintf("started and listening at http://localhost%s", cfg.Addr),
		)

		return server.ListenAndServe(cfg.Addr)
	}, func(err error) {
		logger.LogError("fast http server", err)
		serverCleanUp()
	})

	g.Add(func() error {
		logger.LogInfo(
			"prometheus server",
			fmt.Sprintf("started and listening for incoming requests at: "+
				"http://localhost%s/metrics", cfg.Metrics.Addr))

		return prometheusServer.ListenAndServe()
	}, func(err error) {
		logger.LogError("prometheus server", err)
		prometheusServerCleanUp()
	})

	{
		logger.LogInfo("app started")
		logger.LogError("error", g.Run())
	}

	{
		logger.LogInfo("exited")
		redisConnCleanUp()
		cleanupZapLogger()
	}

}
