package bd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

	query = tools.CreateQuery(query, "Prod_Title", "S", 0, 0, product.ProdTitle)
	query = tools.CreateQuery(query, "Prod_Description", "S", 0, 0, product.ProdDescription)
	query = tools.CreateQuery(query, "Prod_Price", "F", 0, product.ProdPrice, "")
	query = tools.CreateQuery(query, "Prod_CategoryId", "N", product.ProdCategId, 0, "")
	query = tools.CreateQuery(query, "Prod_Stock", "N", product.ProdStock, 0, "")
	query = tools.CreateQuery(query, "Prod_Path", "S", 0, 0, product.ProdPath)

	query += "WHERE Prod_Id = " + strconv.Itoa(product.ProdId)

	fmt.Println(query)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update product > Sucessfully")
	return nil
}

func DeleteProduct(id int) error {
	fmt.Println("Starting delete product")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "DELETE FROM products WHERE Prod_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete product > Sucessfully")
	return nil
}

func SelectProduct(product models.Product, choice string, page int, pageSize int, orderType string, orderField string) (models.ProductResp, error) {
	fmt.Println("Starting Select product DB")

	var resp models.ProductResp
	var products []models.Product

	err := DbConnect()
	if err != nil {
		return resp, err
	}

	defer Db.Close()

	var query string
	var queryCount string
	var where, limit string

	query = "SELECT Prod_Id, Prod_Title, Prod_Description, Prod_CreatedAt, Prod_Updated, Prod_Price, Prod_Path, Prod_CategoryId, Prod_Stock FROM products"
	queryCount = "SELECT count(*) as items FROM products"

	switch choice {
	case "P":
		where = "WHERE Prod_Id = " + strconv.Itoa(product.ProdId)
	case "S":
		where = "WHERE UCASE(CONCAT(Prod_Title, Prod_Description)) LIKE '%" + strings.ToUpper(product.ProdSearch) + "%' "
	case "C":
		where = "WHERE Prod_CategId = " + strconv.Itoa(product.ProdCategId)
	case "U":
		where = "WHERE UCASE(Prod_Path) LIKE '%" + strings.ToUpper(product.ProdPath) + "%' "
	case "K":
		join := "JOIN category ON Prod_CategoryId = Categ_Id AND Categ_Path LIKE '%" + strings.ToUpper(product.ProdCategPath) + "%' "
		query += join
		queryCount += join
	}

	queryCount += where

	var rows *sql.Rows
	rows, err = Db.Query(queryCount)
	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
		return resp, err
	}

	rows.Next()
	var regs sql.NullInt32
	err = rows.Scan(&regs)

	items := int(regs.Int32)

	if page > 0 {
		if items > pageSize {
			limit = " LIMIT " + strconv.Itoa(pageSize)

			if page > 1 {
				offset := pageSize * (page - 1)
				limit += " OFFSET " + strconv.Itoa(offset)
			}
		} else {
			limit = ""
		}
	}

	var orderBy string
	if len(orderField) > 0 {
		switch orderField {
		case "I":
			orderBy = " ORDER BY Prod_Id "
		case "T":
			orderBy = " ORDER BY Prod_Title "
		case "D":
			orderBy = " ORDER BY Prod_Description "
		case "F":
			orderBy = " ORDER BY Prod_CreatedAt "
		case "P":
			orderBy = " ORDER BY Prod_Price "
		case "S":
			orderBy = " ORDER BY Prod_Stock "
		case "C":
			orderBy = " ORDER BY Prod_CategoryId "
		}

		if orderType == "D" {
			orderBy += " DESC"
		}
	}

	query += where + orderBy + limit

	fmt.Println(query)

	rows, err = Db.Query(query)

	for rows.Next() {
		var product models.Product
		var prodId sql.NullInt32
		var prodTitle sql.NullString
		var prodDescription sql.NullString
		var prodCreatedAt sql.NullTime
		var prodUpdated sql.NullTime
		var prodPrice sql.NullFloat64
		var prodPath sql.NullString
		var prodCategoryId sql.NullInt32
		var prodStock sql.NullInt32

		err := rows.Scan(&prodId, &prodTitle, &prodDescription, &prodCreatedAt, &prodUpdated, &prodPrice, &prodPath, &prodCategoryId, &prodStock)
		if err != nil {
			return resp, err
		}

		product.ProdId = int(prodId.Int32)
		product.ProdTitle = prodTitle.String
		product.ProdDescription = prodDescription.String
		product.ProdCreatedAt = prodCreatedAt.Time.String()
		product.ProdUpdated = prodUpdated.Time.String()
		product.ProdPrice = prodPrice.Float64
		product.ProdPath = prodPath.String
		product.ProdCategId = int(prodCategoryId.Int32)
		product.ProdStock = int(prodStock.Int32)

		products = append(products, product)
	}

	resp.TotalItems = items
	resp.Data = products

	fmt.Println("Select product > Successfully")
	return resp, nil
}
