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

func UpdateProduct(body string, user string, id int) (int, string) {
	var product models.Product

	err := json.Unmarshal([]byte(body), &product)
	if err != nil {
		return 400, "Error in data received" + err.Error()
	}

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	product.ProdId = id

	queryError := bd.UpdateProduct(product)
	if queryError != nil {
		return 400, "Error trying to update product " + strconv.Itoa(product.ProdId) + " > " + queryError.Error()
	}

	return 200, "Update Ok"
}

func DeleteProduct(user string, id int) (int, string) {

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	queryError := bd.DeleteProduct(id)
	if queryError != nil {
		return 400, "Error trying to delete product " + strconv.Itoa(id) + " > " + queryError.Error()
	}

	return 200, "Delete Ok"
}