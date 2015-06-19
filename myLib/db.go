package myLib

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func ConnectDB() *sql.DB {

	db, err := sql.Open("mysql",
		"root:Optus123@tcp(127.0.0.1:3306)/MyDatabase")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetUsers(db *sql.DB) []User {
	log.Println("enters getusers")
	var users []User
	rows, err := db.Query("select id, firstName, lastName,username,password from user")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var r User
		err := rows.Scan(&r.Id, &r.FirstName, &r.LastName, &r.Username, &r.Password)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(r.FirstName, r.LastName)

		users = append(users, r)
	}

	return users
}

func getLastId(db *sql.DB) int {
	log.Println("enters getusers")
	var pk int
	err := db.QueryRow("select max(id) from user").Scan(&pk)

	if err != nil {
		log.Fatal(err)
	}

	return pk

}

func GetUser(db *sql.DB, id int) User {
	var r User

	stmt, err := db.Prepare("select id, firstName, lastName, username, password from user where id=?")

	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(id).Scan(&r.Id, &r.FirstName, &r.LastName, &r.Username, &r.Password)

	if err != nil {
		log.Fatal(err)
	}

	return r

}

func InsertUser(db *sql.DB, user User) {
	log.Println("enters insert user")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := db.Prepare("insert into user values (?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}
	log.Println("error when saving user")
	pk := getLastId(db) + 1
	_, err = stmt.Exec(pk, user.FirstName, user.LastName, user.Username, user.Password)

	if err != nil {
		log.Println("error when saving user")
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	stmt.Close()
}

func UpdateUser(db *sql.DB, user User) {
	log.Println("enters update user")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := db.Prepare("update user set firstName=? , lastName=?,username=?,password=? where id=?")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("error when updating user")

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Username, user.Password, user.Id)

	if err != nil {
		log.Println("error when updating user")
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	stmt.Close()
}

func CloseDB(db *sql.DB) {
	defer db.Close()
}
