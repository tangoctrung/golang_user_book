package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangoctrung/golang_api_v2/config"
	"github.com/tangoctrung/golang_api_v2/controller"
	"github.com/tangoctrung/golang_api_v2/docs"
	"github.com/tangoctrung/golang_api_v2/middleware"
	"github.com/tangoctrung/golang_api_v2/repository"
	"github.com/tangoctrung/golang_api_v2/service"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var (
	db             *gorm.DB                  = config.ConnectDatabase()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	bookRepository repository.BookRepository = repository.NewBookRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	authService    service.AuthService       = service.NewAuthService(userRepository)
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
	userService    service.UserService       = service.NewUserService(userRepository)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
	bookService    service.BookService       = service.NewBookService(bookRepository)
	bookController controller.BookController = controller.NewBookController(bookService, jwtService)
)

func main() {

	docs.SwaggerInfo.Title = "Pragmatic Reviews - User book API"
	docs.SwaggerInfo.Description = "Pragmatic Reviews - User, book API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Server is running",
		})
	})

	comonRoutes := r.Group("api/v1")
	{
		authRoutes := comonRoutes.Group("/auth")
		{
			authRoutes.POST("/login", authController.Login)
			authRoutes.POST("/register", authController.Register)
		}

		userRoutes := comonRoutes.Group("/user", middleware.AuthorizeJWT(jwtService))
		{
			userRoutes.PUT("/update", userController.Update)
			userRoutes.GET("/profile", userController.Profile)
		}

		bookRoutes := comonRoutes.Group("/book", middleware.AuthorizeJWT(jwtService))
		{
			bookRoutes.POST("/create-book", bookController.InsertBook)
			bookRoutes.PUT("/update-book", bookController.UpdateBook)
			bookRoutes.DELETE("/delete-book", bookController.DeleteBook)
			bookRoutes.GET("/:id", bookController.FindBookByID)
			bookRoutes.GET("/all", bookController.GetAllBooks)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8000")
}
