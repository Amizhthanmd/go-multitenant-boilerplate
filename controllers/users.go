package controllers

import (
	"go-multitenant-boilerplate/helpers"
	tenantmodel "go-multitenant-boilerplate/models/tenant"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *Controller) AddUser(ctx *gin.Context) {
	var users tenantmodel.User
	if err := ctx.ShouldBindJSON(&users); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	users.Password = helpers.HashPassword(users.Password)

	err := c.userService.Create(&tenantmodel.User{
		FirstName:    users.FirstName,
		LastName:     users.LastName,
		Role:         users.Role,
		Email:        users.Email,
		Password:     users.Password,
		Organization: users.Organization,
	}, users.Organization)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "user added successfully"})

}

func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	organization := ctx.GetString("organization")

	var user tenantmodel.User

	err := c.userService.GetUserById(&user, id, organization)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "user fetched successfully", "data": user})

}

func (c *Controller) ListUsers(ctx *gin.Context) {
	organization := ctx.GetString("organization")
	var user []tenantmodel.User

	limit := ctx.DefaultQuery("limit", "10")
	offset := ctx.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid limit"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid offset"})
		return
	}

	err = c.userService.ListUsers(&user, limitInt, offsetInt, organization)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "user fetched successfully", "data": user})
}

func (c *Controller) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	organization := ctx.GetString("organization")

	var user tenantmodel.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user.ID = id

	err := c.userService.UpdateUser(&user, organization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "User updated successfully", "data": user})
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	organization := ctx.GetString("organization")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user ID"})
		return
	}

	err := c.userService.DeleteUser(id, organization)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "User deleted successfully"})
}
