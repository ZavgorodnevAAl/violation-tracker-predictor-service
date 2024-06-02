package http_server

import (
	"service/internal/features/logging"
	"service/internal/gen/restapi"
	"service/internal/gen/restapi/operations"
	"service/internal/http_server/handlers"
	"service/internal/predictor"

	"github.com/go-openapi/loads"
)

type HttpServer struct {
	host string
	port int
	log  logging.Logger

	predictor *predictor.Predictor

	api    *operations.PredictorServiceAPI
	server *restapi.Server
}

func (h *HttpServer) init() {
	spec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		panic(err)
	}

	api := operations.NewPredictorServiceAPI(spec)
	h.api = api

	// Подготавливаем пакет handlers для обработки запросов со стороны клиента
	handlers := handlers.NewHandler(
		h.predictor,
		h.log,
	)
	handlers.Register(h.api)

	server := restapi.NewServer(h.api)
	server.Port = h.port
	server.Host = h.host

	h.server = server
}

func (h *HttpServer) Start() {
	wait := make(chan struct{})
	go func() {
		close(wait)
		h.log.Info("start rest service")
		if err := h.server.Serve(); err != nil {
			panic(err)
		}
	}()

	<-wait
}

func (h *HttpServer) Stop() {
	h.log.Info("shutdown rest service")
	err := h.server.Shutdown()
	if err != nil {
		panic(err)
	}

}

func NewHttpServer(
	host string,
	port int,
	predictor *predictor.Predictor,
	log logging.Logger,
) *HttpServer {
	out := &HttpServer{
		host:      host,
		port:      port,
		predictor: predictor,
		log:       log,
	}
	out.init()

	return out
}
