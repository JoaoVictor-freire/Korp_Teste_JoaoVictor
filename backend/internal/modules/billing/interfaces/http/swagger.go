package http

const billingSwaggerDoc = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Billing Service API",
    "version": "1.0.0",
    "description": "Invoice endpoints for the Korp backend challenge."
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
          "service": {"type": "string", "example": "billing-service"},
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
      "InvoiceItem": {
        "type": "object",
        "required": ["product_code", "quantity"],
        "properties": {
          "product_code": {"type": "string", "example": "P001"},
          "quantity": {"type": "integer", "example": 2}
        }
      },
      "CreateInvoiceRequest": {
        "type": "object",
        "required": ["number", "items"],
        "properties": {
          "number": {"type": "integer", "example": 1},
          "items": {
            "type": "array",
            "items": {"$ref": "#/components/schemas/InvoiceItem"}
          }
        }
      },
      "Invoice": {
        "type": "object",
        "properties": {
          "owner_id": {"type": "string"},
          "number": {"type": "integer"},
          "status": {"type": "string", "example": "OPEN"},
          "items": {
            "type": "array",
            "items": {"$ref": "#/components/schemas/InvoiceItem"}
          },
          "created_at": {"type": "string", "format": "date-time"}
        }
      },
      "InvoiceEnvelope": {
        "type": "object",
        "properties": {
          "data": {"$ref": "#/components/schemas/Invoice"}
        }
      },
      "InvoiceListEnvelope": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": {"$ref": "#/components/schemas/Invoice"}
          }
        }
      },
      "MessageEnvelope": {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "properties": {
              "message": {"type": "string", "example": "invoice closed successfully"}
            }
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
    "/api/v1/invoices": {
      "get": {
        "summary": "List invoices",
        "security": [{"BearerAuth": []}],
        "responses": {
          "200": {
            "description": "Invoices list",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/InvoiceListEnvelope"}
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
        "summary": "Create invoice",
        "security": [{"BearerAuth": []}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreateInvoiceRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Invoice created",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/InvoiceEnvelope"}
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
    },
    "/api/v1/invoices/{number}": {
      "get": {
        "summary": "Get invoice by number",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "name": "number",
            "in": "path",
            "required": true,
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Invoice details",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/InvoiceEnvelope"}
              }
            }
          },
          "404": {
            "description": "Invoice not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      },
      "put": {
        "summary": "Update invoice",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "name": "number",
            "in": "path",
            "required": true,
            "schema": {"type": "integer"}
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreateInvoiceRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Invoice updated",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/InvoiceEnvelope"}
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
          "404": {
            "description": "Invoice not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          },
          "409": {
            "description": "Closed invoice or duplicate number",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      },
      "delete": {
        "summary": "Delete invoice",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "name": "number",
            "in": "path",
            "required": true,
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Invoice deleted",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/MessageEnvelope"}
              }
            }
          },
          "404": {
            "description": "Invoice not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          },
          "409": {
            "description": "Closed invoice cannot be deleted",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          }
        }
      }
    },
    "/api/v1/invoices/{number}/close": {
      "patch": {
        "summary": "Close invoice",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {
            "name": "number",
            "in": "path",
            "required": true,
            "schema": {"type": "integer"}
          }
        ],
        "responses": {
          "200": {
            "description": "Invoice closed",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/MessageEnvelope"}
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
          "404": {
            "description": "Invoice not found",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorEnvelope"}
              }
            }
          },
          "409": {
            "description": "Invoice already closed",
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
