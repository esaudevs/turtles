package routers

import (
	"bytes"
	"encoding/json"
	"strconv"

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

	stringResult := jsonPrettyPrint("{ CategID: " + strconv.Itoa(int(result)) + " }")

	return 200, string(stringResult)
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
