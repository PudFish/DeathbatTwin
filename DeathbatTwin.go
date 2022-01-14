package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

//errors
var (
	ErrInvalidTokenId   = errors.New("invalid token id, must be between 1 and 10000 inclusive")
	ErrDeathbatNotFound = errors.New("Deathbat not found")
)

//apis
const (
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

//OpenSeaDeathbat represents a partial structure of a Deathbat as listed on OpenSea.io
type OpenSeaDeathbat struct {
	Owner struct {
		User struct {
			Username string `json:"username"`
		} `json:"user"`
	} `json:"owner"`
}

//Deathbats is the global memory storage for all loaded deathbats
var Deathbats []Deathbat

//main handles the high level function calls for now
func main() {
	filename := "deathbats1-1000.json"
	if err := loadDeathbats(filename); err != nil {
		fmt.Printf("err: main: %s", err)
		return
	}

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

	if err = sourceDeathbat.loadOwner(); err != nil {
		fmt.Printf("err: main: %s", err)
	}

	fmt.Printf("Source Deathbat: ")
	sourceDeathbat.print()

	//TODO: implement findTwin
	twinDeathbat, err := findTwin(sourceDeathbat)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	if err = twinDeathbat.loadOwner(); err != nil {
		fmt.Printf("err: main: %s", err)
	}

	fmt.Printf("\nTwin Deathbat: ")
	twinDeathbat.print()
}

//loadDeathbats reads the Deathbats json file and loads it to memory
func loadDeathbats(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("loadDeathbats: %w", err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("loadDeathbats: %w", err)
	}

	if err = json.Unmarshal(data, &Deathbats); err != nil {
		return fmt.Errorf("loadDeathbats: %w", err)
	}

	for i, deathbat := range Deathbats {
		Deathbats[i].Hyperlink = openSeaLink + strconv.Itoa(deathbat.Id)
	}

	err = file.Close()
	return err
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

//getDeathbat retrieves a Deathbat from memory
func getDeathbat(tokenId int) (deathbat Deathbat, err error) {
	//check memory
	for _, deathbat := range Deathbats {
		if deathbat.Id == tokenId {
			return deathbat, nil
		}
	}

	return Deathbat{}, fmt.Errorf("getDeathbat: %s: %w", ErrDeathbatNotFound, err)
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

//loadOwner checks who the current owner of the Deathbat is on OpenSea.io
func (deathbat *Deathbat) loadOwner() (err error) {
	deathbat.Owner = "Unknown"

	URL := ownerAPI + strconv.Itoa(deathbat.Id)

	response, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf("loadOwner: %w", err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("loadOwner: %w", err)
	}

	var jsonData OpenSeaDeathbat
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("loadOwner: %w", err)
	}

	deathbat.Owner = jsonData.Owner.User.Username

	return nil
}

//findTwin finds and returns another Deathbat most alike the provided Deathbat
func findTwin(deathbat Deathbat) (twin Deathbat, err error) {
	return deathbat, nil
}
