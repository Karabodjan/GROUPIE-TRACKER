package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var (
	input string
)

type Data struct {
	Input string
	Name  string
	T1    string
	T2    []string
	T3    string
	T4    int
	T5    []string
	T6    []string
	T8    []Data1
	T7    map[string][]string
}
type Data2 struct {
	Image    string
	Name     string
	ID       int
	FirstAlb string
	Creation int
}
type Data1 struct {
	T5 string
	T6 string
}

type Cookie struct {
	Name  string
	Value int
}

type Artist struct {
	ID          int      `json:"id"`
	Image       string   `json:"image"`
	Name        string   `json:"name"`
	Members     []string `json:"members"`
	Creation    int      `json:"creationDate"`
	FirstAlbum  string   `json:"firstAlbum"`
	Location    string   `json:"locations"`
	ConcertDate string   `json:"concertDates"`
	Relation    string   `json:"relations"`
}
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}
type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type APIResponse struct {
	Data interface{} `json:"data"`
}

func fetchArtist(Id int, ch chan Artist) {
	url := "https://groupietrackers.herokuapp.com/api/artists/" + strconv.Itoa(Id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des informations de l'artiste :", err)
		ch <- Artist{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse HTTP :", err)
		ch <- Artist{}
		return
	}

	var artist Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation des données JSON de l'artiste :", err)
		ch <- Artist{}
		return
	}

	ch <- artist
}
func fetchLocation(Id int, ch chan Location) {
	url := "https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(Id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des informations de l'artiste :", err)
		ch <- Location{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse HTTP :", err)
		ch <- Location{}
		return
	}

	var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation des données JSON de l'artiste :", err)
		ch <- Location{}
		return
	}

	ch <- location
}
func fetchDates(Id int, ch chan Dates) {
	url := "https://groupietrackers.herokuapp.com/api/dates/" + strconv.Itoa(Id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des informations de l'artiste :", err)
		ch <- Dates{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse HTTP :", err)
		ch <- Dates{}
		return
	}

	var dates Dates
	err = json.Unmarshal(body, &dates)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation des données JSON de l'artiste 3 :", err)
		ch <- Dates{}
		return
	}

	ch <- dates
}
func fetchRelation(Id int, ch chan Relation) {
	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(Id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des informations de l'artiste :", err)
		ch <- Relation{}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse HTTP :", err)
		ch <- Relation{}
		return
	}

	var relation Relation
	err = json.Unmarshal(body, &relation)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation des données JSON de l'artiste 4 :", err)
		ch <- Relation{}
		return
	}

	ch <- relation
}

func Index(w http.ResponseWriter, r *http.Request) {
	data := Data{}
	http.Redirect(w, r, "/result", http.StatusSeeOther)
	input := r.FormValue("input")
	cookie := &http.Cookie{
		Name:  "input",
		Value: input,
	}
	if r.Method == "POST" {
		tmpl, _ := template.ParseFiles("./template/index.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("nope")
		}
	}
	http.SetCookie(w, cookie)
	//http.Redirect(w, r, "/result", http.StatusSeeOther)
}
func ValidQuery(input int) bool {

	if input <= 0 || input >= 52 {
		return false
	} else {
		return true
	}
}
func isInContainer(container []string, target string) bool {
	for _, str := range container {
		cleanedStr := strings.ReplaceAll(str, "+", " ")
		if cleanedStr == target {
			return true
		}
	}
	return false
}
func FindID(input string) int {
	fmt.Println(input)
	input1, _ := strconv.Atoi(input) // str to int
	//input2 := strconv.Itoa(input)    // str to int
	apiURL_1 := "https://groupietrackers.herokuapp.com/api/artists"
	apiURL_2 := "https://groupietrackers.herokuapp.com/api/location"
	id := 0
	var Re []Artist
	var Ri []Location
	response, err := http.Get(apiURL_1)
	if err != nil {
		fmt.Println("error 1")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error 2")
	}
	err = json.Unmarshal(body, &Re)
	if err != nil {
		fmt.Println("error 3")
	}
	response2, err := http.Get(apiURL_2)
	if err != nil {
		fmt.Println("error 1")
	}
	defer response2.Body.Close()
	body2, err := ioutil.ReadAll(response2.Body)
	if err != nil {
		fmt.Println("error 2")
	}
	err = json.Unmarshal(body2, &Ri)
	if err != nil {
		fmt.Println("error 3")
	}

	for _, Relation := range Re {
		if Relation.Name == input {
			fmt.Println(input, "rzrzrz")
			id = Relation.ID
		} else if Relation.Creation == input1 {
			id = Relation.ID
		} else if isInContainer(Relation.Members, input) == true {
			id = Relation.ID
		} else if Relation.ID == input1 {
			id = Relation.ID
		} else if Relation.FirstAlbum == input {
			id = Relation.ID
		} else {
			for _, Location := range Ri {
				if isInContainer(Location.Locations, input) == true {
					id = Location.ID
				}
			}
		}
	}
	return id
}

