package onnxsession

import (
	"image"
	"service/internal/features/logging"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/yalue/onnxruntime_go"
	"gocv.io/x/gocv"
)

type OnnxSession struct {
	width       int
	height      int
	yoloClasses []string
	model       string

	inputTensor  *onnxruntime_go.Tensor[float32]
	outputTensor *onnxruntime_go.Tensor[float32]

	session *onnxruntime_go.AdvancedSession
	log     logging.Logger
}

func (ns *OnnxSession) init() {
	inputShape := []int64{1, 3, int64(ns.width), int64(ns.height)}
	outputShape := []int64{1, 84, 8400}

	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		err = errors.Wrap(err, "[onnxruntime_go.NewEmptyTensor[float32](outputShape)]")
		panic(err)
	}

	NewShapeInput := onnxruntime_go.NewShape(inputShape...)

	var n int64
	n = 1
	for _, v := range inputShape {
		n *= v
	}

	ns.inputTensor, err = onnxruntime_go.NewTensor(NewShapeInput, make([]float32, n))
	if err != nil {
		err = errors.Wrapf(err, "[onnxruntime_go.NewTensor(inputShape, []float32{})] - width: %d, height: %d", ns.width, ns.height)
		panic(err)
	}

	newShapeOutput := onnxruntime_go.NewShape(outputShape...)
	ns.outputTensor, err = onnxruntime_go.NewEmptyTensor[float32](newShapeOutput)
	if err != nil {
		err = errors.Wrap(err, "[onnxruntime_go.NewEmptyTensor[float32](outputShape)]")
		panic(err)
	}

	options, err := onnxruntime_go.NewSessionOptions()
	if err != nil {
		err = errors.Wrap(err, "[onnxruntime_go.NewSessionOptions()]")
		panic(err)
	}

	ns.session, err = onnxruntime_go.NewAdvancedSession(ns.model,
		[]string{"images"},
		[]string{"output0"},
		[]onnxruntime_go.ArbitraryTensor{ns.inputTensor},
		[]onnxruntime_go.ArbitraryTensor{ns.outputTensor},
		options,
	)
	if err != nil {
		err = errors.Wrap(err, "[ort.NewAdvancedSession()]")
		panic(err)
	}
}

func (onxs *OnnxSession) Predict(mat gocv.Mat, persent float64) (out []predictOutput, err error) {
	t := time.Now()

	// Изменяем размер изображения
	newMap := gocv.NewMat()
	gocv.Resize(mat, &newMap, image.Point{onxs.height, onxs.width}, 0, 0, gocv.InterpolationLinear)

	input := onxs.PrepareInput(&newMap)

	inTensor := onxs.inputTensor.GetData()
	copy(inTensor, input)
	err = onxs.session.Run()
	if err != nil {
		err = errors.Wrap(err, "[session.Run()] - unable to calc worker data")
		onxs.log.Error(err.Error())
		return nil, err
	}
	outTens := onxs.outputTensor.GetData()

	out = onxs.PrepareOutput(outTens, persent)
	onxs.log.Info(strconv.FormatInt(time.Since(t).Milliseconds(), 10))

	return out, nil
}

func NewONNXSession(width int, height int, yoloClasses []string, model string, log logging.Logger) *OnnxSession {
	out := &OnnxSession{
		width:       width,
		height:      height,
		yoloClasses: yoloClasses,
		model:       model,
		log:         log,
	}
	out.init()

	return out
}
