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
	ErrInvalidTokenId      = errors.New("invalid token id, must be between 1 and 10000 inclusive")
	ErrDeathbatNotFound    = errors.New("Deathbat not found")
	ErrUnknownTraitType    = errors.New("unknown trait type")
	ErrOpenSeaUnresponsive = errors.New("opensea.io unresponsive")
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
	Traits struct {
		Background      string `json:"background,omitempty"`
		BrooksWackerman string `json:"brooks_wackerman,omitempty"`
		Eyes            string `json:"eyes,omitempty"`
		FacialHair      string `json:"facial_hair,omitempty"`
		Head            string `json:"head,omitempty"`
		JohnnyChrist    string `json:"johnny_christ,omitempty"`
		Mask            string `json:"mask,omitempty"`
		Mouth           string `json:"mouth,omitempty"`
		Shadows         string `json:"shadows,omitempty"`
		Nose            string `json:"nose,omitempty"`
		Perk            string `json:"perk,omitempty"`
		Skin            string `json:"skin,omitempty"`
		SynysterGates   string `json:"synyster_gates,omitempty"`
		ZackyVengeance  string `json:"zacky_vengeance,omitempty"`
	} `json:"traits"`
	Hyperlink string `json:"hyperlink"`
	Owner     string `json:"owner"`
}

//Deathbats is the global memory storage for all loaded deathbats
var Deathbats []Deathbat

//main handles the high level function calls for now
func main() {
	filename := "deathbats.json"
	if err := loadDeathbats(filename); err != nil {
		fmt.Printf("err: main: %s\n", err)
		return
	}

	tokenId, err := getSourceDeathbat()
	if err != nil {
		fmt.Printf("err: main: %s\n", err)
		return
	}

	if err = checkTokenId(tokenId); err != nil {
		fmt.Printf("err: main: %s\n", err)
		return
	}

	sourceDeathbat, err := getDeathbat(tokenId)
	if err != nil {
		fmt.Printf("err: main: %s\n", err)
		return
	}

	if err = sourceDeathbat.loadOwner(); err != nil {
		fmt.Printf("err: main: %s\n", err)
	}

	fmt.Printf("Source Deathbat: ")
	sourceDeathbat.print()

	twinDeathbat, err := findTwin(sourceDeathbat)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if err = twinDeathbat.loadOwner(); err != nil {
		fmt.Printf("err: main: %s\n", err)
	}

	fmt.Printf("\nTwin Deathbat: ")
	twinDeathbat.print()
}

//loadDeathbats reads the Deathbats json file and loads it to memory
func loadDeathbats(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("loadDeathbats: Open: %w", err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("loadDeathbats: ReadAll: %w", err)
	}

	if err = json.Unmarshal(data, &Deathbats); err != nil {
		return fmt.Errorf("loadDeathbats: Unmarshal: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("loadDeathbats: %w", err)
	}
	return nil
}

//getSourceDeathbat prompts the user for a Deathbat tokenId to use as the source for comparison
func getSourceDeathbat() (tokenId int, err error) {
	fmt.Printf("What number Deathbat do you want to find a twin for? ")
	var input string
	if _, err = fmt.Scanln(&input); err != nil {
		return 0, fmt.Errorf("getSourceDeathbat: Scanln: %w", err)
	}

	tokenId, err = strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("getSourceDeathbat: Atoi: %w", err)
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
	//quick check if ordered
	if Deathbats[tokenId-1].Id == tokenId {
		return Deathbats[tokenId-1], nil
	}

	//check all memory in case unordered
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
	ownerAPI := "https://api.opensea.io/api/v1/asset/0x1D3aDa5856B14D9dF178EA5Cab137d436dC55F1D/"

	URL := ownerAPI + strconv.Itoa(deathbat.Id)

	response, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf("loadOwner: Get: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return ErrOpenSeaUnresponsive
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("loadOwner: ReadAll: %w", err)
	}

	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("loadOwner: Unmarshal: %w", err)
	}

	//TODO: pull out owner data

	return nil
}

//findTwin finds and returns another Deathbat most alike the provided Deathbat
func findTwin(deathbat Deathbat) (twin Deathbat, err error) {
	twin = deathbat

	//weights to determine which matching traits matter more
	weight := map[string]int{
		"Mask":        6,
		"Facial Hair": 5,
		"Eyes":        4,
		"Mouth":       4,
		"Nose":        4,
		"Head":        3,
		"Skin":        2,
		"Background":  1,
	}

	//1/1 (Brooks Wackerman, Johnny Christ, M. Shadows, Synyser Gates, Zacky Vengeance)
	if deathbat.Traits.BrooksWackerman != "" || deathbat.Traits.JohnnyChrist != "" || deathbat.Traits.Shadows != "" || deathbat.Traits.SynysterGates != "" || deathbat.Traits.ZackyVengeance != "" {
		return twin, fmt.Errorf("you got a 1/1, there is no twin")
	}

	score := 0
	for _, check := range Deathbats {
		//can't be the same deathbat
		if deathbat.Id == check.Id {
			continue
		}

		checkScore := 0

		//Mask
		if deathbat.Traits.Mask != "" && deathbat.Traits.Mask == check.Traits.Mask {
			checkScore += weight["Mask"]
		}

		//Facial Hair
		if deathbat.Traits.FacialHair != "" && deathbat.Traits.FacialHair == check.Traits.FacialHair {
			checkScore += weight["Facial Hair"]
		}

		//Eyes
		if deathbat.Traits.Eyes != "" && deathbat.Traits.Eyes == check.Traits.Eyes {
			checkScore += weight["Eyes"]
		}

		//Mouth
		if deathbat.Traits.Mouth != "" && deathbat.Traits.Mouth == check.Traits.Mouth {
			checkScore += weight["Mouth"]
		}

		//Nose
		if deathbat.Traits.Nose != "" && deathbat.Traits.Nose == check.Traits.Nose {
			checkScore += weight["Nose"]
		}

		//Head
		if deathbat.Traits.Head != "" && deathbat.Traits.Head == check.Traits.Head {
			checkScore += weight["Head"]
		}

		//Skin
		if deathbat.Traits.Skin != "" && deathbat.Traits.Skin == check.Traits.Skin {
			checkScore += weight["Skin"]
		}

		//Background
		if deathbat.Traits.Background != "" && deathbat.Traits.Background == check.Traits.Background {
			checkScore += weight["Background"]
		}

		//update if new twin
		if checkScore > score {
			twin = check
			score = checkScore
		}

		//if equal, pick closer
		if checkScore == score {
			if diff(deathbat.Id, check.Id) < diff(deathbat.Id, twin.Id) {
				twin = check
			}
		}
	}

	return twin, nil
}

//diff gets the absolute difference between two integers
func diff(x int, y int) (diff int) {
	diff = x - y
	if diff < 0 {
		return -diff
	}
	return diff
}
