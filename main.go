package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("templates/notFound.html", "templates/header.html", "templates/footer.html", "templates/index.html"))

//A Page structure
type Page struct {
	PageTitle string
}

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		PageTitle: "Home",
		ListTitle: "My TODO List",
		Todos: []Todo{
			{Title: "Project Setup", Done: true},
			{Title: "Setup CI", Done: false},
			{Title: "Setup CD", Done: true},
			{Title: "Make Test Cases", Done: false},
			{Title: "Basic Dice Rolls Based on URL", Done: true},
			{Title: "Flushout Navbar", Done: false},
			{Title: "Create 404", Done: true},
			{Title: "Create Common Dice Set Page", Done: false},
			{Title: "Create Custom Dice Set Page", Done: false},
			{Title: "Add Dice Graphics", Done: false},
			{Title: "Add How to Use Page", Done: false},
			{Title: "Flushout Dedicated Dice Page with Reroll button", Done: false},
			{Title: "Make a logo", Done: false},
		},
	}
	display(w, "index", data)
}

func commonSetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the common set of dice!", r.URL.Path[1:])
}

func customSetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Make your own set of dice!")
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
	// If empty show the home page
	// If a number simulate a dice with that many sides
	// If static page show the static package
	// Else show the 404 page
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else if _, err := strconv.Atoi(r.URL.Path[1:]); err == nil {
		fmt.Fprintf(w, "You rolled a D%s! \n", r.URL.Path[1:])
		inputInt, err := strconv.Atoi(r.URL.Path[1:])
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Fprintf(w, "it was a %d", randomValue(1, inputInt))
	} else if strings.HasSuffix(r.URL.Path[1:], ".html") {
		http.ServeFile(w, r, "static/html/"+r.URL.Path[1:])
	} else {
		fmt.Println("Sorry but it seems this page does not exist...")
		errorHandler(w, r, http.StatusNotFound)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		display(w, "404", &Page{PageTitle: "404"})
	} else {
		http.ServeFile(w, r, "static/html/issue.html")
	}
}

// Gets a random value from the low to high values. This will include the low and high values.
func randomValue(low int, high int) int {
	var scaledInt int = high - low + 1 // The +1 is to offset the values so it can be the high value.
	return rand.Intn(scaledInt) + low
}

func getPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return ":" + value
	}
	return ":8080"
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	ListTitle string
	Todos     []Todo
}

// Basic Dice object.
// D20 would have a High of 20 and Low of 1.
type Dice struct {
	High int
	Low  int
}

func main() {
	// http.HandleFunc("", homeHandler)
	http.HandleFunc("/common/", commonSetHandler)
	http.HandleFunc("/custom/", customSetHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", randomPageHandler)

	var port string = getPort()
	fmt.Println("Now listening to port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
