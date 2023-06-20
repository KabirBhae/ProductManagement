package productHelpers

import (
	"ProductManagement/helpers"
	"ProductManagement/models"
	"ProductManagement/responses"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type productDetails struct {
	Id       string `json:"productID" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

func BuyProduct(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var buyer models.User
	err := userCollection.FindOne(ctx, bson.M{"username": usernameFromClaims, "status": "Active"}).Decode(&buyer)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "buyer doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}

	var requestedProduct productDetails
	//validate the request body
	if err := c.Bind(&requestedProduct); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error while binding", Data: &echo.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&requestedProduct); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "please provide ProductID, SellerID, Name, Price and Quantity correctly", Data: &echo.Map{"data": validationErr.Error()}})
	}
	//non positive not allowed
	if requestedProduct.Quantity < 1 {
		return c.JSON(http.StatusNotAcceptable, responses.UserResponse{Status: http.StatusNotAcceptable, Message: "please provide quantity greater than 0", Data: &echo.Map{"data": ""}})
	}

	requestedProductObjectId, errr := primitive.ObjectIDFromHex(requestedProduct.Id)
	if errr != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "string cannot be parsed to ObjectID", Data: &echo.Map{"data": err.Error()}})
	}

	var existingProduct models.Product
	err2 := productCollection.FindOne(ctx, bson.M{"_id": requestedProductObjectId, "isavailable": true}).Decode(&existingProduct)
	if err2 != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "product doesn't exists in DB", Data: &echo.Map{"data": err2.Error()}})
	}
	if existingProduct.Quantity < 1 {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "product doesn't exists in DB", Data: &echo.Map{"data": ""}})
	}
	if requestedProduct.Quantity > existingProduct.Quantity {
		return c.JSON(http.StatusNotAcceptable, responses.UserResponse{Status: http.StatusNotAcceptable, Message: "sufficient number of products doesn't exist in DB", Data: &echo.Map{"data": ""}})
	}
	if existingProduct.Price*float32(requestedProduct.Quantity) > buyer.Balance {
		return c.JSON(http.StatusNotAcceptable, responses.UserResponse{Status: http.StatusNotAcceptable, Message: "user does not have sufficient balance", Data: &echo.Map{"data": ""}})
	} else {
		//decrease balance of buyer
		buyer.Balance -= existingProduct.Price * float32(requestedProduct.Quantity)
		_, err := userCollection.UpdateOne(ctx, bson.M{"username": buyer.Username}, bson.M{"$set": buyer})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
		//increase balance of seller
		var seller models.User
		sellerID := existingProduct.SellerID
		err3 := userCollection.FindOne(ctx, bson.M{"_id": sellerID, "status": "Active"}).Decode(&seller)
		if err3 != nil {
			return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "cannot buy product, seller doesn't exists in DB", Data: &echo.Map{"data": err3.Error()}})
		}
		seller.Balance += existingProduct.Price * float32(requestedProduct.Quantity)
		_, err4 := userCollection.UpdateOne(ctx, bson.M{"username": seller.Username}, bson.M{"$set": seller})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err4.Error()}})
		}

		//decrease quantity of products
		existingProduct.Quantity -= requestedProduct.Quantity
		result, err2 := productCollection.UpdateOne(ctx, bson.M{"_id": requestedProductObjectId}, bson.M{"$set": existingProduct})
		if err2 != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err2.Error()}})
		}

		return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "product buying successful", Data: &echo.Map{"data": result}})
	}
}
