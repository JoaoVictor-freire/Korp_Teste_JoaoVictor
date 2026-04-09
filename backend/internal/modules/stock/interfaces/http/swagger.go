package http

const stockSwaggerDoc = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Stock Service API",
    "version": "1.0.0",
    "description": "Product and auth endpoints for the Korp backend challenge."
  },
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "HealthResponse": {
        "type": "object",
        "properties": {
          "service": {"type": "string", "example": "stock-service"},
          "status": {"type": "string", "example": "ok"}
        }
      },
      "ErrorEnvelope": {
        "type": "object",
        "properties": {
          "error": {
            "type": "object",
            "properties": {
              "message": {"type": "string"}
            }
          }
        }
      },
      "RegisterRequest": {
        "type": "object",
        "required": ["name", "email", "password"],
        "properties": {
          "name": {"type": "string", "example": "Joao Victor"},
          "email": {"type": "string", "example": "joao@example.com"},
          "password": {"type": "string", "example": "123456"}
        }
      },
      "LoginRequest": {
        "type": "object",
        "required": ["email", "password"],
        "properties": {
          "email": {"type": "string", "example": "joao@example.com"},
          "password": {"type": "string", "example": "123456"}
        }
      },
      "User": {
        "type": "object",
        "properties": {
          "id": {"type": "string"},
          "name": {"type": "string"},
          "email": {"type": "string"},
          "created_at": {"type": "string", "format": "date-time"}
        }
      },
      "AuthResponseEnvelope": {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "properties": {
              "user": {"$ref": "#/components/schemas/User"},
              "token": {"type": "string"}
            }
          }
        }
      },
      "CreateProductRequest": {
        "type": "object",
        "required": ["code", "description", "stock"],
        "properties": {
          "code": {"type": "string", "example": "P001"},
          "description": {"type": "string", "example": "Caneta Azul"},
          "stock": {"type": "integer", "example": 10}
        }
      },
      "Product": {
        "type": "object",
        "properties": {
          "owner_id": {"type": "string"},
          "code": {"type": "string"},
          "description": {"type": "string"},
          "stock": {"type": "integer"},
          "created_at": {"type": "string", "format": "date-time"}
        }
      },
      "ProductEnvelope": {
        "type": "object",
        "properties": {
          "data": {"$ref": "#/components/schemas/Product"}
        }
      },
      "ProductListEnvelope": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": {"$ref": "#/components/schemas/Product"}
          }
        }
      }
    }
  },
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "Service health",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/HealthResponse"}
              }
            }
          }
        }
      }
    },
    "/whoami": {
      "get": {
        "summary": "Auth usage hint",
        "responses": {
          "200": {
            "description": "Usage hint",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "hint": {"type": "string"}
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/v1/auth/register": {
      "post": {
        "summary": "Register user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/RegisterRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "User registered",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/AuthResponseEnvelope"}
              }
            }
          },
          "400": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      }
    },
    "/api/v1/auth/login": {
      "post": {
        "summary": "Login user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/LoginRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login success",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/AuthResponseEnvelope"}
              }
            }
          },
          "401": {
            "description": "Invalid credentials",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      }
    },
    "/api/v1/products": {
      "get": {
        "summary": "List products",
        "security": [{"BearerAuth": []}],
        "responses": {
          "200": {
            "description": "Products list",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ProductListEnvelope"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create product",
        "security": [{"BearerAuth": []}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreateProductRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Product created",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ProductEnvelope"}
              }
            }
          },
          "400": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          },
          "401": {
            "description": "Unauthorized",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          },
          "409": {
            "description": "Conflict",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      }
    }
  }
}`
