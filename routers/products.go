package routers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

func SelectProduct(request events.APIGatewayV2HTTPRequest) (int, string) {
	var product models.Product
	var page, pageSize int
	var orderType, orderField string

	param := request.QueryStringParameters

	page, _ = strconv.Atoi(param["page"])
	pageSize, _ = strconv.Atoi(param["pageSize"])
	orderType = param["orderType"] // D = DESC, A or Nil = ASC
	orderField = param["orderField"] // I Id, T Title, D Desc, F CreatedAt, P Price, C categId, S Stock

	if !strings.Contains("ITDFPCS", orderField) {
		orderField = ""
	}

	var choice string
	if len(param["prodId"]) > 0 {
		choice="P"
		product.ProdId, _ = strconv.Atoi(param["prodId"])
	}
	if len(param["search"]) > 0 {
		choice="S"
		product.ProdSearch, _ = param["search"]
	}
	if len(param["categId"]) > 0 {
		choice="C"
		product.ProdCategId, _ = strconv.Atoi(param["categId"])
	}
	if len(param["slug"]) > 0 {
		choice="U"
		product.ProdPath, _ = param["slug"]
	}
	if len(param["slugCateg"]) > 0 {
		choice="K"
		product.ProdCategPath, _ = param["slugCateg"]
	}

	fmt.Println(param)

	result, queryErr := bd.SelectProduct(product, choice, page, pageSize, orderType, orderField)
	if queryErr != nil {
		return 400, "Error trying to get result of type '" + choice + "' " + queryErr.Error() 
	}

	productList, formatErr := json.Marshal(result)
	if formatErr != nil {
		return 400, "Error trying to parse products to Json"
	}

	return 200, string(productList)
}

func UpdateProductStock(body string, user string, id int) (int, string) {
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

	queryError := bd.UpdateProductStock(product)
	if queryError != nil {
		return 400, "Error trying to update product stock " + strconv.Itoa(product.ProdId) + " > " + queryError.Error()
	}

	return 200, "Update Ok"
}