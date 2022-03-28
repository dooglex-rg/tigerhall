package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var client *asynq.Client

func main() {
	err := godotenv.Load(".env")
	CheckError(err, nil)

	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	//Creating DB connection
	DB, err = sql.Open("postgres", DSN)
	CheckError(err, nil)

	defer DB.Close()

	//Ping test to db
	CheckError(DB.Ping(), nil)

	//Based on free tier limit
	//https://www.elephantsql.com/plans.html
	DB.SetConnMaxLifetime(time.Minute * 1)
	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	//Fiber app instance
	app := fiber.New()
	middleware_config(app)
	url_router(app)

	redis_url := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redis_url,
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	})
	defer client.Close()

	server_url := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	log.Fatal(app.Listen(server_url))
}

//Routing of incoming URLs with its handlers
func url_router(app *fiber.App) {
	app.Get("/", index_page)

	app.Post("/tiger/add", create_tiger)
	app.Post("/tiger/show", show_tigers)

	app.Post("/sighting/add", create_sighting)
	app.Post("/sighting/show", show_sighting)
}

//Middlewares configuration
func middleware_config(app *fiber.App) {

	//CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:40619",
	}))
}

//General error checking function with provision for given exemption
func CheckError(err, exempt error) {
	switch err {
	case nil, exempt:
		return
	default:
		log.Panic(err)
	}
}

func schedule_image_resize(file_name, id, file_extention string) {
	i := map[string]string{
		"filename": file_name,
		"id":       id,
		"ext":      file_extention,
	}
	json_byte, _ := json.Marshal(i)
	task_item := asynq.NewTask(os.Getenv("TASK_NAME"), json_byte, asynq.MaxRetry(1))
	client.Enqueue(task_item, asynq.Queue("critical"))
}

func save_tiger_image(c *fiber.Ctx, id int64) (string, error) {
	//file content
	file_stream, _ := c.FormFile("image")
	//name of the file
	file_extention := filepath.Ext(file_stream.Filename)
	file_name := uuid.New().String() + file_extention
	//file destination path
	file_path := os.Getenv("IMAGE_FOLDER") + file_name
	//scheduling for image resizing queue
	schedule_image_resize(file_name, strconv.FormatInt(id, 10), file_extention)
	return file_path, c.SaveFile(file_stream, file_path)
}
