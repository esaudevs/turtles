package routers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/esaudevs/turtles/bd"
	"github.com/esaudevs/turtles/models"
)

func UpdateUser(body string, user string) (int, string) {
	var userModel models.User
	err := json.Unmarshal([]byte(body), &userModel)

	if err != nil {
		return 400, "Error in data received " + err.Error()
	}

	if len(userModel.UserFirstName) == 0 && len(userModel.UserLastName) == 0 {
		return 400, "UserFirstName or UserLastName are needed"
	}

	_, userFound := bd.UserExists(user)
	if !userFound {
		return 400, "There is no user with this userId " + user
	}

	err = bd.UpdateUser(userModel, user)
	if err != nil {
		return 400, "There was an error trying to update user data " + user + " " + err.Error()
	}

	return 200, "Update userOk"
}

func SelectUser(body string, user string) (int, string) {
	_, userFound := bd.UserExists(user)
	if !userFound {
		return 400, "There is no user with this userId " + user
	}

	row, err := bd.SelectUser(user)

	fmt.Println(row)
	if err != nil {
		return 400, "There was an error trying to get the user data " + err.Error()
	}

	respJson, err := json.Marshal(row)
	if err != nil {
		return 500, "Error trying to parse user data"
	}

	return 200, string(respJson)
}

func SelectUsers(body string, user string, request events.APIGatewayV2HTTPRequest) (int, string) {
	var page int
	if len(request.QueryStringParameters["page"]) == 0 {
		page = 1
	} else {
		page, _ = strconv.Atoi(request.QueryStringParameters["page"])
	}

	isAdmin, errorMessage := bd.IsUserAdmin(user)
	if !isAdmin {
		return 400, errorMessage
	}

	users, err := bd.SelectUsers(page)

	if err != nil {
		return 400, "There was an error trying to get the user list > " + err.Error()
	}

	respJson, err := json.Marshal(users)
	if err != nil {
		return 500, "Error trying to parse users to Json"
	}

	return 200, string(respJson)
}
