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
)

var productCollection = configs.GetCollection(configs.DB, "products")

func DeleteAUser(c echo.Context) error {
	claims := helpers.GetClaimsFromJwt(c)
	usernameFromClaims := claims.Username

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	usernamefromURL := c.Param("username")
	if usernamefromURL != usernameFromClaims {
		return c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "unauthorised user", Data: &echo.Map{"data": ""}})
	}

	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": usernamefromURL, "status": "Active"}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "username doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}

	cur, err := productCollection.Find(ctx, bson.M{"sellerusername": user.Username})
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "this seller doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}
	if cur.RemainingBatchLength() > 0 {
		fmt.Println("cur is bigger than 0")
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var singleProduct models.Product
			err := cur.Decode(&singleProduct)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error decoding single data from cursor", Data: &echo.Map{"data": err.Error()}})
			}
			singleProduct.IsAvailable = false

			_, err5 := productCollection.UpdateOne(ctx, bson.M{"_id": singleProduct.ProductID}, bson.M{"$set": singleProduct})
			if err5 != nil {
				return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error deleting products of the user you want to delete", Data: &echo.Map{"data": err.Error()}})
			}

		}
	}

	user.Status = "Inactive"
	result, err := userCollection.UpdateOne(ctx, bson.M{"username": usernameFromClaims}, bson.M{"$set": user})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error deleting user", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "User successfully deleted!", Data: &echo.Map{"data": result}})
}
