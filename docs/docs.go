// Package docs provides Swagger documentation.
// Para regenerar ejecuta: swag init -g main.go
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "title": "Coleccion Backend API",
        "description": "API para gestión de vehículos de colección con imágenes en Google Drive.",
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "paths": {
        "/auth/login": {
            "post": {
                "tags": ["auth"],
                "summary": "Iniciar sesión",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/LoginRequest"}}],
                "responses": {"200": {"description": "OK"}, "401": {"description": "Unauthorized"}}
            }
        },
        "/auth/logout": {
            "post": {
                "tags": ["auth"],
                "summary": "Cerrar sesión",
                "security": [{"BearerAuth": []}],
                "responses": {"204": {"description": "No Content"}}
            }
        },
        "/auth/refresh": {
            "post": {
                "tags": ["auth"],
                "summary": "Renovar token",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "parameters": [{"in": "body", "name": "body", "required": true, "schema": {"$ref": "#/definitions/RefreshRequest"}}],
                "responses": {"200": {"description": "OK"}, "401": {"description": "Unauthorized"}}
            }
        },
        "/coleccion": {
            "get": {
                "tags": ["coleccion"],
                "summary": "Listar vehículos de colección",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "responses": {"200": {"description": "OK"}}
            },
            "post": {
                "tags": ["coleccion"],
                "summary": "Crear vehículo de colección",
                "security": [{"BearerAuth": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {"in": "formData", "name": "nombre", "type": "string", "required": true, "description": "Nombre del vehículo"},
                    {"in": "formData", "name": "marca", "type": "string", "required": true, "description": "Marca del vehículo"},
                    {"in": "formData", "name": "modelo", "type": "string", "required": true, "description": "Modelo del vehículo"},
                    {"in": "formData", "name": "imagenes", "type": "file", "description": "Imagen (máx 3, repetir parámetro)"}
                ],
                "responses": {"201": {"description": "Created"}, "400": {"description": "Bad Request"}}
            }
        },
        "/coleccion/{id}": {
            "get": {
                "tags": ["coleccion"],
                "summary": "Obtener vehículo por ID",
                "security": [{"BearerAuth": []}],
                "produces": ["application/json"],
                "parameters": [{"in": "path", "name": "id", "type": "string", "required": true}],
                "responses": {"200": {"description": "OK"}, "404": {"description": "Not Found"}}
            },
            "put": {
                "tags": ["coleccion"],
                "summary": "Actualizar vehículo",
                "security": [{"BearerAuth": []}],
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "parameters": [
                    {"in": "path", "name": "id", "type": "string", "required": true},
                    {"in": "formData", "name": "nombre", "type": "string"},
                    {"in": "formData", "name": "marca", "type": "string"},
                    {"in": "formData", "name": "modelo", "type": "string"},
                    {"in": "formData", "name": "imagenes", "type": "file", "description": "Nuevas imágenes (reemplazan las anteriores)"}
                ],
                "responses": {"200": {"description": "OK"}, "404": {"description": "Not Found"}}
            },
            "delete": {
                "tags": ["coleccion"],
                "summary": "Eliminar vehículo",
                "security": [{"BearerAuth": []}],
                "parameters": [{"in": "path", "name": "id", "type": "string", "required": true}],
                "responses": {"204": {"description": "No Content"}, "404": {"description": "Not Found"}}
            }
        }
    },
    "definitions": {
        "LoginRequest": {
            "type": "object",
            "required": ["email", "password"],
            "properties": {
                "email": {"type": "string", "example": "usuario@email.com"},
                "password": {"type": "string", "example": "contraseña123"}
            }
        },
        "RefreshRequest": {
            "type": "object",
            "required": ["refresh_token"],
            "properties": {
                "refresh_token": {"type": "string"}
            }
        }
    }
}`

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: swag.Name,
		SwaggerTemplate:  docTemplate,
	})
}
