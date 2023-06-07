package routers

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/esaudevs/turtles/bd"
	"github.com/esaudevs/turtles/models"
)

func InsertOrder(body string, user string) (int, string) {
	var order models.Order

	err := json.Unmarshal([]byte(body), &order)
	if err != nil {
		return 400, "Error in data received " + err.Error()
	}

	order.Order_UserUUID = user

	Ok, message := IsOrderValid(order)

	if !Ok {
		return 400, message
	}

	result, queryErr := bd.InsertOrder(order)
	if queryErr != nil {
		return 400, "There was an error trying to insert order " + queryErr.Error()
	}

	return 200, "{ OrderID: " + strconv.Itoa(int(result)) + " }"
}

func SelectOrders(user string, request events.APIGatewayV2HTTPRequest) (int, string) {
	var fromDate, untilDate string
	var orderId int
	var page int

	if len(request.QueryStringParameters["fromDate"]) > 0 {
		fromDate = request.QueryStringParameters["fromDate"]
	}
	if len(request.QueryStringParameters["untilDate"]) > 0 {
		untilDate = request.QueryStringParameters["untilDate"]
	}
	if len(request.QueryStringParameters["page"]) > 0 {
		page, _ = strconv.Atoi(request.QueryStringParameters["page"])
	}
	if len(request.QueryStringParameters["orderId"]) > 0 {
		orderId, _ = strconv.Atoi(request.QueryStringParameters["orderId"])
	}

	result, queryErr := bd.SelectOrders(user, fromDate, untilDate, page, orderId)
	if queryErr != nil {
		return 400, "There was an error trying to get orders"
	}

	orders, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		return 400, "There was an error trying to parse to Json"
	}

	return 200, string(orders)

}

func IsOrderValid(order models.Order) (bool, string) {
	if order.Order_Total == 0 {
		return false, "Order total field is needed"
	}

	count := 0

	for _, od := range order.OrderDetails {
		if od.OD_ProdId == 0 {
			return false, "ProductId is needed"
		}
		if od.OD_Quantity == 0 {
			return false, "Quantity is needed"
		}
		count++
	}
	if count == 0 {
		return false, "Items are needed to create an order"
	}

	return true, ""
}
