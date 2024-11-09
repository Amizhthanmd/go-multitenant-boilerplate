package controllers

import (
	"fmt"
	tenantmodel "go-multitenant-boilerplate/models/tenant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) ListPermissions(ctx *gin.Context) {
	organization := ctx.GetString("organization")
	var permissions []tenantmodel.Permission

	err := c.userService.ListPermissions(&permissions, organization)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "permissions fetched successfully", "data": permissions})
}

func (c *Controller) AddRoles(ctx *gin.Context) {
	organization := ctx.GetString("organization")

	var roles tenantmodel.AddRoles
	if err := ctx.ShouldBindJSON(&roles); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	c.TenantDB.Exec(fmt.Sprintf(`SET search_path="%s"`, organization))

	var permissions []tenantmodel.Permission
	err := c.userService.GetPermissionsByIds(&permissions, roles.Permissions)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	err = c.userService.CreateRoles(&tenantmodel.Role{
		Name:        roles.Name,
		Permissions: permissions,
	})
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "Roles created successfully", "data": roles})
}

func (c *Controller) ListRoles(ctx *gin.Context) {

}

func (c *Controller) UpdateRoles(ctx *gin.Context) {

}

func (c *Controller) DeleteRoles(ctx *gin.Context) {

}
