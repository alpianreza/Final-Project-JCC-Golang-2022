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

type postsInput struct {
	PostTitle   string            `json:"post_title" binding:"required" gorm:"unique"`
	PostContent string            `json:"post_content" binding:"required"`
	Publish     bool              `json:"publish" binding:"required"`
	CategoryID  int               `json:"category_id" binding:"required"`
	PostMeta    map[string]string `json:"post_meta" binding:"required"`
}

// GetAllPost godoc
// @Summary Get All Posts
// @Description Get a list of Posts
// @Tags Posts
// @Produce json
// @Success 200 {object} []models.Posts
// @Router /posts [get]
func GetAllPost(c *gin.Context) {
	//get db from fin context
	db := c.MustGet("db").(*gorm.DB)
	var posts []models.Posts
	q := c.Request.URL.Query().Get("q")

	db.Preload("User").Preload("Category").Preload("Favorit").Where("title LIKE ? OR content LIKE ? OR meta_title LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%").Find(&posts)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// CreatePost godoc
// @Summary Create New Posts
// @Description Creating a new Posts
// @Tags Posts
// @Param Body body postsInput true "the body to create a new Posts"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Param Body body postsInput true "the body create a new Posts"
// @Produce json
// @Success 200 {object} []models.Posts
// @Router /posts [post]
func CreatePost(c *gin.Context) {
	var input postsInput
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

	// Create Posts
	posts := models.Posts{PostTitle: input.PostTitle,
		PostContent: input.PostContent,
		Publish:     input.Publish,
		UserID:      UserIdCurrent,
		CategoryID:  input.CategoryID,
		CreatedAt:   time.Now(),
	}

	db := c.MustGet("db").(*gorm.DB)
	_, err := posts.PSave(input.PostMeta, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// GetPostById godoc
// @Summary Get Posts
// @Description Get Post by Id
// @Tags Post
// @Produce json
// @Param id path string trur "post id"
// @Success 200 {object} []models.Posts
// @Router /posts/{id} [get]
func GetPostById(c *gin.Context) {
	var posts models.Posts
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("id = ?", c.Param("id")).First(&posts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
		return
	}

	db.Preload("User").Preload("Category").Preload("Comments").Preload("Meta").Find(&posts)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// UpdatePost godoc
// @Summary Update Post.
// @Description Update Post by id.
// @Tags Post
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Param id path string true "Post id"
// @Param Body body postsInput true "the body to update an Posts"
// @Success  200 {object} models.Posts
// @Router /posts/{id} [patch]
func UpdatePost(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//Get model if exist
	var Posts models.Posts
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
		if err := db.Where("id = ? ", c.Param("id"), UserIdCurrent).First(&Posts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	} else {
		if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), UserIdCurrent).First(&Posts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	}

	//Validate Input
	var input postsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedInput := models.Posts{PostTitle: input.PostTitle,
		PostContent: input.PostContent,
		Publish:     input.Publish,
		UserID:      UserIdCurrent,
		CategoryID:  input.CategoryID,
		UpdatedAt:   time.Now(),
	}

	var err error = db.Model(&Posts).Updates(updatedInput).Error

	db.Model(&Posts).Association("Meta").Replace(Posts.Meta)

	for key, value := range input.PostMeta {
		metaInput := models.PostMeta{
			PostID:    Posts.ID,
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

	c.JSON(http.StatusOK, gin.H{"data": Posts})

}

//DeletePost godoc
// @Summary Delete Posts.
// @Description Delete a Posts by id.
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Tags Posts
// @Produce json
// @Param id path string true "Posts id"
// @Success 200 {object} map[string]boolean
// @Router /posts/{id} [delete]
func DeletePost(c *gin.Context) {
	//Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var Posts models.Posts
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
		if err := db.Where("id = ? ", c.Param("id"), UserIdCurrent).First(&Posts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	} else {
		if err := db.Where("id = ? AND user_id = ? ", c.Param("id"), UserIdCurrent).First(&Posts).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not found!"})
			return
		}
	}
	db.Model(&Posts).Association("Comments").Replace(Posts.Comments)
	db.Delete(&Posts)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
