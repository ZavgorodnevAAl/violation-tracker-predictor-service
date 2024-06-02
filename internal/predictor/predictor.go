package predictor

import (
	"service/internal/features/logging"
	onnxsession "service/internal/onnx_session"
	"sync"

	"github.com/pkg/errors"
)

type Predictor struct {
	workersNum             uint64
	videoCaptureFolderPath string
	log                    logging.Logger

	query    chan *work
	allIDMap map[uint64]*work

	rwLocker sync.RWMutex
	incID    uint64
}

func (p *Predictor) init() {
	p.query = make(chan *work, 32)
	p.allIDMap = make(map[uint64]*work)

	p.removeFilesInFolder()

	onnxSession := onnxsession.NewONNXSession(640, 640, []string{"human"}, "./data/yolov8n.onnx", p.log)

	for i := uint64(0); i < p.workersNum; i++ {
		go func() {
			for {
				select {
				case w, ok := <-p.query:
					if !ok {
						err := errors.New("unable to get work from p.query")
						w.Error(err)
					}

					w.Run(onnxSession, 0.6, 10)
					w.Wait()
				}
			}
		}()
	}

}

func NewPredictor(
	workersNum uint64,
	videoCaptureFolderPath string,
	log logging.Logger,
) *Predictor {
	out := &Predictor{
		workersNum:             workersNum,
		videoCaptureFolderPath: videoCaptureFolderPath,
		log:                    log,
	}
	out.init()

	return out
}
