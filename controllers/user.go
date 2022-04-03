package controllers

import (
	"finalproject/models"
	"finalproject/utils"
	"finalproject/utils/token"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type users struct {
	FullName   string `json:"full_name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password"`
	AdminToken string `json:"admin_token"`
}

// GetAllUser godoc
// @Summary Get All User
// @Description Get a list of User
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Success 200 {object} []models.Guest
// @Router /users [get]
func GetAllUser(c *gin.Context) {
	//get db form gin context
	db := c.MustGet("db").(*gorm.DB)
	var users []models.Guest
	db.Table("users").Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser godoc
// @Summary Create New User
// @Description Creating a new User
// @Tags User
// @Param Body body users true "the body to create a new User"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input users
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		if input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be filled!"})
			return
		}
	}

	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errPassword.Error()})
		return
	}

	user := models.User{FullName: input.FullName,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "guest",
	}
	if input.AdminToken == utils.Getenv("ADMINISTRATOR", "JCCGOLANG") {
		user.Role = "admin"
	}
	db := c.MustGet("db").(*gorm.DB)
	var err error = db.Create(&user).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

//GetUserById godoc
// @Summary Get User
// @Description Get A User by id
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "user id"
// @Success  200 {object} models.Guest
// @Router /users/{id} [get]
func GetUserById(c *gin.Context) {
	var user models.Guest
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Table("users").Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID not found!"})
		return
	}
	db.Preload("Posts").Find(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser godoc
// @Summary Update User.
// @Description Update user by id.
// @Tags User
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "user id"
// @Param Body body users true "the body to update an user"
// @Success  200 {object} models.User
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//Get model if exist
	var user models.User

	UserIdCurrent, errToken := token.ExtractTokenID(c)

	if errToken != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login first"})
		return
	}
	roleToken, errToken := token.ExtractTokenRole(c)
	if errToken != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error()})
		return
	}
	if roleToken == "admin" {
		fmt.Println("Admin")
		if err := db.Where("id = ? ", c.Param("id")).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID not found!"})
			return
		}
	} else {
		if err := db.Where("id = ? ", UserIdCurrent).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID not found!"})
			return
		}
	}

	var input users
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUser := models.User{FullName: input.FullName,
		Username:  input.Username,
		Email:     input.Email,
		UpdatedAt: time.Now(),
	}

	var err error = db.Model(&user).Updates(updateUser).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if input.Password != "" {
		// password hash
		hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if errPassword != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errPassword.Error()})
			return
		}
		updatePassword := models.User{
			Password: string(hashedPassword),
		}
		var err error = db.Model(&user).Updates(updatePassword).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

//DeleteUser godoc
// @Summary Delete one user.
// @Description Delete a user by id.
// @Tags User
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} map[string]boolean
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	//Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var user models.User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data not found!"})
		return
	}
	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
