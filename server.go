package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Users struct {
	Group   int
	Surname string
	Name    string
}
type Objects struct {
	First  string
	Second string
	Third  string
	Fourth string
}

var Lists = []Users{}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.ExecuteTemplate(w, "login", nil)
}

func survey(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/survey.html")
	tmpl.ExecuteTemplate(w, "survey", nil)
}

func (u Users) Db_insert(w http.ResponseWriter, r *http.Request) {
	u.Group, _ = strconv.Atoi(r.FormValue("group"))
	u.Surname = r.FormValue("surname")
	u.Name = r.FormValue("name")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `students`(`Groupp`, `Surname`, `Name`, `Object`) VALUES ('%d','%s','%s','%s')", u.Group, u.Surname, u.Name, "Базы данных"))
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func authorize(w http.ResponseWriter, r *http.Request) {
	Users.Db_insert(Users{}, w, r)
	http.Redirect(w, r, "/survey", http.StatusSeeOther)
}

func (obj *Objects) savefirstround(w http.ResponseWriter, r *http.Request) {
	obj.First = r.FormValue("firstpair")
	obj.Second = r.FormValue("secondpair")
	obj.Third = r.FormValue("thirdpair")
	obj.Fourth = r.FormValue("fourthpair")
	http.Redirect(w, r, fmt.Sprintf("/nextround?first=%s&second=%s&third=%s&fourth=%s", obj.First, obj.Second, obj.Third, obj.Fourth), http.StatusSeeOther)
}

func (obj *Objects) savesecondround(w http.ResponseWriter, r *http.Request) {
	obj.First = r.FormValue("firstpair")
	obj.Second = r.FormValue("secondpair")
	http.Redirect(w, r, fmt.Sprintf("/final?first=%s&second=%s&", obj.First, obj.Second), http.StatusSeeOther)
}

func final(w http.ResponseWriter, r *http.Request) {
	data := Objects{
		First:  r.URL.Query().Get("first"),
		Second: r.URL.Query().Get("second"),
	}
	tmpl, _ := template.ParseFiles("templates/final.html")
	tmpl.ExecuteTemplate(w, "final", data)
}

func nextround(w http.ResponseWriter, r *http.Request) {
	data := Objects{
		First:  r.URL.Query().Get("first"),
		Second: r.URL.Query().Get("second"),
		Third:  r.URL.Query().Get("third"),
		Fourth: r.URL.Query().Get("fourth"),
	}
	tmpl, _ := template.ParseFiles("templates/nextround.html")
	tmpl.ExecuteTemplate(w, "nextround", data)
}

func (obj *Objects) savefinalround(w http.ResponseWriter, r *http.Request) {
	obj.First = r.FormValue("firstpair")
	http.Redirect(w, r, fmt.Sprintf("/results?first=%s&", obj.First), http.StatusSeeOther)
}

func (u Users) results(w http.ResponseWriter, r *http.Request) {
	data := Objects{
		First: r.URL.Query().Get("first"),
	}
	tmpl, _ := template.ParseFiles("templates/results.html")
	tmpl.ExecuteTemplate(w, "results", data)

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var lastID int
	err = db.QueryRow("SELECT MAX(newid) FROM students").Scan(&lastID)
	if err != nil {
		panic(err.Error())
	}

	update, err := db.Exec("UPDATE students SET Object=? WHERE newid=?", data.First, lastID)
	if err != nil {
		panic(err.Error())
	}
	rowsAffected, _ := update.RowsAffected()
	fmt.Printf("Затронуто %d строка.\n", rowsAffected)
}

func handleRequest() {
	var obj = &Objects{}
	var u = &Users{}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", login).Methods("GET")
	rtr.HandleFunc("/authorize", authorize).Methods("POST")
	rtr.HandleFunc("/survey", survey).Methods("GET", "POST")
	rtr.HandleFunc("/savefirstround", obj.savefirstround).Methods("POST")
	rtr.HandleFunc("/nextround", nextround).Methods("GET", "POST")
	rtr.HandleFunc("/savesecondround", obj.savesecondround).Methods("POST")
	rtr.HandleFunc("/final", final).Methods("GET", "POST")
	rtr.HandleFunc("/savefinalround", obj.savefinalround).Methods("POST")
	rtr.HandleFunc("/results", u.results).Methods("GET", "POST")
	http.Handle("/", rtr)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}
