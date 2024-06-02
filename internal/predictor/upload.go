package predictor

import (
	"encoding/base64"
	"service/internal/gen/models"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

func (p *Predictor) Upload(videoBase64 string) (token uint64, err error) {
	p.rwLocker.Lock()
	defer p.rwLocker.Unlock()

	p.incID++
	token = p.incID
	// Декодируем строку base64 в байтовый массив
	videoBytes, err := base64.StdEncoding.DecodeString(videoBase64)
	if err != nil {
		err = errors.Wrap(err, "[base64.StdEncoding.DecodeString(params.Body.Base64Image)]")
		return 0, err
	}

	err = p.uploadFile(p.incID, videoBytes)
	if err != nil {
		err = errors.Wrap(err, "[p.uploadFile(p.c, videoBytes)")
		return 0, err
	}

	fp := p.makeFullFilePath(p.incID)

	video, err := gocv.VideoCaptureFile(fp)
	if err != nil {
		err = errors.Wrap(err, "[gocv.VideoCaptureFile(fp)")
		return 0, err
	}

	worker := newWork(p.incID, video)

	select {
	case p.query <- worker:
	default:
		return 0, errors.New("server overload, channel is full")
	}

	p.allIDMap[p.incID] = worker

	return
}

func (p *Predictor) GetInfo(token uint64) (model *models.Info200, err error) {
	p.rwLocker.RLock()
	defer p.rwLocker.RUnlock()

	info, ok := p.allIDMap[token]
	if !ok {
		return nil, errors.New("unable to find")
	}

	model, err = info.Info()
	if err != nil {
		err = errors.Wrap(err, "[info.Info()]")
		return nil, err
	}

	return
}
