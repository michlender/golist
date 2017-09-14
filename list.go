package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
  "regexp"
  "errors"
)

// regexp used to validate user input
var validPath = regexp.MustCompile("^/(add|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
  Title string
  Body []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t, err := template.ParseFiles(tmpl + ".html")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  err = t.Execute(w, p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// Validates user input of the title
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("Invalid Page Title")
  }
  return m[2], nil // The title is the second subexpression.
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title, err := getTitle(w, r)
  if err != nil {
    return
    }
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/add/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
  title, err := getTitle(w, r)
  if err != nil {
    return
  }
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "add", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title , err := getTitle(w, r)
  if err != nil {
    return
  }
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err = p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/add/", addHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}