func Result(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	//id2,_:=strconv.Atoi(id)
	//Id:=0
	//if ValidQuery(id2){
	//	Id = strconv.Itoa(id)
	//}
	Id := FindID(id)
	fmt.Println(id, "gdggd")

	artistChannel := make(chan Artist)
	locationChannel := make(chan Location)
	datesChannel := make(chan Dates)
	relationChannel := make(chan Relation)

	go fetchArtist(Id, artistChannel)
	go fetchLocation(Id, locationChannel)
	go fetchDates(Id, datesChannel)
	go fetchRelation(Id, relationChannel)

	artist := <-artistChannel
	relation := <-relationChannel
	location := <-locationChannel

	data := Data{
		Input: "ACDC",
		Name:  artist.Name,
		T1:    artist.Image,
		T2:    artist.Members,
		T3:    artist.FirstAlbum,
		T4:    artist.Creation,
		T5:    location.Locations,
		T7:    relation.DatesLocations,
	}

	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		http.Error(w, "Erreur de rendu du template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur de rendu du template", http.StatusInternalServerError)
		return
	}
}

func Api(w http.ResponseWriter, r *http.Request) {
	MinCreation, _ := strconv.Atoi(r.FormValue("MinCreation"))
	MaxCreation, _ := strconv.Atoi(r.FormValue("MaxCreation"))
	//AlbulmCreation, _ := strconv.Atoi(r.FormValue("AlbulmCreation"))
	//MinYear, _ := strconv.Atoi(r.FormValue("MinYear"))
	//Country := r.FormValue("Country")

	var vec1 = []string{}
	var vec2 = []string{}
	var vec3 = []int{}
	var vec4 = []int{}
	var vec5 = []string{}

	for i := 1; i <= 52; i++ {
		artistChannel := make(chan Artist)
		go fetchArtist(i, artistChannel)
		artist := <-artistChannel
		if MinCreation != 0 && MaxCreation != 0 {
			if artist.Creation < MinCreation || artist.Creation > MaxCreation {
				vec1 = append(vec1, artist.Image)
				vec2 = append(vec2, artist.Name)
				vec3 = append(vec3, artist.ID)
				vec4 = append(vec4, artist.Creation)
				vec5 = append(vec5, artist.FirstAlbum)
			}
		} else {
			vec1 = append(vec1, artist.Image)
			vec2 = append(vec2, artist.Name)
			vec3 = append(vec3, artist.ID)
			vec4 = append(vec4, artist.Creation)
			vec5 = append(vec5, artist.FirstAlbum)
		}
	}

	data_arr := []Data2{}
	//artist := <-artistChannel
	for i := 0; i < len(vec1); i++ {
		data_arr = append(data_arr, Data2{
			Image:    vec1[i],
			Name:     vec2[i],
			ID:       vec3[i],
			FirstAlb: vec5[i],
			Creation: vec4[i],
		})
	}
	if r.Method == "POST" {
		Query := r.FormValue("Query")
		http.SetCookie(w, &http.Cookie{
			Name:  "Query",
			Value: Query,
		})
		http.Redirect(w, r, "/result", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("./template/Api.html")
	if err != nil {
		http.Error(w, "Erreur de rendu du template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data_arr)
	if err != nil {
		http.Error(w, "Erreur de rendu du templateee", http.StatusInternalServerError)
		return
	}
}

func Map(w http.ResponseWriter, r *http.Request) {
	url := "https://google-maps-geocoding.p.rapidapi.com/geocode/json?address=164%20Townsend%20St.%2C%20San%20Francisco%2C%20CA&language=en"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "0c396aa89emsh0b8b177e474fb34p1b1734jsn177f3aa11bd6")
	req.Header.Add("X-RapidAPI-Host", "google-maps-geocoding.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	data := Data{}

	tmpl, err := template.ParseFiles("./template/result.html")
	if err != nil {
		http.Error(w, "Erreur de rendu du template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur de rendu du templateee", http.StatusInternalServerError)
		return
	}

}
func main() {

	/////////////////////////////// SERVEUR /////////////////
	http.HandleFunc("/", Api)
	http.HandleFunc("/result", Result)
	http.HandleFunc("/Map", Map)
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe("localhost:8080", nil)
}
