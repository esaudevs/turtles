package bd

import (
	"database/sql"
	"fmt"

	"github.com/esaudevs/turtles/models"
	_ "github.com/go-sql-driver/mysql"
)

func InsertCategory(category models.Category) (int64, error) {
	fmt.Println("Starting insertCategory DB")

	err := DbConnect()
	if err != nil {
		return 0, err
	}

	defer Db.Close()

	query := "INSER INTO category (Categ_Name, Categ_Path) VALUES ('" + category.CategName + "','" + category.CategPath + "')"

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
