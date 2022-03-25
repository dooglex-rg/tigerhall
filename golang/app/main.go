package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func main() {
	const (
		host     = "kandula.db.elephantsql.com"
		port     = 5432
		user     = "nrkisatu"
		password = "rzrCj-txNNA1zJhC8xdGy6hmXvr7vlkG"
		dbname   = "nrkisatu"
	)
	DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error

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

	log.Fatal(app.Listen("localhost:40619"))
}

//Routing of incoming URLs with its handlers
func url_router(app *fiber.App) {
	app.Get("/", index_page)
	app.Post("/tiger/new", create_tiger)
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
		panic(err)
	}
}
