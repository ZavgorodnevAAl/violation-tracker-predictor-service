package onnxsession

import (
	"image"
	"math"
	"sort"

	"gocv.io/x/gocv"
)

func (s *OnnxSession) PrepareInput(mat *gocv.Mat) []float32 {
	resizedImg := gocv.NewMat()
	gocv.Resize(*mat, &resizedImg, image.Pt(s.width, s.height), 0, 0, gocv.InterpolationCubic)

	red := []float32{}
	green := []float32{}
	blue := []float32{}
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			rgb := resizedImg.GetVecbAt(y, x)
			red = append(red, float32(rgb[2])/255.0)
			green = append(green, float32(rgb[1])/255.0)
			blue = append(blue, float32(rgb[0])/255.0)
		}
	}

	input := append(red, green...)
	input = append(input, blue...)
	return input
}

func union(box1, box2 []interface{}) float64 {
	box1_x1, box1_y1, box1_x2, box1_y2 := box1[0].(float64), box1[1].(float64), box1[2].(float64), box1[3].(float64)
	box2_x1, box2_y1, box2_x2, box2_y2 := box2[0].(float64), box2[1].(float64), box2[2].(float64), box2[3].(float64)
	box1_area := (box1_x2 - box1_x1) * (box1_y2 - box1_y1)
	box2_area := (box2_x2 - box2_x1) * (box2_y2 - box2_y1)
	return box1_area + box2_area - intersection(box1, box2)
}

func intersection(box1, box2 []interface{}) float64 {
	box1_x1, box1_y1, box1_x2, box1_y2 := box1[0].(float64), box1[1].(float64), box1[2].(float64), box1[3].(float64)
	box2_x1, box2_y1, box2_x2, box2_y2 := box2[0].(float64), box2[1].(float64), box2[2].(float64), box2[3].(float64)
	x1 := math.Max(box1_x1, box2_x1)
	y1 := math.Max(box1_y1, box2_y1)
	x2 := math.Min(box1_x2, box2_x2)
	y2 := math.Min(box1_y2, box2_y2)
	return (x2 - x1) * (y2 - y1)
}

func iou(box1, box2 []interface{}) float64 {
	return intersection(box1, box2) / union(box1, box2)
}

type predictOutput struct {
	PointMin image.Point
	PointMax image.Point
	Label    string
	Percent  float32
}

func (s *OnnxSession) PrepareOutput(output []float32, persent float64) []predictOutput {
	boxes := [][]interface{}{}
	for index := 0; index < 8400; index++ {
		class_id, prob := 0, float32(0.0)
		for col := 0; col < len(s.yoloClasses); col++ {
			if output[8400*(col+4)+index] > prob {
				prob = output[8400*(col+4)+index]
				class_id = col
			}
		}
		if prob < 0.5 {
			continue
		}
		label := s.yoloClasses[class_id]
		xc := output[index]
		yc := output[8400+index]
		w := output[2*8400+index]
		h := output[3*8400+index]
		x1 := (xc - w/2) / 640 * float32(s.width)
		y1 := (yc - h/2) / 640 * float32(s.height)
		x2 := (xc + w/2) / 640 * float32(s.width)
		y2 := (yc + h/2) / 640 * float32(s.height)
		boxes = append(boxes, []interface{}{float64(x1), float64(y1), float64(x2), float64(y2), label, prob})
	}

	sort.Slice(boxes, func(i, j int) bool {
		return boxes[i][5].(float32) < boxes[j][5].(float32)
	})
	result := [][]interface{}{}
	for len(boxes) > 0 {
		result = append(result, boxes[0])
		tmp := [][]interface{}{}
		for _, box := range boxes {
			if iou(boxes[0], box) < persent {
				tmp = append(tmp, box)
			}
		}
		boxes = tmp
	}

	res := []predictOutput{}

	for _, v := range result {
		res = append(res,
			predictOutput{
				PointMin: image.Point{
					X: int(v[0].(float64)),
					Y: int(v[1].(float64)),
				},
				PointMax: image.Point{
					X: int(v[2].(float64)),
					Y: int(v[3].(float64)),
				},
				Label:   v[4].(string),
				Percent: v[5].(float32),
			},
		)
	}

	return res
}
