package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type TokenJSON struct {
	Sub       string
	Event_Id  string
	Token_use string
	Scope     string
	Auth_time int
	Iss       string
	Exp       int
	Iat       int
	Client_id string
	Username  string
}

func IsTokenValid(token string) (bool, error, string) {
	parts := strings.Split(token, ".")
	
	if len(parts) != 3 {
		fmt.Println("Invalid token")
		return false, nil, "Invalid token"
	}

	userInfo, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Cannot decode token: " + err.Error())
		return false, err, err.Error()
	}

	var tokenJson TokenJSON
	err = json.Unmarshal(userInfo, &tokenJson)
	if err != nil {
		fmt.Println("Cannot parse json: " + err.Error())
		return false, err, err.Error()
	}

	timeNow := time.Now()
	tokenExpiration := time.Unix(int64(tokenJson.Exp), 0)

	if tokenExpiration.Before(timeNow) {
		fmt.Println("Token expired on: " + tokenExpiration.String())
		fmt.Println("Token expired")
		return false, err, "Token expired"
	}

	return true, nil, string(tokenJson.Username)
}