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

func UpdateProduct(product models.Product) error {
	fmt.Println("Starting updateProduct DB")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "UPDATE products SET "

	query += tools.CreateQuery(query, "Prod_Title", "S", 0, 0, product.ProdTitle)
	query += tools.CreateQuery(query, "Prod_Description", "S", 0, 0, product.ProdDescription)
	query += tools.CreateQuery(query, "Prod_Price", "F", 0, product.ProdPrice, "")
	query += tools.CreateQuery(query, "Prod_CategoryId", "N", product.ProdCategId, 0, "")
	query += tools.CreateQuery(query, "Prod_Stock", "N", product.ProdStock, 0, "")
	query += tools.CreateQuery(query, "Prod_Path", "S", 0,0, product.ProdPath)

	query += "WHERE Prod_Id = " + strconv.Itoa(product.ProdId)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update product > Sucessfully")
	return nil
}
