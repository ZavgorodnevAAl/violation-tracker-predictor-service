package handlers

import (
	"service/internal/features/logging"
	"service/internal/gen/restapi/operations"
	"service/internal/gen/restapi/operations/predictor_tag"
	"service/internal/predictor"
)

type Handler struct {
	predictor *predictor.Predictor
	log       logging.Logger
}

func NewHandler(predictor *predictor.Predictor, log logging.Logger) *Handler {
	return &Handler{
		predictor: predictor,
		log:       log,
	}
}

func (h *Handler) Register(api *operations.PredictorServiceAPI) {
	api.PredictorTagUploadPostHandler = predictor_tag.UploadPostHandlerFunc(h.uploadPost)
	api.PredictorTagInfoGetHandler = predictor_tag.InfoGetHandlerFunc(h.infoGet)
}
