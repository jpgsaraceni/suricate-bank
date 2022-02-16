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
        "contact": {
            "name": "João Saraceni",
            "url": "https://www.linkedin.com/in/joaosaraceni/",
            "email": "jpgome@id.uff.br"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/jpgsaraceni/suricate-bank/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/accounts": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Get all accounts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/schemas.FetchedAccount"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates an account with BRL$10.00 initial balance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Create a new account",
                "parameters": [
                    {
                        "description": "Account",
                        "name": "account",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.CreateAccountRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Idempotency key",
                        "name": "Idempotency-Key",
                        "in": "header"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.CreateAccountResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    }
                }
            }
        },
        "/accounts/{id}/balance": {
            "get": {
                "description": "Gets the balance of the account matching account ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "Get account balance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.GetBalanceResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Login"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Login Credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schemas.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    }
                }
            }
        },
        "/transfers": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transfer"
                ],
                "summary": "Get all transfers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/schemas.FetchedTransfer"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a transfer from origin account matching bearer token account id\nto account with request body account ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transfer"
                ],
                "summary": "Create transfer",
                "parameters": [
                    {
                        "description": "Transfer",
                        "name": "transfer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.CreateTransferRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Idempotency key",
                        "name": "Idempotency-Key",
                        "in": "header"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/schemas.CreateTransferResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorPayload"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.ErrorPayload": {
            "type": "object",
            "properties": {
                "title": {
                    "type": "string",
                    "example": "Message for some error"
                }
            }
        },
        "schemas.CreateAccountRequest": {
            "type": "object",
            "properties": {
                "cpf": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "secret": {
                    "type": "string"
                }
            }
        },
        "schemas.CreateAccountResponse": {
            "type": "object",
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "balance": {
                    "type": "string"
                },
                "cpf": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "schemas.CreateTransferRequest": {
            "type": "object",
            "properties": {
                "account_destination_id": {
                    "type": "string"
                },
                "amount": {
                    "type": "integer"
                }
            }
        },
        "schemas.CreateTransferResponse": {
            "type": "object",
            "properties": {
                "account_destination_id": {
                    "type": "string"
                },
                "account_origin_id": {
                    "type": "string"
                },
                "amount": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "transfer_id": {
                    "type": "string"
                }
            }
        },
        "schemas.FetchedAccount": {
            "type": "object",
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "balance": {
                    "type": "string"
                },
                "cpf": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "schemas.FetchedTransfer": {
            "type": "object",
            "properties": {
                "account_destination_id": {
                    "type": "string"
                },
                "account_origin_id": {
                    "type": "string"
                },
                "amount": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "transfer_id": {
                    "type": "string"
                }
            }
        },
        "schemas.GetBalanceResponse": {
            "type": "object",
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "balance": {
                    "type": "string"
                }
            }
        },
        "schemas.LoginRequest": {
            "type": "object",
            "properties": {
                "cpf": {
                    "type": "string"
                },
                "secret": {
                    "type": "string"
                }
            }
        },
        "schemas.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Access token": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.2.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Suricate Bank API",
	Description:      "Suricate Bank is an api that creates accounts and transfers money between them.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}