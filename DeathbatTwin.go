package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

//tmpl represents the frontend html template
var tmpl *template.Template

//main handles the high level function calls for now
func main() {
	//load deathbat data
	filename := "deathbats.json"
	if err := loadDeathbats(filename); err != nil {
		log.Printf("err: main: %s", err)
		return
	}

	//prep backend and frontend goroutines to be linked
	wg := new(sync.WaitGroup)
	wg.Add(2)

	//backend
	http.HandleFunc("/twin", twin)
	go func() {
		log.Fatal(http.ListenAndServe(":6660", nil))
		wg.Done()
	}()

	//frontend
	mux := http.NewServeMux()
	tmpl = template.Must(template.ParseFiles("templates/index.gohtml"))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", home)

	go func() {
		log.Fatal(http.ListenAndServe(":6661", mux))
		wg.Done()
	}()

	wg.Wait()
}

//home is the frontend function handler
func home(w http.ResponseWriter, r *http.Request) {
	_ = tmpl.Execute(w, r)
}

//twin in the backend function handler
func twin(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	tokenIdString := r.FormValue("token_id")
	tokenId, err := strconv.Atoi(tokenIdString)
	if err != nil {
		log.Printf("twin: %s, %s", tokenIdString, ErrInvalidTokenId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if tokenId < 1 || tokenId > 10000 {
		log.Printf("twin: %d, %s", tokenId, ErrInvalidTokenId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sourceDeathbat, err := getDeathbat(tokenId)
	if err != nil {
		log.Printf("twin: %d, %s", tokenId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = sourceDeathbat.loadOwner(); err != nil {
		log.Printf("twin: %d, %s", tokenId, err)
	}

	twinDeathbat, err := sourceDeathbat.findTwin()
	if err != nil {
		log.Printf("twin: %d, %s", tokenId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = twinDeathbat.loadOwner(); err != nil {
		log.Printf("twin: %d, %s", tokenId, err)
	}

	data := struct {
		Source Deathbat
		Twin   Deathbat
	}{
		Source: sourceDeathbat,
		Twin:   twinDeathbat,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("twin: %d, %s", tokenId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(jsonData); err != nil {
		log.Printf("twin: %d, %s", tokenId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		return fmt.Errorf("loadDeathbats: Close: %w", err)
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
	for _, deathbat = range Deathbats {
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
		return fmt.Errorf("loadOwner: Get: %w", ErrOpenSeaUnresponsive)
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
func (deathbat *Deathbat) findTwin() (twin Deathbat, err error) {
	twin = *deathbat

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
	if deathbat.Traits.BrooksWackerman != "" || deathbat.Traits.JohnnyChrist != "" || deathbat.Traits.Shadows != "" ||
		deathbat.Traits.SynysterGates != "" || deathbat.Traits.ZackyVengeance != "" {
		return twin, nil
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
