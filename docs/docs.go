// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

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
    "paths": {
        "/api/v1/auth/login": {
            "post": {
                "description": "Login a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.Login.loginRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"message\": \"User logged in successfully\", \"user\": {email: \"email\", organization_id: \"organization_id\"}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "{\"error\": \"Unauthorized\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to login user\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/logout": {
            "post": {
                "description": "Logout a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Logout a user",
                "responses": {
                    "200": {
                        "description": "{\"message\": \"User logged out successfully\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to logout user\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration request",
                        "name": "Creadentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.Register.registerRequestBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "{\"message\": \"User registered successfully\", \"user\": {email: \"\temail\", organization_id: \"organization_id\"}}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to register user\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/validate": {
            "get": {
                "description": "Validate a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Validate a user",
                "responses": {
                    "200": {
                        "description": "{\"message\": \"User logged in successfully\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to validate user\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/organization/": {
            "get": {
                "description": "Get organization information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "summary": "Get organization information",
                "responses": {
                    "200": {
                        "description": "{\"organizations\": \"organizations\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to get organization\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Update organization information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "summary": "Update organization information",
                "parameters": [
                    {
                        "description": "Organization",
                        "name": "Organization",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.organizationBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"organizations\": \"organizations\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to get organization\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/organization/list-project": {
            "get": {
                "description": "List projects",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "summary": "List projects",
                "responses": {
                    "200": {
                        "description": "{\"projects\": \"projects\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to list projects\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/organization/new-project": {
            "post": {
                "description": "Create new project for organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "summary": "Create new project",
                "parameters": [
                    {
                        "description": "Project",
                        "name": "Project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.projectBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"project\": \"project\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to create project\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/organization/{project_id}": {
            "delete": {
                "description": "Delete project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Organization"
                ],
                "summary": "Delete project",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"message\": \"Project deleted\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to delete project\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/project_detail/{project_id}": {
            "get": {
                "description": "Get project overview",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project Detail / Overview"
                ],
                "summary": "Get project overview",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "project_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"project\": \"project\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to get project detail\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/projects/{project_id}": {
            "put": {
                "description": "Update project information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project Detail / Overview"
                ],
                "summary": "Update project information",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "project_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Project",
                        "name": "Project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.projectBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"project\": \"project\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{\"error\": \"Bad request\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "{\"error\": \"Failed to update project\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Login.loginRequestBody": {
            "type": "object",
            "required": [
                "email",
                "organization_name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "organization_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "controllers.Register.registerRequestBody": {
            "type": "object",
            "required": [
                "email",
                "organization_name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "organization_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "controllers.organizationBody": {
            "type": "object",
            "properties": {
                "aliasName": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "establishmentDate": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.projectBody": {
            "type": "object",
            "properties": {
                "currentSprint": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "repoURL": {
                    "type": "string"
                },
                "startAt": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
