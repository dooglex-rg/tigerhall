package main

import (
	"github.com/gofiber/fiber/v2"
)

//A sample JSON response for index page
func index_page(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Ping": "Pong",
	})
}

//Create a new tiger along with last seen info
func create_tiger(c *fiber.Ctx) error {
	//payload variable
	var p PayloadNewTiger
	var response ResponseNewTiger

	//Validating the payload fields
	switch {
	case p.Name == "":
		response.Status.Message = "name field should not be blank"
	case p.Dob.IsZero() || p.LastSeen.IsZero():
		response.Status.Message = "birthday & last_seen fields should not be blank/zero/nil value"
	case len(p.GeoLocation) == 2:
		response.Status.Message = "cordinates field should contain exactly 2 values in the array. ie., [lat,lon]"
	default:
		response.Status.Success = true
	}

	if !response.Status.Success {
		return c.JSON(response)
	}

	sql := `WITH rows AS ( 
		INSERT INTO tiger_bio(name,dob) 
		VALUES( $1, $2) 
		RETURNING id
	)

	INSERT INTO 
		last_seen(seen_time,latitude,longitude,tiger_id) 
		VALUES( $3, $4, $5, (SELECT id FROM rows) )
	RETURNING id;`

	row := DB.QueryRow(sql, p.Name, p.Dob, p.LastSeen, p.GeoLocation[0], p.GeoLocation[1])

	err := row.Scan(&response.Data.Id)
	CheckError(err)

	return c.JSON(response)
}
