package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func MySQLDate() string {
	t := time.Now()
	return fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func EscapeString(s string) string {
	desc := strings.ReplaceAll(s, "'", "")
	desc = strings.ReplaceAll(desc, "\"", "")
	return desc
}

func CreateQuery(queryString string, fieldName string, typeField string, ValueN int, ValueF float64, ValueS string) string {
	if (typeField == "S" && len(ValueS) == 0) ||
		(typeField == "F" && ValueF == 0) ||
		(typeField == "N" && ValueN == 0) {
		return queryString
	}

	if !strings.HasSuffix(queryString, "SET ") {
		queryString += ", "
	}

	switch typeField {
	case "S":
		queryString += fieldName + " = '" + EscapeString(ValueS) + "'"
	case "N":
		queryString += fieldName + " = " + strconv.Itoa(ValueN)
	case "F":
		queryString += fieldName + " = " + strconv.FormatFloat(ValueF, 'e', -1, 64)
	}

	return queryString
}
