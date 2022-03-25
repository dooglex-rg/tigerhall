package main

import "time"

/*"models.go" file defines various structs used in this app*/

//A basic API Response containing default fields for error status
type ErrorStatus struct {
	//status of the error occurence in the current response
	Status struct {
		//Whether the current response processed successfully
		Success bool `json:"success"`
		//Error message incase of any error.
		Message string `json:"message"`
	} `json:"status"`
}

//Parsing Request payload to add a new tiger
type PayloadNewTiger struct {
	//Name of the tiger
	Name string `json:"name"`
	//Date of birth of the tiger
	Dob time.Time `json:"birthday"`
	//Timestamp when the tiger was last seen
	LastSeen time.Time `json:"last_seen"`
	//Last seen cordinate points [lat,lon]
	GeoLocation [2]float64 `json:"cordinates"`
}

//Outgoing Response format for adding a new tiger
type ResponseNewTiger struct {
	ErrorStatus
	//Data field
	Data struct {
		//id of the newly created Tiger
		Id int64 `json:"tiger_id"`
	} `json:"data"`
}
