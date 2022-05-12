package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/xavimg/Turing/apituringserver/internal/config"
	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/helper"
	"github.com/xavimg/Turing/apituringserver/internal/service"
)

const (
	urlAndreba = "https://api-turing-api.herokuapp.com"
)

// AuthController interface is a contract what this controller can do
type AuthController interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
	Logout(context *gin.Context)
	VerifyAccount(context *gin.Context)
	GoogleLogin(context *gin.Context)
}
type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}
type JsonAndreba struct {
	Isvalid bool `json:"isvalid"`
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

// Login godoc
// @Title Login
// @Description  Login to the server. Check token with backend.
// @Param request body dto.LoginDTO true "Body to login"
// @Tags Auth
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/auth/login [post]
func (c *authController) Login(context *gin.Context) {
	var loginDTO dto.LoginDTO

	// validation from request
	if errDTO := context.ShouldBindJSON(&loginDTO); errDTO != nil {
		response := helper.BuildErrorResponse("User login failed", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Verify of credentials exists
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		if strconv.FormatBool(v.Active) == "false" {
			context.JSON(http.StatusBadRequest, "User has been banned")
			return

		}

		generateToken := c.jwtService.GenerateTokenLogin(v.ID)
		v.Token = fmt.Sprintf("Bearer %v", generateToken)
		c.authService.SaveToken(v, v.Token)

		errEnv := godotenv.Load(".env")
		if errEnv != nil {
			log.Println("impossible get .env")
		}
		// Check .env file to change debug mode.
		if os.Getenv("DEBUG_MODE") == "off" {
			client := &http.Client{}
			url := fmt.Sprintf("%v/player/signin", urlAndreba)
			req, err := http.NewRequest("POST", url, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", generateToken))

			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			if resp.StatusCode != 200 {
				log.Println("Something went wrong")
				return
			}
		}

		response := helper.BuildResponseSession(true, "User login successfully", generateToken)
		context.JSON(http.StatusOK, response)
		return
	}

	response := helper.BuildErrorResponse("User login failed", "Invalid credential", helper.EmptyObj{})
	context.AbortWithStatusJSON(http.StatusUnauthorized, response)

}

// Register godoc
// @Title Register
// @Description  Register to the server as a new user. Sends token to backend.
// @Param request body dto.RegisterDTO true "Body to register"
// @Tags Auth
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/auth/register [post]
func (c *authController) Register(context *gin.Context) {
	var registerDTO dto.RegisterDTO

	if errDTO := context.ShouldBind(&registerDTO); errDTO != nil {
		response := helper.BuildErrorResponse("User register failed", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	duplicated, err := c.authService.IsDuplicateEmail(registerDTO.Email)
	if err != nil {
		return
	}
	if !duplicated {
		response := helper.BuildErrorResponse("User register failed", "Duplicate email", helper.EmptyObj{})
		context.JSON(http.StatusConflict, response)
		return
	}

	// getCode := service.SendEmailCodeVerify(registerDTO.Name, registerDTO.Email)
	// registerDTO.CodeVerify = getCode

	createdUser := c.authService.CreateUser(registerDTO)

	token := c.jwtService.GenerateTokenRegister(createdUser.ID)
	createdUser.Token = fmt.Sprintf("Bearer %v", token)

	json_data, err := json.Marshal(createdUser.ID)
	if err != nil {
		return
	}

	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Println("impossible get .env")
	}
	// Check .env file to change debug mode.
	if os.Getenv("DEBUG_MODE") == "off" {
		url := fmt.Sprintf("%v/player/signup", urlAndreba)

		resp, err := http.Post(url, "application/json", bytes.NewReader(json_data))
		if err != nil {
			log.Println(err)
		}
		fmt.Println(resp.StatusCode)
		defer resp.Body.Close()
	}

	response := helper.BuildResponse(true, "Check your email !", createdUser)
	context.JSON(http.StatusCreated, response)
}

// Logout godoc
// @Title Logout
// @Description  Logout to the server
// @Param id path string true "ID from query"
// @Tags Auth
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/auth/logout [post]
func (c *authController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		response := helper.BuildErrorResponse("Failed to process request", "No token found", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	token, _ := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		return []byte("turingoffworld"), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	id := claims["user_id"]

	fmt.Println(id)

	authResult := c.authService.VerifyUserExist(id)
	fmt.Println(authResult)
	if authResult != false {
		ctx.JSON(http.StatusBadRequest, "user doesn't exist")
		return
	}

	if v, ok := authResult.(entity.User); ok {
		token, err := c.authService.GetToken(id)
		if err != nil {
			return
		}

		url := fmt.Sprintf("%v/player/signout", urlAndreba)
		client := &http.Client{}
		req, err := http.NewRequest("POST", url, nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token))

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusBadRequest, "No authorization key found")
			return
		}

		c.authService.DeleteToken(v, "")
	}

	ctx.JSON(http.StatusOK, "session logout")
}

// verifyAccount godoc
// @Title verifyAccount
// @Description  Verify the account with code send to email.
// @Param request body dto.CodeVerifyDTO true "Body to verify account"
// @Tags Auth
// @Success      200 {object} helper.Response
// @Failure      400 body is empty or missing param
// @Failure      500 "internal server error"
// @Router       /api/auth/verifyaccount [post]
func (c *authController) VerifyAccount(ctx *gin.Context) {
	var req dto.CodeVerifyDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Fatal("Error binding")
		return
	}

	if req.Email == "" {
		log.Println("email and code are required")
		return
	}
	if req.Code <= 0 {
		log.Println("email and code are required")
		return
	}

	exist, err := c.authService.VerifyCode(req.Email, req.Code)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	if !exist {
		ctx.JSON(http.StatusBadRequest, "invalid code !")
		return
	}

	ctx.JSON(http.StatusOK, "you've been verified !")
}

func (c *authController) GoogleLogin(ctx *gin.Context) {
	googleConfig := config.SetupConfigGoogle()
	url := googleConfig.AuthCodeURL("randomstate")

	ctx.Redirect(303, url)
}
