package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	_ "github.com/dooglex-rg/tigerhall/app/docs"
	_ "github.com/lib/pq"
)

//database instance
var DB *sql.DB

//asynq task scheduler client
var client *asynq.Client

// @title Tigerhall test API
// @version 1.0
// @description This is an swagger documentation of simple test API task given by tigerhall
// @termsOfService https://www.example.com/terms
// @contact.name Tech Support
// @contact.email rg@dooglex.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host tigerhall.dooglex.com
// @BasePath /
func main() {
	//load env variables from .env file
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

	DB.SetConnMaxLifetime(time.Minute * 1)
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)

	//Fiber app instance
	app := fiber.New()

	//middlewares setup
	middleware_config(app)

	//url routing
	url_router(app)

	redis_url := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redis_url,
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	})
	defer client.Close()

	server_url := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	//server listening
	log.Fatal(app.Listen(server_url))
}

//Routing of incoming URLs with its handlers
func url_router(app *fiber.App) {
	app.Get("/", index_page)

	app.Post("/tiger/add", create_tiger)
	app.Post("/tiger/show", show_tigers)

	app.Post("/sighting/add", create_sighting)
	app.Post("/sighting/show", show_sighting)

	app.Get("/api-docs/*", swagger.New(swagger.Config{
		Title:                    "Swagger API docs",
		URL:                      "/api-docs/doc.json",
		DefaultModelsExpandDepth: -1,
		Layout:                   "StandaloneLayout",
		Plugins: []template.JS{
			template.JS("SwaggerUIBundle.plugins.DownloadUrl"),
		},
		Presets: []template.JS{
			template.JS("SwaggerUIBundle.presets.apis"),
			template.JS("SwaggerUIStandalonePreset"),
		},
		DeepLinking:             true,
		DefaultModelExpandDepth: 1,
		DefaultModelRendering:   "example",
		DocExpansion:            "list",
		SyntaxHighlight: swagger.SyntaxHighlightConfig{
			Activate: true,
			Theme:    "agate",
		},
		ShowMutatedRequest:   true,
		ShowExtensions:       true,
		ShowCommonExtensions: true,
	}))

}

//Middlewares configuration
func middleware_config(app *fiber.App) {

	//CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: fmt.Sprintf("http://%s:%s, https://%s",
			os.Getenv("APP_HOST"), os.Getenv("APP_PORT"), os.Getenv("PUBLIC_HOST")),
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

//Schedule image resizing function for later queue to
//make responses quicker & to avoid memory crash
func schedule_image_resize(file_name, id, file_extention string) {
	i := map[string]string{
		"filename": file_name,
		"id":       id,
		"ext":      file_extention,
	}
	json_byte, _ := json.Marshal(i)
	//create task
	task_item := asynq.NewTask(os.Getenv("TASK_NAME"), json_byte, asynq.MaxRetry(1))
	//send to queue
	client.Enqueue(task_item, asynq.Queue("critical"))
}

//saves uploaded image to storage
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
