package routers

import (
	"encoding/json"
	"strconv"

	"github.com/esaudevs/turtles/bd"
	"github.com/esaudevs/turtles/models"
)

func InsertProduct(body string, user string) (int, string) {
	var product models.Product
	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return 400, "Error in product data"
	}

	if len(product.ProdTitle) == 0 {
		return 400, "ProductTitle is needed"
	}

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	result, queryError := bd.InsertProduct(product)
	if queryError != nil {
		return 400, "Error inserting product in DB " + queryError.Error()
	}

	return 200, "{ ProductID: " + strconv.Itoa(int(result)) + " }"
}
