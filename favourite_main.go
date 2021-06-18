package main

import (
	"database/sql"
  "fmt"

	"html/template"
	"log"
	"net/http"
  "io/ioutil"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
  
)

var DB *gorm.DB
var err error
var db *sql.DB
var updateId string

type list struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Link   string `json:"Link"`
	Favourite string `json:"Favourite"`
}

func Search(w http.ResponseWriter, r *http.Request) {
  response, err := http.Get("https://serpapi.com/playground?q=Apple&tbm=isch&ijn=0")
  if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }
}

func Index(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT * from lists")
	var lists []list
	for rows.Next() {
		var tem list
		rows.Scan(&tem.Id, &tem.Title, &tem.Link, &tem.Favourite)
		lists = append(lists, tem)
	}
	data := map[string]interface{}{
		"lists": lists,
	}
	temp, _ := template.ParseFiles("file:///C:/Users/ANIKET/Desktop/R/2ndexo.html")
	temp.Execute(w, data)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("file:///C:/Users/ANIKET/Desktop/R/2ndexo.html")
	temp.Execute(w, nil)
}

func ProcessInsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	link := r.Form.Get("link")
	DB.Exec("INSERT INTO lists(title, link, favourite) VALUES (?,?,?)", title, link, "pending")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	temp := params["id"]
	DB.Exec("DELETE FROM lists WHERE id=?", temp)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Favourite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	updateId = params["id"]
	DB.Exec("UPDATE lists SET favourite ='favourited' WHERE id=?", updateId)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/favourites")
	DB, err = gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/favourites"))
	if err != nil {
		log.Fatal(err)
	}

  r := mux.NewRouter()
  r.HandleFunc("/search", Search)
  r.HandleFunc("/changeFavourite/{id}", Favourite)
  r.HandleFunc("/", Index)
  r.HandleFunc("/insert", Insert)
  r.HandleFunc("/process", ProcessInsert)
  r.HandleFunc("/delete/{id}", Delete)
  log.Fatal(http.ListenAndServe(":9000", r))

}