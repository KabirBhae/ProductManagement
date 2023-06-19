package userHelpers

import (
	"ProductManagement/configs"
	"ProductManagement/helpers"
	"ProductManagement/models"
	"ProductManagement/responses"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type loginUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var loginUser loginUser
	if err := c.Bind(&loginUser); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error while binding", Data: &echo.Map{"data": err.Error()}})
	}
	//use the validator library to validate required fields
	if validationErr := validate.Struct(&loginUser); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error while validating request from user", Data: &echo.Map{"data": validationErr.Error()}})
	}

	var user models.User

	err := userCollection.FindOne(ctx, bson.M{"username": loginUser.Username, "status": "Active"}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "No such user exists", Data: &echo.Map{"data": err.Error()}})
	}

	if user.Password != loginUser.Password {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "User Credentials don't match", Data: &echo.Map{"data": ""}})
	}

	// Set custom claims
	claims := &helpers.JwtCustomClaims{
		user.Username,
		user.Password,
		user.IsAdmin,
		user.ID.Hex(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	secretKey := configs.SecretKey()
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"user": user}, Token: t})
}
