package bd

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/esaudevs/turtles/models"
	"github.com/esaudevs/turtles/tools"
	_ "github.com/go-sql-driver/mysql"
)

func UpdateUser(userModel models.User, user string) error {
	fmt.Println("Starting UpdateUser DB")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "UPDATE users SET "

	colon := ""
	if len(userModel.UserFirstName) > 0 {
		colon = ","
		query += "User_FirstName = '" + userModel.UserFirstName + "'"
	}

	if len(userModel.UserLastName) > 0 {
		query += colon + "User_LastName = '" + userModel.UserLastName + "'"
	}

	query += ", User_DateUpg = '" + tools.MySQLDate() + "' WHERE User_UUID = '" + user + "'"

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update user > Sucessfully")
	return nil
}

func SelectUser(userId string) (models.User, error) {
	fmt.Println("Starting SelectUser DB")

	userModel := models.User{}

	err := DbConnect()
	if err != nil {
		return userModel, err
	}

	defer Db.Close()

	query := "SELECT * FROM users WHERE User_UUID = '" + userId + "'"

	var rows *sql.Rows
	rows, err = Db.Query(query)

	if err != nil {
		fmt.Println(err.Error())
		return userModel, err
	}

	defer rows.Close()

	rows.Next()

	var firstName sql.NullString
	var lastName sql.NullString
	var dateUpg sql.NullTime

	rows.Scan(&userModel.UserUUID, &userModel.UserEmail, &firstName, &lastName, &userModel.UserStatus, &userModel.UserDateAdd, &dateUpg)

	userModel.UserFirstName = firstName.String
	userModel.UserLastName = lastName.String
	userModel.UserDateUpd = dateUpg.Time.String()

	fmt.Println("Select User > Sucessfully")

	return userModel, nil
}

func SelectUsers(page int) (models.UsersList, error) {
	fmt.Println("Starting Select Users DB")

	var usersListResp models.UsersList
	usersList := []models.User{}

	err := DbConnect()
	if err != nil {
		return usersListResp, err
	}

	defer Db.Close()

	var offset int = (page + 10) - 10
	var query string
	var queryCount string = "SELECT count(*) FROM users"

	query = "SELECT * FROM users LIMIT 10"
	if offset > 0 {
		query = " OFFSET " + strconv.Itoa(offset)
	}

	var rowsCount *sql.Rows
	rowsCount, err = Db.Query(queryCount)
	if err != nil {
		return usersListResp, err
	}

	defer rowsCount.Close()

	rowsCount.Next()
	var totalItems int
	rowsCount.Scan(&totalItems)
	usersListResp.TotalItems = totalItems

	var rows *sql.Rows
	rows, err = Db.Query(query)

	if err != nil {
		return usersListResp, err
	}

	for rows.Next() {
		var userItem models.User

		var firstName sql.NullString
		var lastName sql.NullString
		var dateUpg sql.NullTime

		rows.Scan(&userItem.UserUUID, &userItem.UserEmail, &firstName, &lastName, &userItem.UserStatus, &userItem.UserDateAdd, &dateUpg)

		userItem.UserFirstName = firstName.String
		userItem.UserLastName = lastName.String
		userItem.UserDateUpd = dateUpg.Time.String()
		usersList = append(usersList, userItem)
	}

	fmt.Println("Select users > Sucessfully")

	usersListResp.Data = usersList
	return usersListResp, nil
}
