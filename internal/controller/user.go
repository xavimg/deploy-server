package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	user, err := c.userService.Profile(id)
	if err != nil {
		return
	}

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

	if err := ctx.ShouldBindJSON(&userUpdateDTO); err != nil {
		res := helper.BuildErrorResponse(
			"Update user failed", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	if _, err := c.userService.Update(userUpdateDTO, id, userUpdateDTO); err != nil {
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, "succesfully updated")
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
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"].(float64)

	if err := c.userService.DeleteAccount(id); err != nil {
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, "your account is gonna be delated once you log out")
}
