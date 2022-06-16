package controller

import (
	"log"
	"net/http"
	"strconv"

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

	AddFriend(ctx *gin.Context)
	ShowFriendlist(ctx *gin.Context)
	RemoveFriend(ctx *gin.Context)
	IsFriend(ctx *gin.Context)
	SendMessage(ctx *gin.Context)
	ListMessages(ctx *gin.Context)
	MessageDetail(ctx *gin.Context)
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

// AddFriend godoc
// @Title AddFriend
// @Description AddFriend.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/friend [put]
func (c *userController) AddFriend(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	idFriend := ctx.Param("id")

	user, err := c.userService.Profile(idFriend)
	if err != nil {
		res := helper.BuildErrorResponse(
			"add user failed", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if user.Email == "" {
		ctx.JSON(http.StatusBadRequest, "user you want to add it doesn't exists")
		return
	}

	if err := c.userService.AddFriend(id, user); err != nil {
		return
	}

	ctx.JSON(200, "New friend added")
}

// ShowFriendList godoc
// @Title ShowFriendlist
// @Description Show all friends from one user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/friendlist [get]
func (c *userController) ShowFriendlist(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	friendList, err := c.userService.ShowFriendlist(id)
	if err != nil {
		return
	}

	ctx.JSON(200, friendList)
}

// Deletefriend godoc
// @Title ShowFriendlist
// @Description Show all friends from one user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/friendlist [get]
func (c *userController) RemoveFriend(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	idFriend := ctx.Param("id")
	idInt, _ := strconv.Atoi(idFriend)

	if err := c.userService.RemoveFriend(id, uint64(idInt)); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, "friend removed")
}

// IsFriend godoc
// @Title ShowFriendlist
// @Description Show all friends from one user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/remove-friend/{id} [get]
func (c *userController) IsFriend(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	idFriend := claims["user_id"]

	_, err := c.userService.IsFriend(idFriend)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, "friend exists in the friendlist")
}

// SendMessage godoc
// @Title SendMessage
// @Description Show all friends from one user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/send-message/{id} [post]
func (c *userController) SendMessage(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	var message dto.MessageDTO
	if err := ctx.ShouldBindJSON(&message); err != nil {
		res := helper.BuildErrorResponse(
			"body empty", err.Error(),
			helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	message.From = id.(float64)

	idFriend := ctx.Param("id")
	idInt, _ := strconv.Atoi(idFriend)
	message.To = float64(idInt)

	_, err := c.userService.IsFriend(idInt)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := c.userService.SendMessage(message); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, "message sended")
}

// ListMessages godoc
// @Title ListMessages
// @Description Show all friends from one user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/remove-friend/{id} [get]
func (c *userController) ListMessages(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	notifications, err := c.userService.ListMessages(id)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, notifications)
}

// MessageDetail godoc
// @Title MessageDetail
// @Description Show all friends from one user.
// @Param Authorization header string true "Token acces login"
// @Tags User
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/user/list-messages/{id} [get]
func (c *userController) MessageDetail(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	idMessage := ctx.Param("id")
	idM, err := strconv.Atoi(idMessage)
	if err != nil {
		return
	}

	notifications, err := c.userService.MessageDetail(idM, id)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, notifications)
}
