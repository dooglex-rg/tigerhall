package main

import (
	"database/sql"
	"log"
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
	sighting_info(seen_time,latitude,longitude,image,tiger_id) 
		VALUES( $3, $4, $5, $6, (SELECT id FROM rows) )
	RETURNING tiger_id, id;`

	img_path, _ := save_tiger_image(c, r.Data.SightingId)
	stmt, err := DB.Prepare(sql_code)
	CheckError(err, nil)
	defer stmt.Close()
	row := stmt.QueryRow(p.Name, p.Dob, p.LastSeen, p.Latitude, p.Longitude, img_path)

	err = row.Scan(&r.Data.TigerId, &r.Data.SightingId)
	CheckError(err, sql.ErrNoRows)

	return c.JSON(r)
}

//Check if the tiger already exists in the database
func show_tigers(c *fiber.Ctx) error {
	page_no, _ := strconv.Atoi(c.Query("page", "1"))

	sql_code := `
	SELECT 
		si.tiger_id, 
		tb.name, 
		tb.dob, 
		si.seen_time, 
		si.latitude, 
		si.longitude, 
		si.image,
		count(si.id) OVER() AS full_count
	FROM sighting_info si
	
	LEFT JOIN tiger_bio tb
	ON si.tiger_id = tb.id 
	
	WHERE si.seen_time = (
		SELECT MAX(si2.seen_time)
		FROM sighting_info si2
		WHERE si2.tiger_id = si.tiger_id
	)
	ORDER BY dob DESC
	LIMIT 10
	OFFSET $1;`

	stmt, err := DB.Prepare(sql_code)
	CheckError(err, nil)
	defer stmt.Close()

	//offset logic
	rows, err := stmt.Query((page_no - 1) * 10)
	CheckError(err, sql.ErrNoRows)

	var r ResponseShowTigers
	for rows.Next() {
		var data ShowTigerModel
		err := rows.Scan(&data.TigerId, &data.Name, &data.Dob, &data.LastSeen, &data.Latitude, &data.Longitude, &data.Image, &r.Data.Count)
		r.Data.Tigers = append(r.Data.Tigers, data)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	if r.Data.Count == 0 {
		r.Status.Error = true
		r.Status.Message = "No results found!"
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
	sighting_info(seen_time,latitude,longitude,image,tiger_id) 
	VALUES( $1, $2, $3, $4, $5 )
	RETURNING id;`

	img_path, _ := save_tiger_image(c, r.Data.SightingId)
	stmt, err := DB.Prepare(sql_code)
	CheckError(err, nil)
	defer stmt.Close()
	row := stmt.QueryRow(p.LastSeen, p.Latitude, p.Longitude, img_path, p.TigerId)
	err = row.Scan(&r.Data.SightingId)

	if err != nil {
		r.Data.SightingId = 0
		r.Status.Message = "There was an error creating the record. Make sure whether the given tiger_id already exists."
		c.Status(400)
		r.Status.Error = true
		return c.JSON(r)
	}

	r.Data.TigerId = p.TigerId
	//save_tiger_image(c, r.Data.SightingId)

	return c.JSON(r)
}

//Check if the tiger already exists in the database
func show_sighting(c *fiber.Ctx) error {
	page_no, _ := strconv.Atoi(c.Query("page", "1"))
	var tiger_id TigerIdModel
	c.BodyParser(&tiger_id)
	sql_code := `
	SELECT 		
		si.seen_time, 
		si.latitude, 
		si.longitude, 
		si.image,
		count(si.id) OVER() AS full_count
	FROM sighting_info si

	WHERE si.tiger_id = $1
	LIMIT 10
	OFFSET $2;`

	stmt, err := DB.Prepare(sql_code)
	CheckError(err, nil)
	defer stmt.Close()

	//offset logic
	rows, err := stmt.Query(tiger_id.TigerId, (page_no-1)*10)
	CheckError(err, sql.ErrNoRows)

	var r ResponseShowSighting
	for rows.Next() {
		var data SightingInfo
		err := rows.Scan(&data.LastSeen, &data.Latitude, &data.Longitude, &data.Image, &r.Data.Count)
		r.Data.Sightings = append(r.Data.Sightings, data)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	if r.Data.Count == 0 {
		r.Status.Error = true
		r.Status.Message = "No results found!"
	}

	return c.JSON(r)
}
