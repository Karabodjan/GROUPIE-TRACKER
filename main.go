package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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
	Image string
	Name  string
	ID    int
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

func FindID(name string) int {
	apiURL_1 := "https://groupietrackers.herokuapp.com/api/artists"
	id := 0
	var Re []Artist
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
	for _, Relation := range Re {
		if Relation.Name == name {
			id = Relation.ID
		}
	}
	return id
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
		tmpl, _ := template.ParseFiles("index.html")
		err := tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("nope")
		}
	}
	http.SetCookie(w, cookie)
	//http.Redirect(w, r, "/result", http.StatusSeeOther)
}
func ValidQuery(input int) bool {
	fmt.Println(input)
	if input <= 0 || input >= 52 {
		return false
	} else {
		return true
	}
}

func Result(w http.ResponseWriter, r *http.Request) {
	Id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if ValidQuery(Id) == false {
		Id = FindID(r.URL.Query().Get("id"))
	}
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

	data := Data{
		Input: "ACDC",
		Name:  artist.Name,
		T1:    artist.Image,
		T2:    artist.Members,
		T3:    artist.FirstAlbum,
		T4:    artist.Creation,
		//T8:    data_arr2,
		//T5: location.Locations,
		//T6: dates.Dates,
		T7: relation.DatesLocations,
	}

	tmpl, err := template.ParseFiles("index.html")
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

	var mappa1 = []string{}
	var mappa2 = []string{}
	var mappa3 = []int{}
	for i := 1; i <= 52; i++ {
		artistChannel := make(chan Artist)
		go fetchArtist(i, artistChannel)
		artist := <-artistChannel
		mappa1 = append(mappa1, artist.Image)
		mappa2 = append(mappa2, artist.Name)
		mappa3 = append(mappa3, artist.ID)
	}
	//artistChannel := make(chan Artist)
	//Id := FindID("ACDC")
	//go fetchArtist(Id, artistChannel)

	data_arr := []Data2{}
	//artist := <-artistChannel
	for i := 0; i < len(mappa1); i++ {
		data_arr = append(data_arr, Data2{
			Image: mappa1[i],
			Name:  mappa2[i],
			ID:    mappa3[i],
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

	tmpl, err := template.ParseFiles("Api.html")
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

	tmpl, err := template.ParseFiles("result.html")
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
