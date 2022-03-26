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
	var r ResponseTiger

	//parsing payload JSON to struct
	c.BodyParser(&p)
	//Validating the payload fields
	if p.Name == "" || p.Dob == "" || p.LastSeen == "" || p.Latitude == 0 || p.Longitude == 0 {
		r.Status.Message = "name/birthday/last_seen/geo fields should not be blank/zero/nil value"
		c.Status(400)
		r.Status.Error = true
		return c.JSON(r)
	}

	sql_code := `WITH rows AS ( 
		INSERT INTO tiger_bio(name,dob) 
		VALUES( $1, $2) 
		RETURNING id
	)

	INSERT INTO 
	sighting_info(seen_time,latitude,longitude,tiger_id) 
		VALUES( $3, $4, $5, (SELECT id FROM rows) )
	RETURNING tiger_id, id;`

	row := DB.QueryRow(sql_code, p.Name, p.Dob, p.LastSeen, p.Latitude, p.Longitude)

	err := row.Scan(&r.Data.TigerId, &r.Data.SightingId)
	CheckError(err)

	save_tiger_image(c, r.Data.SightingId)

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
	var r ResponseTiger

	//parsing payload JSON to struct
	c.BodyParser(&p)
	//Validating the payload fields
	if p.Name == "" || p.Dob == "" {
		r.Status.Message = "name/birthday field should not be blank/zero/nil value"
		c.Status(400)
		r.Status.Error = true
		return c.JSON(r)
	}

	sql_code := `SELECT id FROM tiger_bio WHERE name=$1 AND dob=$2 LIMIT 1;`

	stmt, err := DB.Prepare(sql_code)
	CheckError(err)
	defer stmt.Close()

	row := stmt.QueryRow(p.Name, p.Dob)

	err = row.Scan(&r.Data.TigerId)

	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}

	return c.JSON(r)
}

//Create a new sighting of existing tiger
func create_sighting(c *fiber.Ctx) error {
	//payload variable
	var p PayloadAddSighting
	var r ResponseTiger

	//parsing payload JSON to struct
	c.BodyParser(&p)

	switch {
	case p.TigerId == 0:
		r.Status.Message = "tiger_id is required to add new sighting of existing tiger. To add a record for a new tiger & its sighting, use '/tiger/add' endpoint."
		c.Status(400)
		r.Status.Error = true
		return c.JSON(r)
	case p.LastSeen == "" || p.Latitude == 0 || p.Longitude == 0:
		r.Status.Message = "last_seen/geo fields should not be blank/zero/nil value"
		c.Status(400)
		r.Status.Error = true
		return c.JSON(r)
	}

	sql_code := `INSERT INTO 
	sighting_info(seen_time,latitude,longitude,tiger_id) 
	VALUES( $1, $2, $3, $4 )
	RETURNING id;`

	row := DB.QueryRow(sql_code, p.LastSeen, p.Latitude, p.Longitude, p.TigerId)
	err := row.Scan(&r.Data.SightingId)

	if err != nil {
		r.Data.SightingId = 0
		r.Status.Message = "There was an error creating the record. Make sure whether the given tiger_id already exists."
		c.Status(400)
		r.Status.Error = true
		return c.JSON(r)
	}

	r.Data.TigerId = p.TigerId
	save_tiger_image(c, r.Data.SightingId)

	return c.JSON(r)
}
