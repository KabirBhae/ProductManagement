package productHelpers

import (
	"ProductManagement/models"
	"ProductManagement/responses"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func ViewAllProducts(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//usernamefromURL := c.Param("username")

	defer cancel()

	//err := productCollection.Find(ctx, bson.M{"username": usernamefromURL}).Decode(&user)
	cur, err := productCollection.Find(ctx, bson.M{"quantity": bson.M{"$gt": 0}})

	var products []models.Product

	if err != nil {
		return c.JSON(http.StatusNoContent, responses.UserResponse{Status: http.StatusNoContent, Message: "username doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var singleProduct models.Product
		err := cur.Decode(&singleProduct)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error decoding single data from cursor", Data: &echo.Map{"data": err.Error()}})
		}
		products = append(products, singleProduct)
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": products}})
}

func ViewOwn(c echo.Context) error {
	usernamefromURL := c.Param("sellerUsername")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//usernamefromURL := c.Param("username")

	defer cancel()
	cur, err := productCollection.Find(ctx, bson.M{"sellerusername": usernamefromURL, "quantity": bson.M{"$gte": 0}})
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "this seller doesn't exists in DB", Data: &echo.Map{"data": err.Error()}})
	}
	if cur.RemainingBatchLength() < 1 {
		return c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNoContent, Message: "this seller doesn't exists in DB", Data: &echo.Map{"data": ""}})
	}

	var products []models.Product
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var singleProduct models.Product
		err := cur.Decode(&singleProduct)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error decoding single data from cursor", Data: &echo.Map{"data": err.Error()}})
		}
		products = append(products, singleProduct)
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": products}})
}
