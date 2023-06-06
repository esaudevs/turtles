package bd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/esaudevs/turtles/models"
	_ "github.com/go-sql-driver/mysql"
)

func InsertAddress(addr models.Address, user string) error {
	fmt.Println("Starting Insert Address DB")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "INSERT INTO addresses (Add_UserId, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name)"
	query += "VALUES ('" + user + "', '" + addr.AddAddress + "', '" + addr.AddCity + "', '" + addr.AddState + "', '"
	query += addr.AddPostalCode + "', '" + addr.AddPhone + "', '" + addr.AddTitle + "', '" + addr.AddName + "')"

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Insert Address > Sucessfully")
	return nil
}

func AddressExists(user string, id int) (error, bool) {
	fmt.Println("Starting address exists")

	err := DbConnect()
	if err != nil {
		return err, false
	}

	defer Db.Close()

	query := "SELECT 1 FROM addresses WHERE Add_Id = " + strconv.Itoa(id) + " AND Add_UserId = '" + user + "'"
	fmt.Println(query)

	rows, err := Db.Query(query)
	if err != nil {
		return err, false
	}

	var value string
	rows.Next()
	rows.Scan(&value)

	fmt.Println("Address exists > Sucessfully")

	return nil, value == "1"
}

func UpdateAddress(addr models.Address) error {
	fmt.Println("Starting Update Address DB")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "UPDATE addresses SET "

	if addr.AddAddress != "" {
		query += "Add_Address = '" + addr.AddAddress + "'"
	}
	if addr.AddCity != "" {
		query += "Add_City = '" + addr.AddCity + "'"
	}
	if addr.AddName != "" {
		query += "Add_Name = '" + addr.AddName + "'"
	}
	if addr.AddPhone != "" {
		query += "Add_Phone = '" + addr.AddPhone + "'"
	}
	if addr.AddPostalCode != "" {
		query += "Add_PostalCode = '" + addr.AddPostalCode + "'"
	}
	if addr.AddState != "" {
		query += "Add_State = '" + addr.AddState + "'"
	}
	if addr.AddTitle != "" {
		query += "Add_Title = '" + addr.AddTitle + "'"
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE Add_Id = " + strconv.Itoa(addr.AddId)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Address > Sucessfully")
	return nil
}

func DeleteAddress(id int) error {
	fmt.Println("Starting delete address")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "DELETE FROM addresses WHERE Add_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Address > Sucessfully")
	return nil
}

func SelectAddress(user string) ([]models.Address, error) {
	fmt.Println("Starting Select Address DB")

	addresses := []models.Address{}

	err := DbConnect()
	if err != nil {
		return addresses, err
	}

	defer Db.Close()

	query := "SELECT Add_Id, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name FROM addresses WHERE Add_UserId = '" + user + "'"

	var rows *sql.Rows
	rows, err = Db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return addresses, err
	}

	defer rows.Close()

	for rows.Next() {
		var address models.Address
		var addId sql.NullInt16
		var addAddress sql.NullString
		var addCity sql.NullString
		var addState sql.NullString
		var addPostalCode sql.NullString
		var addPhone sql.NullString
		var addTitle sql.NullString
		var addName sql.NullString

		err := rows.Scan(&addId, &addAddress, &addCity, &addState, &addPostalCode, &addPhone, &addTitle, &addName)
		if err != nil {
			return addresses, err
		}

		address.AddId = int(addId.Int16)
		address.AddAddress = addAddress.String
		address.AddCity = addCity.String
		address.AddState = addState.String
		address.AddPostalCode = addPostalCode.String
		address.AddPhone = addPhone.String
		address.AddTitle = addTitle.String
		address.AddName = addName.String

		addresses = append(addresses, address)
	}

	fmt.Println("Select addresses > Sucessfully")
	return addresses, nil
}
