package database

import (
	_ "database/sql"
	"io/ioutil"
	"log"
	"oauth2-gmail/model"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitTable() bool {
	f, err := os.Open("./db.sqlite.sql")
	defer f.Close()
	if err != nil {
		return false
	}
	lines, err := ioutil.ReadAll(f)
	if err!= nil {
		return false
	}
	db.Exec(string(lines[:]))
	return true
}

func GetUsers() []model.User {
	var users []model.User
	rows, err := db.Query(model.GetUsersQuery)
	user := model.User{}

	if err != nil {
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&user.Id,
			&user.DisplayName,
			&user.Mail,
			&user.JobTitle,
			&user.UserPrincipalName,
			&user.AccessToken,
			&user.AccessTokenActive,
			&user.RefreshToken,
			&user.MailIds)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}

func InsertUser(user model.User) {

	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.InsertUserQuery)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(user.Id,
		user.DisplayName,
		user.Mail,
		user.JobTitle,
		user.UserPrincipalName,
		user.AccessToken,
		user.AccessTokenActive,
		user.RefreshToken,
		user.MailIds)
	tx.Commit()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}

func UpdateUserMailIds(user model.User)  {
	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.UpdateMailIds)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(user.MailIds, user.Mail)
	tx.Commit()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}


func UpdateUserTokens(user model.User) {

	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.UpdateTokensQuery)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(user.AccessToken, user.RefreshToken, user.AccessTokenActive, user.Mail)
	tx.Commit()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

}

func GetUser(Mail string) model.User {

	var user model.User

	rows, err := db.Query(model.GetUserByEmailQuery, Mail)

	if err != nil {
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&user.Id,
			&user.DisplayName,
			&user.Mail,
			&user.JobTitle,
			&user.UserPrincipalName,
			&user.AccessToken,
			&user.AccessTokenActive,
			&user.RefreshToken,
			&user.MailIds)
		if err != nil {
			log.Fatal(err)
		}

	}

	return user
}
