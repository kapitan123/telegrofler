package faceFinder

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"
	//visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

//go:embed test.jpg
var test []byte

type FaceFinder struct {
}

func New() *FaceFinder {
	return &FaceFinder{}
}

func (ff *FaceFinder) Find(img []byte) ([]byte, error) {
	return nil, nil
}

func DetectLandmarks() error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	r := bytes.NewReader(test)
	image, err := vision.NewImageFromReader(r)
	if err != nil {
		return err
	}
	//annotations, err := client.DetectLandmarks(ctx, image, nil, 10)
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Println("No landmarks found.")
	} else {
		fmt.Println("Landmarks:")
		for _, annotation := range annotations {
			fmt.Println(annotation.Landmarks)
		}
	}

	return nil
}

func FindMouth() (float32, float32, float32, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	defer client.Close()

	if err != nil {
		return 0, 0, 0, err
	}

	r := bytes.NewReader(test)
	image, err := vision.NewImageFromReader(r)
	if err != nil {
		return 0, 0, 0, err
	}

	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(annotations) == 0 {
		return 0, 0, 0, errors.New("no landmarks found")
	} else {
		for _, annotation := range annotations {
			face := vision.FaceFromLandmarks(annotation.Landmarks)

			width := face.Mouth.Right.X - face.Mouth.Left.Y
			return face.Mouth.Left.X, face.Mouth.Left.Y, width, nil
			// for _, landmark := range annotation.Landmarks {

			// 	if landmark.Type == vision.FaceFromLandmarks() {
			// 		return landmark.Position.X, landmark.Position.Y, nil
			// 	}
			// }
		}
	}

	return 0, 0, 0, nil
}
