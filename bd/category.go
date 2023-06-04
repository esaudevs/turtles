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

func InsertCategory(category models.Category) (int64, error) {
	fmt.Println("Starting insertCategory DB")

	err := DbConnect()
	if err != nil {
		return 0, err
	}

	defer Db.Close()

	query := "INSERT INTO category (Categ_Name, Categ_Path) VALUES ('" + category.CategName + "','" + category.CategPath + "')"

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

	fmt.Println("Insert category > Succesfully completed")
	return LastInsertId, insertErr
}

func UpdateCategory(category models.Category) error {
	fmt.Println("Starting updateCategory DB")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "UPDATE category SET "

	if len(category.CategName) > 0 {
		query += "Categ_Name = '" + tools.EscapeString(category.CategName) + "'"
	}

	if len(category.CategPath) > 0 {
		if !strings.HasSuffix(query, "SET") {
			query += ", "
		}

		query += "Categ_Path = '" + tools.EscapeString(category.CategPath) + "'"
	}

	query += " WHERE Categ_Id = " + strconv.Itoa(category.CategID)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Category > Sucessfully")
	return nil
}

func DeleteCategory(id int) error {
	fmt.Println("Starting deleteCategory DB")

	err := DbConnect()
	if err != nil {
		return err
	}

	defer Db.Close()

	query := "DELETE FROM category WHERE Categ_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Category > Sucessfully")
	return nil
}

func SelectCategories(CategId int, Slug string) ([]models.Category, error) {
	fmt.Println("Starting SelectCategories")

	var Categories []models.Category

	err := DbConnect()
	if err != nil {
		return Categories, err
	}

	defer Db.Close()

	query := "SELECT Categ_Id, Categ_Name, Categ_Path FROM category"

	if CategId > 0 {
		query += " WHERE Categ_Id = " + strconv.Itoa(CategId)
	} else {
		if len(Slug) > 0 {
			query += " WHERE Categ_Path LIKE '%" + Slug + "%'"
		}
	}

	fmt.Println(query)

	var rows *sql.Rows
	rows, err = Db.Query(query)

	for rows.Next() {
		var category models.Category
		var categId sql.NullInt32
		var categName sql.NullString
		var categPath sql.NullString

		err := rows.Scan(&categId, &categName, &categPath)
		if err != nil {
			return Categories, err
		}
		category.CategID = int(categId.Int32)
		category.CategName = categName.String
		category.CategPath = categPath.String

		Categories = append(Categories, category)
	}

	fmt.Println("Select Categories > Successfully")

	return Categories, nil
}
