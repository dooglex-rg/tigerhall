package main

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

func main() {
	//load enviroinment variable
	godotenv.Load(".env")

	redis_url := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	//worker config for listening to new tasks
	worker := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redis_url,
			Username: os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASS"),
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 8,
				"default":  1,
				"low":      1,
			},
		},
	)

	//Queue server created
	server := asynq.NewServeMux()
	//task assignment
	server.HandleFunc(os.Getenv("TASK_NAME"), process_image_resize)

	//server stated listening
	err := worker.Run(server)
	CheckError(err, nil)
}

//image resize function
func process_image_resize(c context.Context, t *asynq.Task) error {
	var i map[string]string
	//fetching task related data
	json.Unmarshal(t.Payload(), &i)
	img_name := i["filename"]
	img_path := os.Getenv("IMAGE_FOLDER") + img_name

	file, err := os.Open(img_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var img image.Image

	switch i["ext"] {
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return errors.New("image format not supported")
	}
	CheckError(err, nil)

	// resize to using Lanczos resampling
	// and preserves aspect ratio
	m := resize.Resize(250, 200, img, resize.Lanczos3)

	out, err := os.Create(img_path)
	CheckError(err, nil)
	defer out.Close()

	// write resized image to file
	jpeg.Encode(out, m, nil)

	return nil
}

//error checking with provision for exemption
func CheckError(err, exempt error) {
	switch err {
	case nil, exempt:
		return
	default:
		log.Panic(err)
	}
}
