package main

import (
  "fmt"
  "log"
  "strconv"
  "strings"
  "math/rand"
  "os"
  // "html/template"
  // "io/ioutil"
  "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Pick your dice!", r.URL.Path[1:])
}

func commonSetHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "This is the common set of dice!", r.URL.Path[1:])
}

func customSetHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Make your own set of dice!")
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
  // If it is a number simulate a dice with that many sides
  if _, err := strconv.Atoi(r.URL.Path[1:]); err == nil {
    fmt.Fprintf(w, "You rolled a D%s! \n", r.URL.Path[1:])
    inputInt, err := strconv.Atoi(r.URL.Path[1:])
    if err != nil {
        // handle error
        fmt.Println(err)
        os.Exit(2)
    }
    fmt.Fprintf(w, "it was a %d", randomValue(1, inputInt))
  } else if strings.HasSuffix(r.URL.Path[1:], ".html") {
    if strings.HasPrefix(r.URL.Path[1:], "static/") {
      http.ServeFile(w, r, r.URL.Path[1:])
    } else {
      http.ServeFile(w, r, "static/" + r.URL.Path[1:])
    }
  } else {
    fmt.Fprintf(w, "Sorry but it seems this page does not exist...")
  }
}

// Gets a random value from the low to high values. This will include the low and high values.
func randomValue(low int, high int) int {
  var scaledInt int = high - low + 1  // The +1 is to offset the values so it can be the high value.
  return rand.Intn(scaledInt) + low
}

// Basic Dice object.
// D20 would have a High of 20 and Low of 1.
type Dice struct {
    High int
    Low  int
}

func main() {
  // http.HandleFunc(nil, homeHandler)
  http.HandleFunc("/common/", commonSetHandler)
  http.HandleFunc("/custom/", customSetHandler)
  http.HandleFunc("/", randomPageHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
