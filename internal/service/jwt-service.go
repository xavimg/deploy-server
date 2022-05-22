package service

import (
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

type JwtCustomClaim struct {
	UserID uint64 `json:"user_id"`
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
	claims := &JwtCustomClaim{
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
	claims := &JwtCustomClaim{
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

	return t
}

func (j *jwtService) ValidateToken(tokenSigned string) (*jwt.Token, error) {
	token, _ := jwt.ParseWithClaims(
		tokenSigned,
		&JwtCustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("turingoffworld"), nil
		},
	)
	token.Valid = true

	return token, nil
}
