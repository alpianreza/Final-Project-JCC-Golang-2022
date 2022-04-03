package controllers

import (
	"finalproject/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type categoryInput struct {
	Category string `json:"category_name" binding:"required"`
}

// GetAllCategory godoc
// @Summary Get All Category
// @Description Get a list of Category
// @Tags Category
// @Produce json
// @Success 200 {object} []models.Category
// @Router /category [get]
func GetAllCategory(c *gin.Context) {
	//get db form gin context
	db := c.MustGet("db").(*gorm.DB)
	var category []models.Category
	db.Find(&category)
	c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateCategory godoc
// @Summary Create New Category
// @Description Creating a new Category
// @Tags Category
// @Param Body body categoryInput true "the body to create a new Category"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Category
// @Router /category [post]
func CreateCategory(c *gin.Context) {
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create Category
	category := models.Category{Category: input.Category,
		CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	var err error = db.Create(&category).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

//GetCategoryById godoc
// @Summary Get Category
// @Description Get A Category by id
// @Tags Category
// @Produce json
// @Param id path string true "category id"
// @Success  200 {object} models.Category
// @Router /category/{id} [get]
func GetCategoryById(c *gin.Context) {
	var category models.Category
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}
	db.Preload("Posts").Find(&category)

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// UpdateCategory godoc
// @Summary Update Category.
// @Description Update category by id.
// @Tags Category
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "category id"
// @Param Body body categoryInput true "the body to update an category"
// @Success  200 {object} models.Category
// @Router /category/{id} [patch]
func UpdateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//Get model if exist
	var category models.Category

	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}

	//Validate Input
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedInput := models.Category{Category: input.Category,
		UpdatedAt: time.Now(),
	}

	var err error = db.Model(&category).Updates(updatedInput).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

//DeleteCategory godoc
// @Summary Delete one category.
// @Description Delete a category by id.
// @Tags Category
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "category id"
// @Success 200 {object} map[string]boolean
// @Router /category/{id} [delete]
func DeleteCategory(c *gin.Context) {
	//Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category

	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}
	db.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
