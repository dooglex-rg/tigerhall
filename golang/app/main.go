package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	err := godotenv.Load(".env")
	CheckError(err)

	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	//Creating DB connection
	DB, err = sql.Open("postgres", DSN)
	CheckError(err)

	defer DB.Close()

	//Ping test to db
	CheckError(DB.Ping())

	//Based on free tier limit
	//https://www.elephantsql.com/plans.html
	DB.SetConnMaxLifetime(time.Minute * 5)
	DB.SetMaxOpenConns(2)
	DB.SetMaxIdleConns(2)

	//Fiber app instance
	app := fiber.New()
	middleware_config(app)
	url_router(app)

	log.Fatal(app.Listen(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")))
}

//Routing of incoming URLs with its handlers
func url_router(app *fiber.App) {
	app.Get("/", index_page)
	app.Post("/tiger/new", create_tiger)
	app.Post("/tiger/check", check_tiger)
}

//Middlewares configuration
func middleware_config(app *fiber.App) {

	//CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:40619",
	}))
}

//General error checking function
func CheckError(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}
