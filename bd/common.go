package bd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/esaudevs/turtles/models"
	"github.com/esaudevs/turtles/secrets"
	_ "github.com/go-sql-driver/mysql"
)

var SecretModel models.SecretRDSJson
var err error
var Db *sql.DB

func ReadSecret() error {
	SecretModel, err = secrets.GetSecret(os.Getenv("SecretName"))
	return err
}

func DbConnect() error {
	Db, err = sql.Open("mysql", ConnStr(SecretModel))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = Db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Connected to DB successfully")
	return nil
}

func ConnStr(key models.SecretRDSJson) string {
	var dbUser, authToken, dbEndpoint, dbName string
	dbUser = key.Username
	authToken = key.Password
	dbEndpoint = key.Host
	dbName = "turtles"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&parseTime=true", dbUser, authToken, dbEndpoint, dbName)
	fmt.Println(dsn)
	return dsn
}

func IsUserAdmin(userUUID string) (bool, string) {
	fmt.Println("Starting IsUserAdmin")

	err := DbConnect()
	if err != nil {
		return false, err.Error()
	}

	defer Db.Close()

	query := "SELECT 1 FROM users WHERE User_UUID = '" + userUUID + "' AND User_Status = 0"
	fmt.Println(query)

	rows, err := Db.Query(query)
	if err != nil {
		return false, err.Error()
	}

	var value string
	rows.Next()
	rows.Scan(&value)

	fmt.Println("IsUserAdmin > Sucessfully - value returned: " + value)
	if value == "1" {
		return true, ""
	}

	return false, "User has not admin privileges"
}
