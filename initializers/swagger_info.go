package initializers

import (
	"os"

	"github.com/swaggo/swag"
)

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {}
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost" + os.Getenv("PORT"),
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Parameter Store BE API",
	Description:      "This is the API for the Parameter Store BE application. It is a RESTful API that allows you to manage parameter in organization",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}
