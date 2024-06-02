package main

import (
	"os"
	"os/signal"
	"service/internal/config"
	"service/internal/cors"
	"service/internal/features/logging"
	"service/internal/http_server"
	"service/internal/predictor"
	"syscall"

	"github.com/pkg/errors"
	"github.com/yalue/onnxruntime_go"
)

func main() {
	// Заводим канал для ожидания сигнала со стороны os
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Подготавливаем конфиг
	conf, err := config.New()
	if err != nil {
		err = errors.Wrap(err, "unable to create config: [config.New]")
		panic(err)
	}

	// Подготавливаем логгирование
	logConfig := logging.NewConfig()
	logConfig.SetInFile(conf.Logger.LogInFile)
	logConfig.SetOutputDir(conf.Logger.OutputDir)
	logConfig.SetLevel(logging.InfoLevel)
	log, err := logging.New(logConfig)
	if err != nil {
		err = errors.Wrap(err, "unable to create logging: [logging.New()]")
		panic(err)
	}

	onnxruntime_go.SetSharedLibraryPath("./data/onnx/libonnxruntime.so.1.17.1")

	predictorN := predictor.NewPredictor(3, "./video_capture", log)

	serv := http_server.NewHttpServer(
		conf.Service.Host,
		conf.Service.Port,
		predictorN,
		log,
	)

	// Запускаем прокси с использованием cors
	log.Info("start cors")
	cors := cors.NewCors(conf.Cors.Target, conf.Cors.Path)
	cors.Run()

	serv.Start()

	<-quit

	serv.Stop()
}
