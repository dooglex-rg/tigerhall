package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {

	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	//Creating DB connection
	var err error
	DB, err = sql.Open("postgres", DSN)
	CheckError(err, nil)

	defer DB.Close()

	DB.SetConnMaxLifetime(time.Minute * 5)
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)

	//Ping test to db
	CheckError(DB.Ping(), nil)

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
	err = worker.Run(server)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func process_image_resize(c context.Context, t *asynq.Task) error {
	var i map[string]string
	json.Unmarshal(t.Payload(), &i)
	img_name := i["filename"]
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
	CheckError(err, nil)
	defer out.Close()

	// write new image to file

	jpeg.Encode(out, m, nil)

	sql_code := `
	UPDATE sighting_info 
	SET image = $1
	WHERE id = $2;`
	DB.Exec(sql_code, output_path, i["id"])
	return nil
}

func CheckError(err, exempt error) {
	switch err {
	case nil, exempt:
		return
	default:
		log.Panic(err)
	}
}
