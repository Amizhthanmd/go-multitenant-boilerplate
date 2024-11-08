package controllers

import (
	"fmt"
	"go-multitenant-boilerplate/db"
	"go-multitenant-boilerplate/db/migrations"
	"go-multitenant-boilerplate/helpers"
	tenantmodel "go-multitenant-boilerplate/models/tenant"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		ID:        uuid.NewString(),
		FirstName: tenant.FirstName,
		LastName:  tenant.LastName,
		Role:      "admin",
		Email:     tenant.Email,
		Password:  tenant.Password,
	})
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

	var user tenantmodel.User
	err := c.userService.ReadByEmail(&user, users.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if user.Email == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
	}

}

func (c *Controller) AddUser(ctx *gin.Context) {

}

func (c *Controller) GetUsers(ctx *gin.Context) {

}

func (c *Controller) UpdateUser(ctx *gin.Context) {

}

func (c *Controller) DeleteUser(ctx *gin.Context) {

}
