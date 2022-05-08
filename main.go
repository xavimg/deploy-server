package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xavimg/Turing/apituringserver/internal/config"
	"github.com/xavimg/Turing/apituringserver/internal/controller"
	"github.com/xavimg/Turing/apituringserver/internal/middleware"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
	"github.com/xavimg/Turing/apituringserver/internal/service"
	"gorm.io/gorm"

	"github.com/xavimg/Turing/apituringserver/docs"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var (
	db              *gorm.DB                   = config.SetupDatabaseConnection()
	userRepository  repository.UserRepository  = repository.NewUserRepository(db)
	adminRepository repository.AdminRepository = repository.NewAdminRepository(db)
	authRepository  repository.AuthRepository  = repository.NewAuthRepository(db)

	jwtService   service.JWTService   = service.NewJWTService()
	userService  service.UserService  = service.NewUserService(userRepository)
	authService  service.AuthService  = service.NewAuthService(userRepository, authRepository)
	adminService service.AdminService = service.NewAdminService(adminRepository)

	authController  controller.AuthController  = controller.NewAuthController(authService, jwtService)
	userController  controller.UserController  = controller.NewUserController(userService, jwtService)
	adminController controller.AdminController = controller.NewAdminController(adminService, authService, jwtService)
)

func main() {

	docs.SwaggerInfo.Title = "Server Turing API"
	docs.SwaggerInfo.Description = "API for testing every endpoint from Turing API server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	// public routes
	authRoutes := r.Group("api/auth")
	{
		authRoutes.GET("/google/login", authController.GoogleLogin)
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout/:id", authController.Logout)
		authRoutes.POST("/verifyaccount", authController.VerifyAccount)
	}

	// private/tokenized routes
	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/update", userController.Update)
		userRoutes.DELETE("/profile/:id", userController.DeleteAccount)
	}

	adminRoutes := r.Group("api/admin")
	{
		adminRoutes.POST("/register", adminController.AdminRegister)
		adminRoutes.POST("/login", adminController.AdminLogin)
		adminRoutes.GET("/users/:typeUser", middleware.CheckRole(userService), adminController.ListAllUsersByParameter)
		adminRoutes.PUT("/ban/:id", middleware.CheckRole(userService), adminController.BanUser)
		adminRoutes.PUT("/unban/:id", middleware.CheckRole(userService), adminController.UnbanUser)
		adminRoutes.POST("/newfeature", middleware.CheckRole(userService), adminController.NewFeature)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
