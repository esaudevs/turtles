package bd

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/esaudevs/turtles/models"
	"github.com/esaudevs/turtles/tools"
	_ "github.com/go-sql-driver/mysql"
)

func InsertProduct(product models.Product) (int64, error) {
	fmt.Println("Starting insert product")

	err := DbConnect()
	if err != nil {
		return 0, err
	}

	defer Db.Close()

	query := "INSERT INTO products (Prod_Title"

	if len(product.ProdDescription) > 0 {
		query += ", Prod_Description"
	}

	if product.ProdPrice > 0 {
		query += ", Prod_Price"
	}

	if product.ProdCategId > 0 {
		query += ", Prod_CategoryId"
	}

	if product.ProdStock > 0 {
		query += ", Prod_Stock"
	}

	if len(product.ProdPath) > 0 {
		query += ", Prod_Path"
	}

	query += ") VALUES ('" + tools.EscapeString(product.ProdTitle) + "'"

	if len(product.ProdDescription) > 0 {
		query += ", '" + tools.EscapeString(product.ProdDescription) + "'"
	}

	if product.ProdPrice > 0 {
		query += ", " + strconv.FormatFloat(product.ProdPrice, 'e', -1, 64)
	}

	if product.ProdCategId > 0 {
		query += ", " + strconv.Itoa(product.ProdCategId)
	}

	if product.ProdStock > 0 {
		query += ", " + strconv.Itoa(product.ProdStock)
	}

	if len(product.ProdPath) > 0 {
		query += ", '" + tools.EscapeString(product.ProdPath) + "'"
	}

	query += ")"

	var result sql.Result
	result, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, insertErr := result.LastInsertId()
	if insertErr != nil {
		return 0, insertErr
	}

	fmt.Println("Insert product > Succesfully completed")
	return LastInsertId, insertErr
}
