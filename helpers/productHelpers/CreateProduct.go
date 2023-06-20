package productHelpers

import (
	"ProductManagement/configs"
	"ProductManagement/helpers"
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
var productCollection = configs.GetCollection(configs.DB, "products")
var validate = validator.New()

func CreateProduct(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var seller models.User
	err := userCollection.FindOne(ctx, bson.M{"username": usernameFromClaims, "status": "Active"}).Decode(&seller)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "seller doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}

	var product models.Product
	//validate the request body
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error while binding", Data: &echo.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&product); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "please provide ProductID, SellerID, Name, Price and Quantity correctly", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newProduct := models.Product{
		ProductID:      primitive.NewObjectID(),
		SellerID:       seller.ID,
		Name:           product.Name,
		Price:          product.Price,
		Quantity:       product.Quantity,
		SellerUserName: seller.Username,
		IsAvailable:    true,
	}

	result, err := productCollection.InsertOne(ctx, newProduct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error inserting product in db", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "successfully created new product", Data: &echo.Map{"created_product": result}})
}
