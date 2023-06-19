package userHelpers

import (
	"ProductManagement/helpers"
	"ProductManagement/models"
	"ProductManagement/responses"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func EditAProfile(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	usernamefromURL := c.Param("username")

	if usernamefromURL != usernameFromClaims {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "unauthorised user", Data: &echo.Map{"data": ""}})
	}

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": usernameFromClaims}).Decode(&existingUser)
	if err != nil {
		return c.JSON(http.StatusNoContent, responses.UserResponse{Status: http.StatusNoContent, Message: "username doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}
	fmt.Println("only name: " + user.Name + " and password: " + user.Password + " were updated")

	updatedUser := models.User{
		ID:       existingUser.ID,
		Name:     user.Name,
		Username: existingUser.Username,
		Email:    existingUser.Email,
		Password: user.Password,
		Balance:  existingUser.Balance,
		Status:   existingUser.Status,
		IsAdmin:  existingUser.IsAdmin,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"username": usernameFromClaims}, bson.M{"$set": updatedUser})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "update successful", Data: &echo.Map{"data": result}})
}
