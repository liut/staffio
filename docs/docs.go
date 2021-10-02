// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/api/login": {
            "post": {
                "description": "login",
                "consumes": [
                    "application/x-www-form-urlencoded",
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "login",
                "operationId": "api-1-login-post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apis.RespDone"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    }
                }
            }
        },
        "/api/password": {
            "post": {
                "description": "change password",
                "consumes": [
                    "application/x-www-form-urlencoded",
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Change password",
                "operationId": "api-1-password-post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Old Password",
                        "name": "old_password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "New Password",
                        "name": "new_password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Confirm Password",
                        "name": "password_confirm",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apis.RespDone"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    }
                }
            }
        },
        "/api/password/forgot": {
            "post": {
                "description": "forgot password",
                "consumes": [
                    "application/x-www-form-urlencoded",
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Forgot password",
                "operationId": "api-1-password-forgot-post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Login name",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Mobile number",
                        "name": "mobile",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apis.RespDone"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    }
                }
            }
        },
        "/api/password/reset": {
            "post": {
                "description": "reset password, form:rt, json:token",
                "consumes": [
                    "application/x-www-form-urlencoded",
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Reset password",
                "operationId": "api-1-password-reset-post",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Login name",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Confirm Password",
                        "name": "password_confirm",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "rt",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apis.RespDone"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apis.RespFail"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apis.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "apis.RespDone": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "main data"
                },
                "extra": {
                    "description": "extra data"
                },
                "ok": {
                    "description": "OK",
                    "type": "boolean"
                }
            }
        },
        "apis.RespFail": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/apis.Error"
                },
                "ok": {
                    "type": "boolean"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
