package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// JWTService is a contract of what jwtService can do
type JWTService interface {
	GenerateTokenLogin(userID uint64) string
	GenerateTokenRegister(userID uint64) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID uint64 `json:"user_id"`
	// Signature string `json:"signature"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

// NewJWTService method is creates a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    "turingoffworld",
		secretKey: goDotEnvVariable("ACCESS_SECRET"),
	}
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func (j *jwtService) GenerateTokenLogin(UserID uint64) string {

	claims := &jwtCustomClaim{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			// ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:   j.issuer,
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err.Error())
	}

	return t
}

func (j *jwtService) GenerateTokenRegister(UserID uint64) string {
	claims := &jwtCustomClaim{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("turingoffworld"))
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(t)
	return t
}

func (j *jwtService) ValidateToken(tokenSigned string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenSigned,
		&jwtCustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("turingoffworld"), nil
		},
	)
	if err != nil {
		return nil, nil
	}

	claims, ok := token.Claims.(jwtCustomClaim)
	if !ok {
		log.Println("couldnt't parse claims")
		return nil, nil
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Println("token expired")
		return nil, nil
	}
	return nil, nil

	// return jwt.Parse(tokenSigned, func(t_ *jwt.Token) (interface{}, error) {
	// 	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
	// 	}

	// 	return []byte(j.secretKey), nil
	// })
}
