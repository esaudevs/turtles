package routers

import (
	"encoding/json"

	"github.com/esaudevs/turtles/bd"
	"github.com/esaudevs/turtles/models"
)

func InsertAddress(body string, user string) (int, string) {
	var address models.Address

	err := json.Unmarshal([]byte(body), &address)
	if err != nil {
		return 400, "Error in data received " + err.Error()
	}

	if address.AddAddress == "" {
		return 400, "Address field is needed"
	}
	if address.AddName == "" {
		return 400, "Address name is needed"
	}
	if address.AddTitle == "" {
		return 400, "Address title is needed"
	}
	if address.AddCity == "" {
		return 400, "Address city is needed"
	}
	if address.AddPhone == "" {
		return 400, "Address phone is needed"
	}
	if address.AddPostalCode == "" {
		return 400, "Address postal code is needed"
	}

	err = bd.InsertAddress(address, user)
	if err != nil {
		return 400, "There was an error trying to insert the address " + err.Error()
	}

	return 200, "Insert Address Ok"
}

func UpdateAddress(body string, user string, id int) (int, string) {
	var address models.Address

	err := json.Unmarshal([]byte(body), &address)
	if err != nil {
		return 400, "Error in data received"
	}

	address.AddId = id

	var addressFound bool
	err, addressFound = bd.AddressExists(user, address.AddId)
	if !addressFound {
		if err != nil {
			return 400, "Error trying to get address for this user " + err.Error()
		}

		return 400, "There is not address with this id asociated to this user"
	}

	err = bd.UpdateAddress(address)
	if err != nil {
		return 400, "There was an error trying to update address for this user " + err.Error()
	}

	return 200, "Update Address Ok"
}

func DeleteAddress(user string, id int) (int, string) {
	err, addressFound := bd.AddressExists(user, id)
	if !addressFound {
		if err != nil {
			return 400, "Error trying to get address for this user " + err.Error()
		}

		return 400, "There is not address with this id asociated to this user"
	}

	err = bd.DeleteAddress(id)
	if err != nil {
		return 400, "There was an error trying to delete the address " + err.Error()
	}

	return 200, "Delete Address Ok"
}

func SelectAddress(user string) (int, string) {
	addr, err := bd.SelectAddress(user)
	if err != nil {
		return 400, "There was an error trying to get user addresses " + err.Error()
	}

	respJson, err := json.Marshal(addr)
	if err != nil {
		return 500, "Error trying to parse addresses to json " + err.Error()
	}

	return 200, string(respJson)
}
