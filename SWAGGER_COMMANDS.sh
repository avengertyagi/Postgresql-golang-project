#!/bin/bash
# Quick Reference Commands for Swagger Documentation

# ============================================================================
# BUILD AND RUN THE API
# ============================================================================

# Build the project
# go build -o api cmd/api/main.go

# Run the API server
# go run ./cmd/api/main.go

# ============================================================================
# REGENERATE SWAGGER DOCS (After modifying endpoints/comments)
# ============================================================================

# From project root:
# go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go

# ============================================================================
# ACCESS SWAGGER UI
# ============================================================================

# Once API is running, visit:
# - http://localhost:8080/swagger/index.html (Recommended)
# - http://localhost:8080/docs (Redirects to above)
# - http://localhost:8080/swagger/doc.json (Raw JSON)

# ============================================================================
# COMMON SWAGGER ANNOTATION TAGS
# ============================================================================

# @Summary       - Brief description (required)
# @Description   - Longer description (optional)
# @Tags          - Endpoint category
# @Accept        - Input content types (json, xml, form)
# @Produce       - Output content types (json, xml)
# @Param         - Parameter documentation
# @Success       - Success response
# @Failure       - Error response
# @Security      - Auth requirement (e.g., BearerAuth)
# @Router        - Route path and HTTP method

# ============================================================================
# PARAMETER TYPES
# ============================================================================

# Query parameter:
# @Param id query int true "User ID"

# Path parameter:
# @Param id path int true "User ID"

# Body parameter:
# @Param request body models.User true "User data"

# ============================================================================
# RESPONSE TYPES
# ============================================================================

# Success response with object:
# @Success 200 {object} models.User

# Success response with array:
# @Success 200 {array} models.User

# Success response with map:
# @Success 200 {object} map[string]interface{}

# ============================================================================
# QUICK EXAMPLE
# ============================================================================

# Example endpoint documentation:
#
# // @Summary      Get user by ID
# // @Description  Retrieve a specific user's information
# // @Tags         Users
# // @Accept       json
# // @Produce      json
# // @Param        id path int true "User ID"
# // @Success      200 {object} models.User
# // @Failure      404 {object} map[string]interface{}
# // @Router       /users/{id} [get]
# func GetUser(c *gin.Context) {
#     // implementation
# }

# ============================================================================
# HELPFUL LINKS
# ============================================================================

# Swagger/OpenAPI Docs: https://swagger.io/docs/specification/
# Swaggo GitHub: https://github.com/swaggo/swag
# Gin-Swagger GitHub: https://github.com/swaggo/gin-swagger
