package controllers

import (
	"finalproject/models"
	"finalproject/utils/token"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type commentInput struct {
	CommentPost string            `json:"post_comment" binding:"required" gorm:"unique"`
	Publish     bool              `json:"publish" binding:"required"`
	CommentMeta map[string]string `json:"post_meta" binding:"required"`
}

// GetAllComments godoc
// @Summary Get All Comments
// @Description Get a list of User
// @Tags Comments
// @Produce json
// @Success 200 {object} []models.Comments
// @Router /comments [get]
func GetAllComment(c *gin.Context) {
	//get db from fin context
	db := c.MustGet("db").(*gorm.DB)

	var comments []models.Comments
	db.Find(&comments)
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

// CreateComments godoc
// @Summary Create New Comments
// @Description Creating a new Comments
// @Tags Comments
// @Param Body body commentInput true "the body to create a new Comments"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} models.Comments
// @Router /users [post]
func CreateComment(c *gin.Context) {
	var input commentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UserIdCurrent, errToken := token.ExtractTokenID(c)

	if errToken != nil {
		fmt.Println(errToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login to continue"})
		return
	}

	inputComment := models.Comments{PostComment: input.CommentPost,
		Publish:   input.Publish,
		UserID:    UserIdCurrent,
		CreatedAt: time.Now(),
	}

	db := c.MustGet("db").(*gorm.DB)
	_, err := inputComment.CommentSave(input.CommentMeta, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": inputComment})
}

// GetCommentById godoc
// @Summary Get Comments
// @Description Get Comment by Id
// @Tags Comments
// @Produce json
// @Param id path string trur "comments id"
// @Success 200 {object} models.Comments
// @Router /comments/{id} [get]
func GetCommentById(c *gin.Context) {
	var comments models.Comments
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", c.Param("id")).First(&comments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}

	db.Preload("User").Preload("Category").Preload("Comments").Preload("Meta").Find(&comments)
	c.JSON(http.StatusOK, gin.H{"data": comments})
}

//GetCommentsById godoc
// @Summary Get Comments
// @Description Get A Comments by id
// @Tags Comments
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param id path string true "comment id"
// @Success  200 {object} models.Comments
// @Router /comments/{id} [get]
func UpdateComment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//Get model if exist
	var Comments models.Comments
	UserIdCurrent, errToken := token.ExtractTokenID(c)

	if errToken != nil {
		fmt.Println(errToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login to continue"})
		return
	}
	roleToken, errToken := token.ExtractTokenRole(c)
	if errToken != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error})
		return

	}
	if roleToken == "admin" {
		if err := db.Where("id = ? ", c.Param("id"), UserIdCurrent).First(&Comments).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	} else {
		if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), UserIdCurrent).First(&Comments).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	}

	//Validate Input
	var input commentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedInput := models.Comments{PostComment: input.CommentPost,
		Publish:   input.Publish,
		UserID:    UserIdCurrent,
		UpdatedAt: time.Now(),
	}

	var err error = db.Model(&Comments).Updates(updatedInput).Error

	db.Model(&Comments).Association("Meta").Replace(Comments.Meta)

	for key, value := range input.CommentMeta {
		metaInput := models.CommentMeta{
			CommentID: Comments.ID,
			MetaKey:   key,
			MetaValue: value,
		}
		var err error = db.Create(&metaInput).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": Comments})

}

//DeleteComment godoc
// @Summary Delete Comment
// @Description Delete a Comment by id
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Tags Comments
// @Produce json
// @Param id path string true "Comments id"
// @Success 200 {object} map[string]boolean
// @Router /comment/{id} [delete]
func DeleteComment(c *gin.Context) {
	//Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var Comments models.Comments
	UserIdCurrent, errToken := token.ExtractTokenID(c)

	if errToken != nil {
		fmt.Println(errToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login to continue"})
		return
	}
	roleToken, errToken := token.ExtractTokenRole(c)
	if errToken != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errToken.Error})
		return

	}
	if roleToken == "admin" {
		if err := db.Where("id = ? ", c.Param("id"), UserIdCurrent).First(&Comments).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	} else {
		if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), UserIdCurrent).First(&Comments).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	}
	db.Delete(&Comments)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
