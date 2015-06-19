// WebAppOne project main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/chrisc40/go-examples/webproject/myLib"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type FormUsers struct {
	Title string
	Users []myLib.User
}

type ValidationError struct {
	FieldName    string
	ErrorMessage string
}

type PostResponse struct {
	Status       string
	ErrorMessage string
	ErrorsList   []ValidationError
}

var tmpl map[string]*template.Template
var db *sql.DB

func validation(r *http.Request) []ValidationError {

	errors := make([]ValidationError, 0, 10)

	fname := r.FormValue("firstName")
	lname := r.FormValue("lastName")
	//uname := r.FormValue("username")
	//pwd := r.FormValue("password")

	if fname == "" {
		errors = append(errors, ValidationError{FieldName: "firstName", ErrorMessage: "FirstName must not be null"})
	}

	if lname == "" {
		errors = append(errors, ValidationError{FieldName: "lastName", ErrorMessage: "LastName must not be null"})
	}

	return errors
}

func userCreateHandleFunc(w http.ResponseWriter, r *http.Request) {
	//validation
	var res = PostResponse{}

	errors := validation(r)

	if len(errors) == 0 {
		res.Status = "OK"

		//validate transaction
		myLib.InsertUser(db, myLib.User{FirstName: r.FormValue("firstName"), LastName: r.FormValue("lastName"), Username: r.FormValue("username"), Password: r.FormValue("password")})

	} else {
		res.Status = "FAIL"
		res.ErrorMessage = "Error when creating user"
		res.ErrorsList = make([]ValidationError, 0, 10)
		res.ErrorsList = errors
	}

	bres, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("res=", string(bres))

	w.Header().Set("Content-Type", "application/json")
	w.Write(bres)
}

func userListHandleFunc(w http.ResponseWriter, r *http.Request) {

	Users := myLib.GetUsers(db)

	err := tmpl["users"].ExecuteTemplate(w, "base", FormUsers{Title: "Users", Users: Users})
	if err != nil {
		log.Fatal(err)
	}

}

func userEditHandleFunc(w http.ResponseWriter, r *http.Request) {

	sid := r.URL.Query().Get("id")

	if sid == "" {
		log.Fatal("User doesn't exist")
	}

	id, err := strconv.Atoi(sid)

	if err != nil {
		log.Fatal(err)
	}

	user := myLib.GetUser(db, id)

	err = tmpl["edit"].ExecuteTemplate(w, "base", user)

	if err != nil {
		log.Fatal(err)
	}
}

func updateUserHandleFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("enters updateUserHandleFunc")

	log.Println("id=", r.FormValue("id"))

	id, err := strconv.Atoi(r.FormValue("id"))

	if err != nil {
		log.Println("error when converting id to int")
		log.Fatal(err)
	}

	var res = PostResponse{}
	errors := validation(r)

	if len(errors) == 0 {
		res.Status = "OK"
		fmt.Println("update user with id=", id)
		//validate transaction
		myLib.UpdateUser(db, myLib.User{Id: id, FirstName: r.FormValue("firstName"), LastName: r.FormValue("lastName"), Username: r.FormValue("username"), Password: r.FormValue("password")})
		//http.Redirect(w, r, "/users", http.StatusFound)

	} else {
		res.Status = "FAIL"
		res.ErrorMessage = "Error when updating user"
		res.ErrorsList = make([]ValidationError, 0, 10)
		res.ErrorsList = errors
	}

	bres, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("res=", string(bres))

	w.Header().Set("Content-Type", "application/json")
	w.Write(bres)
}

func main() {
	tmpl = make(map[string]*template.Template)
	tmpl["users"] = template.Must(template.ParseFiles("templates/users.html", "templates/layout.html", "templates/header.html", "templates/footer.html"))
	tmpl["edit"] = template.Must(template.ParseFiles("templates/edit.html", "templates/layout.html", "templates/header.html", "templates/footer.html"))

	db = myLib.ConnectDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/users", userListHandleFunc)
	http.HandleFunc("/createUser", userCreateHandleFunc)
	http.HandleFunc("/editUser", userEditHandleFunc)
	http.HandleFunc("/updateUser", updateUserHandleFunc)
	//myLib.CloseDB()
	http.ListenAndServe(":8080", nil)
}
