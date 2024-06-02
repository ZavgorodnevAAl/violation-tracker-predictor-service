package predictor

import (
	"fmt"
	"service/internal/gen/models"
	onnxsession "service/internal/onnx_session"
	"sync"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

type work struct {
	id    uint64
	video *gocv.VideoCapture

	inQuery bool
	persent float64
	err     error

	out []int64

	rwinfo sync.RWMutex

	mut  sync.Mutex
	done chan struct{}
}

func (w *work) Info() (*models.Info200, error) {
	w.rwinfo.RLock()
	defer w.rwinfo.RUnlock()

	out := &models.Info200{
		ID:      w.id,
		Persent: w.persent,
		Out:     w.out,
	}

	return out, w.err
}

func (w *work) getID() uint64 {
	w.rwinfo.RLock()
	defer w.rwinfo.RUnlock()

	return w.id
}

func (w *work) setPersent(v float64) {
	w.rwinfo.Lock()
	defer w.rwinfo.Unlock()

	w.persent = v
}

func (w *work) addOut(v int64) {
	w.rwinfo.Lock()
	defer w.rwinfo.Unlock()

	w.out = append(w.out, v)
}

func (w *work) Run(onnxSession *onnxsession.OnnxSession, persent float64, framesPerPeriod int) {
	w.mut.Lock()
	defer w.mut.Unlock()

	defer close(w.done)

	defer w.video.Close()

	w.inQuery = true

	frameRate := w.video.Get(gocv.VideoCaptureFPS)
	frameDuration := int(frameRate) * framesPerPeriod

	// Получаем общее количество кадров в видеопотоке
	frameCount := w.video.Get(gocv.VideoCaptureFrameCount)

	c := 0

	for w.video.IsOpened() {
		frame := gocv.NewMat()
		w.setPersent(frameCount / float64(c))

		if ok := w.video.Read(&frame); !ok {
			break
		}
		if c%frameDuration != 0 {
			c++
			continue
		}

		out, err := onnxSession.Predict(frame, persent)
		if err != nil {
			err = errors.Wrap(err, "[h.onnxSession.Predict(frame, 0.6)")
			w.Error(err)
			return
		}

		if len(out) > 0 {
			w.addOut(int64(c))
		}

		frame.Close()
		c++

		fmt.Println(out)
	}

}

func (w *work) Error(err error) {
	w.mut.Lock()
	defer w.mut.Unlock()
	close(w.done)
	w.err = err
}

func (w *work) Wait() error {
	<-w.done
	w.mut.Lock()
	defer w.mut.Unlock()

	if w.err != nil {
		return w.err
	}
	return nil
}

func (w *work) Close() {

}

func (w *work) init() {
	w.inQuery = true
	w.done = make(chan struct{})
}

func newWork(id uint64, video *gocv.VideoCapture) *work {
	out := &work{
		id:    id,
		video: video,
	}
	out.init()
	return out
}
