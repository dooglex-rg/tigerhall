package main

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//A sample JSON response for index page
func index_page(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Ping": "Pong",
	})
}

//Create a new tiger along with the last seen info
func create_tiger(c *fiber.Ctx) error {
	//payload variable
	var p PayloadAddNewTiger

	//response variable
	var r ResponseNewTiger

	//parsing payload JSON to struct
	c.BodyParser(&p)
	//Validating the payload fields
	switch {
	case p.Name == "" || p.Dob == "" || p.LastSeen == "":
		r.Status.Message = "name/birthday/last_seen fields should not be blank/zero/nil value"
	case len(p.GeoLocation) != 2:
		r.Status.Message = "cordinates field should contain exactly 2 values in the array. ie., [lat,lon]"
	default:
		r.Status.Success = true
	}

	if !r.Status.Success {
		return c.JSON(r)
	}

	sql_code := `WITH rows AS ( 
		INSERT INTO tiger_bio(name,dob) 
		VALUES( $1, $2) 
		RETURNING id
	)

	INSERT INTO 
		last_seen(seen_time,latitude,longitude,tiger_id) 
		VALUES( $3, $4, $5, (SELECT id FROM rows) )
	RETURNING tiger_id;`

	row := DB.QueryRow(sql_code, p.Name, p.Dob, p.LastSeen, p.GeoLocation[0], p.GeoLocation[1])

	err := row.Scan(&r.Data.Id)
	CheckError(err)

	save_tiger_image(c, r.Data.Id)

	return c.JSON(r)
}

func save_tiger_image(c *fiber.Ctx, id int64) error {
	file, _ := c.FormFile("image")
	filename := os.Getenv("IMAGE_FOLDER") + strconv.FormatInt(id, 10) + file.Filename
	return c.SaveFile(file, filename)
}

//Check if the tiger already exists in the database
func check_tiger(c *fiber.Ctx) error {
	//payload variable
	var p PayloadTigerBio

	//response variable
	var r ResponseNewTiger

	//parsing payload JSON to struct
	c.BodyParser(&p)
	//Validating the payload fields
	switch {
	case p.Name == "" || p.Dob == "":
		r.Status.Message = "name/birthday field should not be blank/zero/nil value"
		return c.JSON(r)
	default:
		r.Status.Success = true
	}

	sql_code := `SELECT id FROM tiger_bio WHERE name=$1 AND dob=$2 LIMIT 1;`

	row := DB.QueryRow(sql_code, p.Name, p.Dob)

	err := row.Scan(&r.Data.Id)

	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}

	return c.JSON(r)
}

/*
//Create a new sighting of a tiger
func create_sighting(c *fiber.Ctx) error {

	return c.JSON("response")
}
*/
