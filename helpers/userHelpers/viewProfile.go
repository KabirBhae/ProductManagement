package userHelpers

import (
	"ProductManagement/helpers"
	"ProductManagement/models"
	"ProductManagement/responses"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func ViewProfile(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	usernamefromURL := c.Param("username")

	if usernamefromURL != usernameFromClaims {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "unauthorised user", Data: &echo.Map{"data": ""}})
	}

	var user models.User
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"username": usernamefromURL}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusNoContent, responses.UserResponse{Status: http.StatusNoContent, Message: "username doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": user}})
}
