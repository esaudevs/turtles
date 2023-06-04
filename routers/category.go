package routers

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/esaudevs/turtles/bd"
	"github.com/esaudevs/turtles/models"
)

func InsertCategory(body string, user string) (int, string) {
	var category models.Category

	err := json.Unmarshal([]byte(body), &category)
	if err != nil {
		return 400, "Error in data received" + err.Error()
	}

	if len(category.CategName) == 0 {
		return 400, "Category name needed"
	}

	if len(category.CategPath) == 0 {
		return 400, "Category path needed"
	}

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	result, queryError := bd.InsertCategory(category)
	if queryError != nil {
		return 400, "An error happened when inserting the category" + category.CategName + " > " + queryError.Error()
	}

	return 200, "{ CategID: " + strconv.Itoa(int(result)) + " }"
}

func UpdateCategory(body string, user string, id int) (int, string) {
	var category models.Category

	err := json.Unmarshal([]byte(body), &category)
	if err != nil {
		return 400, "Error in data received" + err.Error()
	}

	if len(category.CategName) == 0 && len(category.CategPath) == 0 {
		return 400, "Category name or path are needed"
	}

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	category.CategID = id

	queryError := bd.UpdateCategory(category)
	if queryError != nil {
		return 400, "Error trying to update category " + strconv.Itoa(category.CategID) + " > " + err.Error()
	}

	return 200, "Update Ok"
}

func DeleteCategory(user string, id int) (int, string) {
	if id == 0 {
		return 400, "Category id is needed"
	}

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	err := bd.DeleteCategory(id)
	if err != nil {
		return 400, "Error trying to delete category " + strconv.Itoa(id) + " > " + err.Error()
	}

	return 200, "Delete Ok"
}

func SelectCategories(body string, request events.APIGatewayV2HTTPRequest) (int, string) {
	var err error
	var CategId int
	var Slug string

	if len(request.QueryStringParameters["categId"]) > 0 {
		CategId, err = strconv.Atoi(request.QueryStringParameters["categId"])
		if err != nil {
			return 500, "Error trying to convert a to int " + request.QueryStringParameters["categId"]
		}
	} else {
		if len(request.QueryStringParameters["slug"]) > 0 {
			Slug = request.QueryStringParameters["slug"]
		}
	}

	categories, queryErr := bd.SelectCategories(CategId, Slug)
	if queryErr != nil {
		return 400, "Error trying to get categories > " + queryErr.Error()
	}

	Categ, jsonErr := json.Marshal(categories)
	if jsonErr != nil {
		return 400, "Error trying to parse categories > " + queryErr.Error()
	}

	return 200, string(Categ)
}
