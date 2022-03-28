package main

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

func main() {
	godotenv.Load(".env")
	redis_url := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	worker := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redis_url},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 8,
				"default":  1,
				"low":      1,
			},
		},
	)

	//Queue server listening for new tasks
	server := asynq.NewServeMux()
	server.HandleFunc(os.Getenv("TASK_NAME"), process_image_resize)
	err := worker.Run(server)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func process_image_resize(c context.Context, t *asynq.Task) error {
	img_name := string(t.Payload())
	img_path := os.Getenv("IMAGE_FOLDER") + img_name
	file, err := os.Open(img_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var img image.Image

	switch filepath.Ext(img_path) {
	case ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	}
	if err != nil {
		log.Panic(err)
	}

	// resize to using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(250, 200, img, resize.Lanczos3)

	output_path := os.Getenv("IMAGE_FOLDER_RESIZED") + img_name
	out, err := os.Create(output_path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file

	jpeg.Encode(out, m, nil)
	return nil
}
