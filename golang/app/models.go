package main

/*"models.go" file defines various structs used in this app*/

//A basic API Response containing default fields for error status
type ErrorStatus struct {
	//status of the error occurence in the current response
	Status struct {
		//Whether the current response processed successfully
		Error bool `json:"error"`
		//Error message incase of any error.
		Message string `json:"message"`
	} `json:"status"`
}

//Bio of the tiger
type PayloadTigerBio struct {
	//Name of the tiger
	Name string `json:"name" form:"name"`
	//Date of birth of the tiger. Must be in YYYY-MM-DD format.
	Dob string `json:"birthday" form:"birthday"`
}

type SightingInfo struct {
	//Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.
	LastSeen string `json:"last_seen" form:"last_seen"`
	//Last seen Latitude point
	Latitude float64 `json:"latitude" form:"latitude"`
	//Last seen Longitude point
	Longitude float64 `json:"longitude" form:"longitude"`
}

//Parsing Request payload to add a new tiger
type PayloadAddNewTiger struct {
	PayloadTigerBio
	SightingInfo
}

type PayloadAddSighting struct {
	SightingInfo
	TigerIdModel
}

//id of the the Tiger
type TigerIdModel struct {
	TigerId int64 `json:"tiger_id" form:"tiger_id"`
}

//Outgoing Response format for adding a new tiger
type ResponseTiger struct {
	ErrorStatus
	//Data field
	Data struct {
		TigerIdModel
		SightingId int64 `json:"sighting_id" form:"sighting_id"`
	} `json:"data"`
}
