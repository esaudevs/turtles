package handlers

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/esaudevs/turtles/auth"
)

func Handlers(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {
	fmt.Println("Handling " + path + " > " + method)

	id := request.PathParameters["id"]
	idn, _ := strconv.Atoi(id)

	isOk, statusCode, user := IsAuthValid(path, method, headers)
	if !isOk {
		return statusCode, user
	}

	switch path[0:4] {
	case "user":
		return UsersHandler(body, path, method, user, id, request)
	case "prod":
		return ProductsHandler(body, path, method, user, idn, request)
	case "stock":
		return StocksHandler(body, path, method, user, idn, request)
	case "addr":
		return AddressesHandler(body, path, method, user, idn, request)
	case "cate":
		return CategoriesHandler(body, path, method, user, idn, request)
	case "orde":
		return OrdersHandler(body, path, method, user, idn, request)
	}

	return 400, "Invalid method"
}

func UsersHandler(body string, path string, method string, user string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Invalid method"
}

func ProductsHandler(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Invalid method"
}

func StocksHandler(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Invalid method"
}

func AddressesHandler(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Invalid method"
}

func CategoriesHandler(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Invalid method"
}

func OrdersHandler(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Invalid method"
}

func IsAuthValid(path string, method string, headers map[string]string) (bool, int, string) {
	if (path == "product" && method == "GET") ||
		(path == "category" && method == "GET") {
		return true, 200, ""
	}

	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "Token required"
	}

	allOk, err, msg := auth.IsTokenValid(token)
	if !allOk {
		if err != nil {
			fmt.Println("Error in token " + err.Error())
			return false, 401, err.Error()
		} else {
			fmt.Println("Error in token " + msg)
			return false, 401, msg
		}
	}

	fmt.Println("Token OK")
	return true, 200, msg
}
