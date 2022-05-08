package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/helper"
	"github.com/xavimg/Turing/apituringserver/internal/service"
)

type AdminController interface {
	AdminRegister(ctx *gin.Context)
	AdminLogin(ctx *gin.Context)
	ListAllUsersByParameter(ctx *gin.Context)
	BanUser(ctx *gin.Context)
	UnbanUser(ctx *gin.Context)
	NewFeature(ctx *gin.Context)
}

type adminController struct {
	adminService service.AdminService
	authService  service.AuthService
	jwtService   service.JWTService
}

func NewAdminController(adminService service.AdminService, authService service.AuthService, jwtService service.JWTService) AdminController {
	return &adminController{
		adminService: adminService,
		authService:  authService,
		jwtService:   jwtService,
	}
}

// Register for admin
// @Title AdminRegister
// @Description  Register to the server as a new admin.
// @Param request body dto.RegisterDTO true "Body to register"
// @Tags Auth
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/admin/register [post]
func (c *adminController) AdminRegister(context *gin.Context) {
	var registerDTO dto.RegisterDTO

	if errDTO := context.ShouldBind(&registerDTO); errDTO != nil {
		response := helper.BuildErrorResponse("User register failed", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("User register failed", "Duplicate email", helper.EmptyObj{})
		context.JSON(http.StatusConflict, response)
		return
	}

	getCode := service.SendEmailCodeVerify(registerDTO.Name, registerDTO.Email)
	registerDTO.CodeVerify = getCode

	createdUser := c.adminService.CreateAdmin(registerDTO)

	token := c.jwtService.GenerateTokenRegister(createdUser.ID)
	createdUser.Token = fmt.Sprintf("Bearer %v", token)

	var routine sync.Mutex
	routine.Lock()
	go service.SendEmail(registerDTO.Name, registerDTO.Email)
	routine.Unlock()

	response := helper.BuildResponse(true, "Check your email !", createdUser)
	context.JSON(http.StatusCreated, response)
}

// Login as admin
// @Title AdminLogin
// @Description entering system as admin.
// @Param request body dto.LoginDTO true "Body to register"
// @Tags Admin
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/admin/login [post]
func (c *adminController) AdminLogin(context *gin.Context) {
	var loginDTO dto.LoginDTO
	if errDTO := context.ShouldBindJSON(&loginDTO); errDTO != nil {
		response := helper.BuildErrorResponse("admin login failed", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response.Message)
		return
	}

	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		if v.TypeUser != "admin" {
			context.JSON(http.StatusBadRequest, "admin doesn't exists")
			return
		}

		generateToken := c.jwtService.GenerateTokenLogin(v.ID)
		v.Token = fmt.Sprintf("Bearer %v", generateToken)
		c.authService.SaveToken(v, fmt.Sprintf("Bearer %v", generateToken))

		response := helper.BuildResponseSession(true, "admin login successfully", generateToken)
		context.JSON(http.StatusOK, response)
		return
	}

	response := helper.BuildErrorResponse("admin login failed", "Invalid credential", helper.EmptyObj{})
	context.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

// Listing Users by param godoc
// @Title Listing users
// @Description  List users depending on param.
// @Param typeUser path string true "typeUser from query"
// @Tags Admin
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/admin/users/:typeUser [get]
func (c *adminController) ListAllUsersByParameter(ctx *gin.Context) {
	tUser := ctx.Param("typeUser")
	var users []entity.User

	switch tUser {
	case "all":
		users = c.adminService.ListAllUsers()
	case "ban":
		users = c.adminService.ListAllUsersByActive()
	case "admin":
		users = c.adminService.ListAllUsersByTypeAdmin()
	case "user":
		users = c.adminService.ListAllUsersByTypeUser()
	default:
		ctx.JSON(http.StatusBadRequest, nil)
	}

	ctx.JSON(http.StatusOK, users)
}

// BanUser godoc
// @Title BanUser
// @Description  Admin ban user for X time.
// @Param Authorization header string true "Token acces admin"
// @Param id path string true "ID from query"
// @Tags Admin
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/admin/ban/{id} [put]
func (c *adminController) BanUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	c.adminService.BanUser(userID)

	res := helper.BuildResponse(true, "User has been banned !", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)

}

// UnbanUser godoc
// @Title UnbanUser
// @Description  Admin unban user.
// @Param Authorization header string true "Token acces admin"
// @Param id path string true "ID from query"
// @Tags Admin
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/admin/unban/{id} [put]
func (c *adminController) UnbanUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	c.adminService.UnbanUser(userID)

	res := helper.BuildResponse(true, "User has been unbanned !", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)

}

// NewFeature godoc
// @Title NewFeature
// @Description  Admin add new feature to show in version of game info.
// @Param Authorization header string true "Token acces admin"
// @Param request body dto.FeatureDTO true "Body to write new features"
// @Tags Admin
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/admin/newfeature [post]
func (c *adminController) NewFeature(ctx *gin.Context) {
	var feature dto.FeatureDTO

	if err := ctx.ShouldBind(&feature); err != nil {
		res := helper.BuildErrorResponse(
			"Feature not created", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	featureCreated := c.adminService.NewFeature(feature)

	response := helper.BuildResponse(true, "Feature has been created", featureCreated)

	ctx.JSON(http.StatusCreated, response.Data)
}
