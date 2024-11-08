package controllers

import (
	"fmt"
	"go-multitenant-boilerplate/db"
	"go-multitenant-boilerplate/db/migrations"
	"go-multitenant-boilerplate/helpers"
	"go-multitenant-boilerplate/middleware"
	tenantmodel "go-multitenant-boilerplate/models/tenant"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *Controller) SignUp(ctx *gin.Context) {
	var tenant tenantmodel.Tenant

	if err := ctx.ShouldBindJSON(&tenant); err != nil {
		ctx.JSON(400, gin.H{"status": false, "message": "Invalid request payload"})
		return
	}

	// Check if email is valid
	if !helpers.CheckValidEmail(tenant.Email) {
		ctx.JSON(400, gin.H{"status": false, "message": "Invalid email address"})
		return
	}

	tenant.Password = helpers.HashPassword(tenant.Password)
	tenant.Dsn = fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), tenant.Organization)

	err := db.CreateSchema(c.TenantDB, tenant.Organization)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	if err := c.TenantService.Create(&tenant); err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.TenantDB.Exec(fmt.Sprintf(`SET search_path="%s"`, tenant.Organization))
	migrate := migrations.MigrateTenantDB(c.TenantDB)
	err = migrate.Migrate()
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}
	err = c.userService.Create(&tenantmodel.User{
		FirstName:    tenant.FirstName,
		LastName:     tenant.LastName,
		Role:         "admin",
		Email:        tenant.Email,
		Password:     tenant.Password,
		Organization: tenant.Organization,
	}, tenant.Organization)
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": true, "message": "Tenant created successfully"})
}

func (c *Controller) Login(ctx *gin.Context) {
	var users tenantmodel.UserLogin
	if err := ctx.ShouldBindJSON(&users); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if users.Email == "" || users.Password == "" || users.Organization == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing credentials"})
		return
	}

	var existingUser tenantmodel.User
	err := c.userService.ReadByEmail(&existingUser, users.Email, users.Organization)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if existingUser.Email == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	if !helpers.VerifyPassword(existingUser.Password, users.Password) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Password is incorrect"})
		return
	}

	jwtToken, err := middleware.GenerateToken(middleware.Claims{
		ID:           existingUser.ID,
		FirstName:    existingUser.FirstName,
		LastName:     existingUser.LastName,
		Organization: existingUser.Organization,
		Email:        existingUser.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": jwtToken, "message": "login successful", "status": true})
}

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
	ctx.JSON(http.StatusOK, gin.H{"message": "user added successfully", "status": true})

}

func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	organization, _ := ctx.Get("organization")

	var user tenantmodel.User

	err := c.userService.GetUserById(&user, id, organization.(string))
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user fetched successfully", "status": true, "data": user})

}

func (c *Controller) ListUsers(ctx *gin.Context) {
	organization, _ := ctx.Get("organization")
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

	err = c.userService.ListUsers(&user, limitInt, offsetInt, organization.(string))
	if err != nil {
		ctx.JSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user fetched successfully", "status": true, "data": user})
}

func (c *Controller) UpdateUser(ctx *gin.Context) {

}

func (c *Controller) DeleteUser(ctx *gin.Context) {

}
