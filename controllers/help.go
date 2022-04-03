package controllers

import (
	"finalproject/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HelpInput struct {
	Content string `json:"content" binding:"required"`
}

// GetAllHelp godoc
// @Summary Get All Help
// @Description Get a list of Help
// @Tags Help
// @Produce json
// @Success 200 {object} []models.Help
// @Router /help [get]
func GetAllHelp(c *gin.Context) {
	//get db form gin context
	db := c.MustGet("db").(*gorm.DB)
	var Help []models.Help
	db.Find(&Help)
	c.JSON(http.StatusOK, gin.H{"data": Help})
}

// CreateHelp godoc
// @Summary Create New Help
// @Description Creating a new Help
// @Tags Help
// @Param Body body HelpInput true "the body to create a new Help"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Help
// @Router /help [post]
func CreateHelp(c *gin.Context) {
	var input HelpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create Help
	help := models.Help{Content: input.Content,
		CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	var err error = db.Create(&help).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": help})
}

//GetHelpById godoc
// @Summary Get Help
// @Description Get A Help by id
// @Tags Help
// @Produce json
// @Param id path string true "Help id"
// @Success  200 {object} models.Help
// @Router /help/{id} [get]
func GetHelpById(c *gin.Context) {
	var Help models.Help
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", c.Param("id")).First(&Help).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}
	db.Preload("Posts").Find(&Help)

	c.JSON(http.StatusOK, gin.H{"data": Help})
}

// UpdateHelp godoc
// @Summary Update Help.
// @Description Update Help by id.
// @Tags Help
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Help id"
// @Param Body body HelpInput true "the body to update an Help"
// @Success  200 {object} models.Help
// @Router /help/{id} [patch]
func UpdateHelp(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//Get model if exist
	var Help models.Help

	if err := db.Where("id = ?", c.Param("id")).First(&Help).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}

	//Validate Input
	var input HelpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateHelp := models.Help{Content: input.Content,
		UpdatedAt: time.Now(),
	}

	var err error = db.Model(&Help).Updates(updateHelp).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": Help})
}

//DeleteHelp godoc
// @Summary Delete one Help.
// @Description Delete a Help by id.
// @Tags Help
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Help id"
// @Success 200 {object} map[string]boolean
// @Router /help/{id} [delete]
func DeleteHelp(c *gin.Context) {
	//Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var Help models.Help

	if err := db.Where("id = ?", c.Param("id")).First(&Help).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}
	db.Delete(&Help)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
