package main

/*"models.go" file defines various structs used in this app*/

//A basic API Response containing default fields for error status
type ErrorStatus struct {
	//status of the error occurence in the current response
	Status struct {
		//Whether the current response processed successfully
		Error bool `json:"error" example:"true"`
		//Error message incase of any error.
		Message string `json:"message" example:"some error information"`
	} `json:"status"`
}

//Bio of the tiger
type PayloadTigerBio struct {
	//Name of the tiger
	Name string `json:"name" form:"name" example:"tiger08"`
	//Date of birth of the tiger. Must be in YYYY-MM-DD format.
	Dob string `json:"birthday" form:"birthday" example:"2005-12-30"`
}

//fields for sighting
type SightingInfo struct {
	//Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.
	LastSeen string `json:"last_seen" form:"last_seen" example:"2019-08-01"`
	//Last seen Latitude point
	Latitude float64 `json:"latitude" form:"latitude" example:"120.51687"`
	//Last seen Longitude point
	Longitude float64 `json:"longitude" form:"longitude" example:"50.894914"`
	ImageId
}

//Parsing Request payload to add a new tiger
type PayloadAddNewTiger struct {
	PayloadTigerBio
	SightingInfo
}

//add sighting payload
type PayloadAddSighting struct {
	SightingInfo
	TigerIdModel
}

//id of the the Tiger
type TigerIdModel struct {
	//id of the tiger
	TigerId int64 `json:"tiger_id" form:"tiger_id" example:"7"`
}

//sighting id field
type SightingIdModel struct {
	//primay key for sighting
	SightingId int64 `json:"sighting_id" form:"sighting_id" example:"12"`
}

//Outgoing Response format for adding a new tiger
type ResponseTiger struct {
	ErrorStatus
	//Data field
	Data struct {
		TigerIdModel
		SightingIdModel
	} `json:"data"`
}

//model for each tiger
type ShowTigerModel struct {
	TigerIdModel
	PayloadTigerBio
	SightingIdModel
	SightingInfo
}

//totals results found response model
type TotalResultsModel struct {
	//totals number of results found for the given query
	Count int `json:"total_results" example:"20"`
}

//Outgoing Response format for show tigers
type ResponseShowTigers struct {
	ErrorStatus
	//Data field
	Data struct {
		TotalResultsModel
		Tigers []ShowTigerModel `json:"tiger_data"`
	} `json:"data"`
}

//Outgoing Response format for show tigers
type ResponseShowSighting struct {
	ErrorStatus
	//Data field
	Data struct {
		TotalResultsModel
		Sightings []SightingInfo `json:"sighting_data"`
	} `json:"data"`
}

//id of the the Tiger
type ImageId struct {
	//Uploaded image id, which can be downloaded at endpoint "/download/image/<generatedimageuuid>"
	Image string `json:"image_id" form:"image_id" example:"generatedimageuuid"`
}
