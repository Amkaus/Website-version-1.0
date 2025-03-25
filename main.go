package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                      uint16
	Title, Anons, Full_text string
}

var posts = []Article{}
var showPOST = Article{}

func contact(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/contact.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "contact", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "index", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не вся информация была заполнена")
	} else {

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang") // connection and error
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//установка данных

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`,`anons`,`full_text`) VALUES('%s','%s','%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/articles", http.StatusSeeOther) //1
	}
}

func articles(w http.ResponseWriter, r *http.Request) { //1
	t, err := template.ParseFiles("templates/art.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang") // connection and error
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "art", posts)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/footer.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang") // connection and error
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPOST = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_text)
		if err != nil {
			panic(err)
		}

		showPOST = post
	}

	t.ExecuteTemplate(w, "show", showPOST)
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/articles", articles).Methods("GET") // 1
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/contacts/", contact).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
