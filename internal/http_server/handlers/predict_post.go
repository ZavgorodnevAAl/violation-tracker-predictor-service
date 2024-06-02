package handlers

import (
	"service/internal/gen/models"
	"service/internal/gen/restapi/operations/predictor_tag"

	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
)

func (h *Handler) uploadPost500(err error) middleware.Responder {
	err = errors.Wrapf(err, "handler error: [predictPost]")
	h.log.Error(err.Error())
	return predictor_tag.NewUploadPostInternalServerError().WithPayload(
		&models.Error500{
			Error: err.Error(),
		},
	)
}

func (h *Handler) uploadPost(params predictor_tag.UploadPostParams) middleware.Responder {
	token, err := h.predictor.Upload(params.Body.VideoBase64)
	if err != nil {
		err = errors.Wrap(err, "[h.predictor.Upload(params.Body.VideoBase64)]")
		return h.uploadPost500(err)
	}

	out := &predictor_tag.UploadPostOKBody{}
	out.Token = token

	return predictor_tag.NewUploadPostOK().WithPayload(out)
}
