package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/helper"
	"github.com/xavimg/Turing/apituringserver/internal/service"
)

type UserController interface {
	Profile(context *gin.Context)
	Update(context *gin.Context)
	DeleteAccount(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Profile godoc
// @Title Profile
// @Description Profile of X user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/profile [get]
func (c *userController) Profile(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	// Send ID to service
	user := c.userService.Profile(userID)

	// response
	res := helper.BuildResponse(true, "Get user profile successfully", user)
	ctx.JSON(http.StatusOK, res)
}

// Update godoc
// @Title Update
// @Description Update profile.
// @Param Authorization header string true "Token acces login"
// @Param request body dto.UserUpdateDTO true "Update profile user"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/update [put]
func (c *userController) Update(ctx *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO

	// NewPass
	err := ctx.ShouldBindJSON(&userUpdateDTO)
	if err != nil {
		res := helper.BuildErrorResponse(
			"Update user failed", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// Get token from userController
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])

	user := c.userService.Update(userUpdateDTO, userID, userUpdateDTO)

	res := helper.BuildResponse(true, "Update user successfully", user)
	ctx.JSON(http.StatusOK, res)

}

// DeleteAccount godoc
// @Title DeleteAccount
// @Description Delete account profile.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/deleteaccount [delete]
func (c *userController) DeleteAccount(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	number, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return
	}

	// Send ID to service
	user := c.userService.DeleteAccount(number)

	// response
	res := helper.BuildResponse(true, "user deleted profile successfully", user)
	ctx.JSON(http.StatusOK, res)
}
