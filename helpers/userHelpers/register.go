package userHelpers

import (
	"ProductManagement/configs"
	"ProductManagement/models"
	"ProductManagement/responses"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

var userCollection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Register(c echo.Context, isAdminParam bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error while binding", Data: &echo.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "please provide name, username, email and password correctly", Data: &echo.Map{"data": validationErr.Error()}})
	}
	var existingUser models.User

	err2 := userCollection.FindOne(ctx, bson.M{"username": user.Username, "status": "Active"}).Decode(&existingUser)
	//username already in use and user status is active
	if err2 == nil {
		return c.JSON(http.StatusConflict, responses.UserResponse{Status: http.StatusConflict, Message: "user with this username already exists", Data: &echo.Map{"data": ""}})
	}
	err3 := userCollection.FindOne(ctx, bson.M{"email": user.Email, "status": "Active"}).Decode(&existingUser)

	//email already in use and user status is active
	if err3 == nil {
		return c.JSON(http.StatusConflict, responses.UserResponse{Status: http.StatusConflict, Message: "user with this email already exists", Data: &echo.Map{"data": ""}})
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Balance:  1000.0,
		Status:   "Active",
		IsAdmin:  isAdminParam,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error inserting user in db", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"created_user": result}})
}
