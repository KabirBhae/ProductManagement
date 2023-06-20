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

func DeleteProduct(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	productIdFromURL := c.Param("productID")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productObjectID, errr := primitive.ObjectIDFromHex(productIdFromURL)
	if errr != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "string cannot be parsed to ObjectID", Data: &echo.Map{"data": errr.Error()}})
	}

	var existingProduct models.Product
	err := productCollection.FindOne(ctx, bson.M{"_id": productObjectID, "isavailable": true}).Decode(&existingProduct)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "product doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}
	if usernameFromClaims != existingProduct.SellerUserName {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "cannot delete other sellers' product", Data: &echo.Map{"data": ""}})
	}

	existingProduct.IsAvailable = false
	result, err := productCollection.UpdateOne(ctx, bson.M{"_id": existingProduct.ProductID}, bson.M{"$set": existingProduct})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error deleting product", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Product successfully deleted!", Data: &echo.Map{"data": result}})
}
