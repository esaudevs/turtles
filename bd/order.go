package bd

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/esaudevs/turtles/models"
	_ "github.com/go-sql-driver/mysql"
)

func InsertOrder(order models.Order) (int64, error) {
	fmt.Println("Starting Insert Order DB")

	err := DbConnect()
	if err != nil {
		return 0, err
	}

	defer Db.Close()

	query := "INSERT INTO orders (Order_UserUUID, Order_Total, Order_AddId) VALUES ('"
	query += order.Order_UserUUID + "', " + strconv.FormatFloat(order.Order_Total, 'f', -1, 64) + "," + strconv.Itoa(order.Order_AddId) + ")"

	var result sql.Result
	result, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, errLastInsert := result.LastInsertId()
	if errLastInsert != nil {
		return 0, errLastInsert
	}

	for _, od := range order.OrderDetails {
		query = "INSERT INTO orders_detail (OD_OrderId, OD_ProdId, OD_Quantity, OD_Price) VALUES (" + strconv.Itoa(int(LastInsertId))
		query += "," + strconv.Itoa(od.OD_ProdId) + "," + strconv.Itoa(od.OD_Quantity) + "," + strconv.FormatFloat(od.OD_Price, 'f', -1, 64) + ")"
		fmt.Println(query)

		_, err = Db.Exec(query)
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
	}

	fmt.Println("Insert Order > Successfully")
	return LastInsertId, nil
}

func SelectOrders(user string, fromDate string, untilDate string, page int, orderId int) ([]models.Order, error) {
	fmt.Println("Starting Select Orders DB")

	var orders []models.Order

	query := "SELECT Order_Id, Order_UserUUID, Order_AddId, Order_Date, Order_Total FROM orders "

	if orderId > 0 {
		query += " WHERE Order_Id = " + strconv.Itoa(orderId)
	} else {
		offset := 0
		if page == 0 {
			page = 1
		}
		if page > 1 {
			offset = (10 * (page - 1))
		}

		if len(untilDate) == 10 {
			untilDate += " 23:59:59"
		}
		var where string
		var whereUser string = " Order_UserUUID = '" + user + "'"

		if len(fromDate) > 0 && len(untilDate) > 0 {
			where += " WHERE Order_Date BETWEEN '" + fromDate + "' AND '" + untilDate + "' "
		}
		if len(where) > 0 {
			where += " AND " + whereUser
		} else {
			where += " WHERE " + whereUser
		}

		limit := " LIMIT 10 "
		if offset > 0 {
			limit += " OFFSET " + strconv.Itoa(offset)
		}

		query += where + limit
	}

	fmt.Println(query)

	err := DbConnect()
	if err != nil {
		return orders, err
	}

	defer Db.Close()

	var rows *sql.Rows
	rows, err = Db.Query(query)
	if err != nil {
		return orders, err
	}

	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var orderAddId sql.NullInt32

		err := rows.Scan(&order.Order_Id, &order.Order_UserUUID, &orderAddId, &order.Order_Date, &order.Order_Total)
		if err != nil {
			return orders, err
		}

		order.Order_AddId = int(orderAddId.Int32)

		var rowsDetails *sql.Rows
		queryDetails := "SELECT OD_Id, OD_ProdId, OD_Quantity, OD_Price FROM orders_detail WHERE OD_OrderId = " + strconv.Itoa(order.Order_Id)
		rowsDetails, err = Db.Query(queryDetails)
		if err != nil {
			return orders, err
		}

		for rowsDetails.Next() {
			var oD_Id int64
			var oD_ProdId int64
			var oD_Quantity int64
			var oD_Price float64

			err = rowsDetails.Scan(&oD_Id, &oD_ProdId, &oD_Quantity, &oD_Price)
			if err != nil {
				return orders, err
			}

			var orderDetail models.OrderDetail
			orderDetail.OD_Id = int(oD_Id)
			orderDetail.OD_ProdId = int(oD_ProdId)
			orderDetail.OD_Quantity = int(oD_Quantity)
			orderDetail.OD_Price = oD_Price

			order.OrderDetails = append(order.OrderDetails, orderDetail)
		}

		orders = append(orders, order)
		rowsDetails.Close()
	}

	fmt.Println("Select Orders > Successfully")
	return orders, nil
}
