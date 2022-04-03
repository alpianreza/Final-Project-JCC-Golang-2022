package routes

import (
	"finalproject/controllers"
	"finalproject/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/gorm"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	// Auth
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// User
	userMiddlewareRoute := r.Group("/users")
	userMiddlewareRoute.Use(middlewares.JwtAuthMiddleware("admin", "guest"))
	userMiddlewareRoute.PATCH("/:id", controllers.UpdateUser)
	usersMiddlewareRoute := r.Group("/users")
	usersMiddlewareRoute.Use(middlewares.JwtAuthMiddleware("admin"))
	usersMiddlewareRoute.GET("/", controllers.GetAllUser)
	usersMiddlewareRoute.GET("/:id", controllers.GetUserById)
	usersMiddlewareRoute.POST("/", controllers.CreateUser)
	usersMiddlewareRoute.DELETE("/:id", controllers.DeleteUser)

	// Post
	r.GET("/posts", controllers.GetAllPost)
	r.GET("/posts/:id", controllers.GetPostById)
	postsMiddlewareRoute := r.Group("/posts")
	postsMiddlewareRoute.Use(middlewares.JwtAuthMiddleware("user", "admin"))
	postsMiddlewareRoute.POST("/", controllers.CreatePost)
	postsMiddlewareRoute.PATCH("/:id", controllers.UpdatePost)
	postsMiddlewareRoute.DELETE("/:id", controllers.DeletePost)

	// Comment
	r.GET("/comments", controllers.GetAllComment)
	r.POST("/comments", controllers.CreateComment)
	r.GET("/comments/:id", controllers.GetCommentById)
	postCommentsMiddlewareRoute := r.Group("/post_comments")
	postCommentsMiddlewareRoute.Use(middlewares.JwtAuthMiddleware("admin"))
	postCommentsMiddlewareRoute.DELETE("/:id", controllers.DeleteComment)

	// Category
	r.GET("/category", controllers.GetAllCategory)
	r.GET("/category/:id", controllers.GetCategoryById)
	categoriesMiddlewareRoute := r.Group("/category")
	categoriesMiddlewareRoute.Use(middlewares.JwtAuthMiddleware("admin"))
	categoriesMiddlewareRoute.POST("/", controllers.CreateCategory)
	categoriesMiddlewareRoute.PATCH("/:id", controllers.UpdateCategory)
	categoriesMiddlewareRoute.DELETE("/:id", controllers.DeleteCategory)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
