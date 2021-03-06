package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	db *gorm.DB = config.SetupDatabaseConnection()

	userRepository  repository.UserRepository  = repository.NewUserRepository(db)
	adminRepository repository.AdminRepository = repository.NewAdminRepository(db)
	authRepository  repository.AuthRepository  = repository.NewAuthRepository(db)
	shopRepository  repository.ShopRepository  = repository.NewShopRepository(db)

	jwtService   service.JWTService   = service.NewJWTService()
	userService  service.UserService  = service.NewUserService(userRepository)
	authService  service.AuthService  = service.NewAuthService(userRepository, authRepository)
	adminService service.AdminService = service.NewAdminService(adminRepository)
	shopService  service.ShopService  = service.NewShopService(shopRepository)

	authController  controller.AuthController  = controller.NewAuthController(authService, jwtService)
	userController  controller.UserController  = controller.NewUserController(userService, jwtService)
	adminController controller.AdminController = controller.NewAdminController(adminService, authService, jwtService)
	shopController  controller.ShopController  = controller.NewShopController(shopService, jwtService, adminService)
)

func main() {

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Println("impossible get .env")
	}

	serverPort := os.Getenv("PORT")

	docs.SwaggerInfo.Title = "Server Turing API"
	docs.SwaggerInfo.Description = "API for testing every endpoint from Turing API server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers"},
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Salle",
		})
	})
	r.GET("/hello", authController.GoogleCallback)
	// public routes
	authRoutes := r.Group("api/auth")
	{
		authRoutes.GET("/google/login", authController.GoogleLogin)
		authRoutes.GET("/google/callback", authController.GoogleCallback)
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
		authRoutes.POST("/verifyaccount", authController.VerifyAccount)
	}

	// private/tokenized routes
	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
		userRoutes.DELETE("/profile", userController.DeleteAccount)

		userRoutes.POST("/add-friend/:id", userController.AddFriend)
		userRoutes.DELETE("/remove-friend/:id", userController.RemoveFriend)
		userRoutes.GET("/friendlist", userController.ShowFriendlist)
		userRoutes.GET("/check-friend/:id", userController.IsFriend)
		userRoutes.POST("/send-message/:id", userController.SendMessage)
		userRoutes.GET("/list-messages", userController.ListMessages)
		userRoutes.GET("/list-messages/:id", userController.MessageDetail)

		userRoutes.GET("/products", shopController.ListProducts)
		userRoutes.GET("/product-detail/:id", shopController.ProductDetail)

		userRoutes.POST("/product/cart/:idcart/product/:id", shopController.AddProductCart)
		userRoutes.POST("/product/cart/add-payment/:idcart", shopController.AddCreditCard)
		userRoutes.PUT("/product/cart/add-payment/:idcart", shopController.AddCreditCard)
		userRoutes.DELETE("/product/cart/:idcart", shopController.DeleteCart)
		userRoutes.POST("/confirm/cart/:idcart", shopController.ConfirmPayment)
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

	shopAdminRoutes := r.Group("api/admin")
	{
		shopAdminRoutes.POST("/product", middleware.CheckRole(userService), shopController.InsertProduct)
		shopAdminRoutes.DELETE("/product/:id", middleware.CheckRole(userService), shopController.DeleteProduct)
		shopAdminRoutes.PUT("/product/:id", middleware.CheckRole(userService), shopController.UpdateProduct)
		shopAdminRoutes.GET("/product", middleware.CheckRole(userService), shopController.InsertProduct)
		shopAdminRoutes.GET("/buys", middleware.CheckRole(userService), shopController.ListBuys)

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf(":%v", serverPort))
}
