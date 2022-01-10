package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

//errors
var (
	ErrInvalidTokenId   = errors.New("invalid token id, must be between 1 and 10000 inclusive")
	ErrDeathbatNotFound = errors.New("Deathbat not found")
)

//apis
const (
	traitsAPI   = "https://avengedsevenfold.io/deathbats/token/"
	ownerAPI    = "https://api.opensea.io/api/v1/asset/0x1D3aDa5856B14D9dF178EA5Cab137d436dC55F1D/"
	openSeaLink = "https://opensea.io/assets/0x1d3ada5856b14d9df178ea5cab137d436dc55f1d/"
)

//Deathbat represents the properties of a Deathbat
type Deathbat struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	Minted      bool        `json:"minted"`
	Image       string      `json:"image"`
	Attributes  []struct {
		TraitType string `json:"trait_type"`
		Value     string `json:"value"`
	} `json:"attributes"`
	Owner     string `json:"owner:"`
	Hyperlink string `json:"hyperlink"`
}

//main handles the high level function calls for now
func main() {
	tokenId, err := getSourceDeathbat()
	if err != nil {
		fmt.Printf("err: main: %s", err)
		return
	}

	if err = checkTokenId(tokenId); err != nil {
		fmt.Printf("err: main: %s", err)
		return
	}

	sourceDeathbat, err := getDeathbat(tokenId)
	if err != nil {
		fmt.Printf("err: main: %s", err)
		return
	}

	sourceDeathbat.print()
	_ = ownerAPI
}

//getSourceDeathbat prompts the user for a Deathbat tokenId to use as the source for comparison
func getSourceDeathbat() (tokenId int, err error) {
	fmt.Printf("What number Deathbat do you want to find a twin for? ")
	var input string
	if _, err = fmt.Scanln(&input); err != nil {
		return 0, fmt.Errorf("getSourceDeathbat: %w", err)
	}

	tokenId, err = strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("getSourceDeathbat: %w", err)
	}

	return tokenId, nil
}

//checkTokenId checks if a tokenId provided is within the valid range
func checkTokenId(tokenId int) (err error) {
	valid := tokenId >= 1 && tokenId <= 10000
	if !valid {
		return ErrInvalidTokenId
	}

	return nil
}

//getDeathbat retrieves a Deathbat from avengedsevenfold.io
func getDeathbat(tokenId int) (deathbat Deathbat, err error) {
	resp, err := http.Get(traitsAPI + strconv.Itoa(tokenId))
	if err != nil {
		return deathbat, fmt.Errorf("getDeathbat: %s: %w", ErrDeathbatNotFound, err)
	}
	if resp.StatusCode != http.StatusOK {
		return deathbat, fmt.Errorf("getDeathbat: %s: response status code %d", ErrDeathbatNotFound, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&deathbat); err != nil {
		return deathbat, fmt.Errorf("getDeathbat: %s: %w", ErrDeathbatNotFound, err)
	}

	deathbat.Hyperlink = openSeaLink + strconv.Itoa(tokenId)

	return deathbat, nil
}

//print displays the deathbat details in a pretty format
func (deathbat *Deathbat) print() {
	traits := ""
	for i, trait := range deathbat.Attributes {
		traits += fmt.Sprintf("%s: %s", trait.TraitType, trait.Value)
		if i != len(deathbat.Attributes)-1 {
			traits += ", "
		}
	}

	fmt.Printf("Deathbat #%d\n%s\nOwner: %s\nOpenSea.io link: %s\n", deathbat.Id, traits, deathbat.Owner, deathbat.Hyperlink)
}
