package handlers

import (
	"service/internal/gen/models"
	"service/internal/gen/restapi/operations/predictor_tag"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
)

func (h *Handler) infoGet500(err error) middleware.Responder {
	err = errors.Wrapf(err, "handler error: [predictPost]")
	h.log.Error(err.Error())
	return predictor_tag.NewInfoGetInternalServerError().WithPayload(
		&models.Error500{
			Error: err.Error(),
		},
	)
}

func (h *Handler) infoGet(params predictor_tag.InfoGetParams) middleware.Responder {
	out, err := h.predictor.GetInfo(params.Token)
	if err != nil {
		err = errors.Wrap(err, "[h.predictor.GetInfo(params.Token)]")
		return h.infoGet500(err)
	}

	return predictor_tag.NewInfoGetOK().WithPayload(out)
}
