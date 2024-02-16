package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var (
	input string
)

type Data struct {
	Input string
	T1    string
	T2    string
	T3    string
	T4    string
	T6    string
	T7    string
}

type Cookie struct {
	Name  string
	Value int
}

type Relation struct {
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

const form = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href ="/static/index.css" rel="stylesheet">
    <title>Barre de Recherche</title>

</head>
<body>
<div class="search-container">
<form action="" method="post" class="form-example">
  <div class="form-example">
    <input type="text" name="input" id="input" required />
  <div class="form-example">
    <input type="submit" value="Subscribe!" />
  </div>
</form>
</div>
</body>
</html>
`
const form2 = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Page avec des Blocs de Texte</title>
    <link href ="/static/result.css" rel="stylesheet">
</head>
<body>
<div class="container">
    <div class="column">
        <div class="block">
            <p>{{.T1}}</p>
        </div>
        <div class="block">
            <p>Texte du Bloc 2</p>
        </div>
        <div class="block">
            <p>Texte du Bloc 3</p>
        </div>
    </div>
    <div class="column">
        <div class="block">
            <p>Texte du Bloc 4</p>
        </div>
        <div class="block">
            <p>Texte du Bloc 5</p>
        </div>
        <div class="block">
            <p>Texte du Bloc 6</p>
        </div>
    </div>
</div>
</body>
</html>
`

// https://groupietrackers.herokuapp.com/api/dates

func Index(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("input")
	cookie := http.Cookie{
		Name:  "input",
		Value: input,
	}
	if r.Method == "GET" {
		tmpl := template.Must(template.New("index").Parse(form))
		tmpl.Execute(w, nil)
		return
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/result", http.StatusSeeOther)
}
func Result(w http.ResponseWriter, r *http.Request) {
	apiUrl := "https://groupietrackers.herokuapp.com/api/artists"

	var Re []Relation

	response, err := http.Get(apiUrl)
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
	//cookie, err := r.Cookie("input")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := Data{
		T2: Re[1].Name,
		T3: Re[1].ConcertDate,
	}
	//word := cookie.Value
	tmpl := template.Must(template.New("result").Parse(form2))
	tmpl.Execute(w, data)
}

func main() {

	///////////////////////////////// SERVEUR /////////////////
	http.HandleFunc("/", Index)
	http.HandleFunc("/result", Result)
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe("localhost:8080", nil)
}

//apiUrl := "https://groupietrackers.herokuapp.com/api/artists"
//
//var Re []Relation
//
//response, err := http.Get(apiUrl)
//if err != nil {
//	fmt.Println("error 1")
//}
//defer response.Body.Close()
//body, err := ioutil.ReadAll(response.Body)
//if err != nil {
//	fmt.Println("error 2")
//}
//err = json.Unmarshal(body, &Re)
//if err != nil {
//	fmt.Println("error 3")
//}
//for _, g := range Re {
//	fmt.Printf("ID: %d\n", g.ID)
//	fmt.Printf("Name: %s\n", g.Name)
//	fmt.Println("Members :")
//	for _, member := range g.Members {
//		fmt.Printf("- %s\n", member)
//	}
//	fmt.Printf("Creation Date: %d\n", g.Creation)
//	fmt.Printf("First Album: %s\n", g.FirstAlbum)
//	fmt.Printf("Locations: %s\n", g.Location)
//	fmt.Printf("Concert Dates: %s\n", g.ConcertDate)
//	fmt.Printf("Relations: %s\n", g.Relation)
//	fmt.Println()
//}

//	for _, entry := range Ya.IndexData {
//		//		fmt.Printf("ID: %d\n", entry.ID)
//		//		fmt.Println("Dates:")
//		//		for _, date := range entry.Dates {
//		//			fmt.Println(date)
//		//		}
//		//		fmt.Println()
//		//	}
//}

//func main() {
//	apiUrl := "https://groupietrackers.herokuapp.com/api/dates"
//
//	var Ya Data
//
//	response, err := http.Get(apiUrl)
//	if err != nil {
//		fmt.Printf("Erreur lors de la requête GET: %v\n", err)
//		return
//	}
//	defer response.Body.Close()
//
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Printf("Erreur lors de la lecture du corps de la réponse: %v\n", err)
//		return
//	}
//
//	err = json.Unmarshal(body, &Ya)
//	if err != nil {
//		fmt.Printf("Erreur lors de la conversion des données JSON: %v\n", err)
//		return
//	}
//
//	fmt.Println("Index:")
//	for _, entry := range Ya.IndexData {
//		fmt.Printf("ID: %d\n", entry.ID)
//		fmt.Println("Dates:")
//		for _, date := range entry.Dates {
//			fmt.Println(date)
//		}
//		fmt.Println()
//	}
//}
