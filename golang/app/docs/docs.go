// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://www.example.com/terms",
        "contact": {
            "name": "Tech Support",
            "email": "rg@dooglex.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/sighting/add": {
            "post": {
                "description": "Create a new sighting of existing tiger",
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tiger"
                ],
                "summary": "Create a new sighting of existing tiger",
                "operationId": "create_sighting",
                "parameters": [
                    {
                        "description": "Request payload",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.PayloadAddSighting"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseTiger"
                        }
                    }
                }
            }
        },
        "/sighting/show": {
            "post": {
                "description": "show the list of sightings of tigers",
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tiger"
                ],
                "summary": "show the list of sightings of tigers",
                "operationId": "show_sighting",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Page number. Default: 1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "description": "Request payload",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.TigerIdModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseShowSighting"
                        }
                    }
                }
            }
        },
        "/tiger/add": {
            "post": {
                "description": "Create a new tiger along with the last seen info",
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tiger"
                ],
                "summary": "Create a new tiger along with the last seen info",
                "operationId": "create_tiger",
                "parameters": [
                    {
                        "description": "Request payload",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.PayloadAddNewTiger"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseTiger"
                        }
                    }
                }
            }
        },
        "/tiger/show": {
            "post": {
                "description": "show the list of tigers sorted by last seen time",
                "consumes": [
                    "application/json",
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tiger"
                ],
                "summary": "show the list of tigers sorted by last seen time",
                "operationId": "show_tigers",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Page number. Default: 1",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseShowTigers"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.PayloadAddNewTiger": {
            "type": "object",
            "properties": {
                "birthday": {
                    "description": "Date of birth of the tiger. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "image": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "last_seen": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "latitude": {
                    "description": "Last seen Latitude point",
                    "type": "number"
                },
                "longitude": {
                    "description": "Last seen Longitude point",
                    "type": "number"
                },
                "name": {
                    "description": "Name of the tiger",
                    "type": "string"
                }
            }
        },
        "main.PayloadAddSighting": {
            "type": "object",
            "properties": {
                "image": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "last_seen": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "latitude": {
                    "description": "Last seen Latitude point",
                    "type": "number"
                },
                "longitude": {
                    "description": "Last seen Longitude point",
                    "type": "number"
                },
                "tiger_id": {
                    "description": "id of the tiger",
                    "type": "integer"
                }
            }
        },
        "main.ResponseShowSighting": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Data field",
                    "type": "object",
                    "properties": {
                        "sighting_data": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.SightingInfo"
                            }
                        },
                        "total_results": {
                            "description": "totals number of results found for the given query",
                            "type": "integer"
                        }
                    }
                },
                "status": {
                    "description": "status of the error occurence in the current response",
                    "type": "object",
                    "properties": {
                        "error": {
                            "description": "Whether the current response processed successfully",
                            "type": "boolean"
                        },
                        "message": {
                            "description": "Error message incase of any error.",
                            "type": "string"
                        }
                    }
                }
            }
        },
        "main.ResponseShowTigers": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Data field",
                    "type": "object",
                    "properties": {
                        "tiger_data": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.ShowTigerModel"
                            }
                        },
                        "total_results": {
                            "description": "totals number of results found for the given query",
                            "type": "integer"
                        }
                    }
                },
                "status": {
                    "description": "status of the error occurence in the current response",
                    "type": "object",
                    "properties": {
                        "error": {
                            "description": "Whether the current response processed successfully",
                            "type": "boolean"
                        },
                        "message": {
                            "description": "Error message incase of any error.",
                            "type": "string"
                        }
                    }
                }
            }
        },
        "main.ResponseTiger": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Data field",
                    "type": "object",
                    "properties": {
                        "sighting_id": {
                            "description": "primay key for sighting",
                            "type": "integer"
                        },
                        "tiger_id": {
                            "description": "id of the tiger",
                            "type": "integer"
                        }
                    }
                },
                "status": {
                    "description": "status of the error occurence in the current response",
                    "type": "object",
                    "properties": {
                        "error": {
                            "description": "Whether the current response processed successfully",
                            "type": "boolean"
                        },
                        "message": {
                            "description": "Error message incase of any error.",
                            "type": "string"
                        }
                    }
                }
            }
        },
        "main.ShowTigerModel": {
            "type": "object",
            "properties": {
                "birthday": {
                    "description": "Date of birth of the tiger. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "image": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "last_seen": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "latitude": {
                    "description": "Last seen Latitude point",
                    "type": "number"
                },
                "longitude": {
                    "description": "Last seen Longitude point",
                    "type": "number"
                },
                "name": {
                    "description": "Name of the tiger",
                    "type": "string"
                },
                "sighting_id": {
                    "description": "primay key for sighting",
                    "type": "integer"
                },
                "tiger_id": {
                    "description": "id of the tiger",
                    "type": "integer"
                }
            }
        },
        "main.SightingInfo": {
            "type": "object",
            "properties": {
                "image": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "last_seen": {
                    "description": "Timestamp when the tiger was last seen. Must be in YYYY-MM-DD format.",
                    "type": "string"
                },
                "latitude": {
                    "description": "Last seen Latitude point",
                    "type": "number"
                },
                "longitude": {
                    "description": "Last seen Longitude point",
                    "type": "number"
                }
            }
        },
        "main.TigerIdModel": {
            "type": "object",
            "properties": {
                "tiger_id": {
                    "description": "id of the tiger",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo_swagger holds exported Swagger Info so clients can modify it
var SwaggerInfo_swagger = &swag.Spec{
	Version:          "1.0",
	Host:             "tigerhall.dooglex.com",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Tigerhall test API",
	Description:      "This is an swagger documentation of simple test API task given by tigerhall",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo_swagger.InstanceName(), SwaggerInfo_swagger)
}