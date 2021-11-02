package database

import (
	_ "database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"oauth2-gmail/model"
)

func GetEmailsByUser(email string) []model.Mail{
	var mails []model.Mail
	rows, err := db.Query(model.GetUserMailsQuery, email)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&mail.Id,
			&mail.User,
			&mail.Subject,
			&mail.SenderEmail,
			&mail.SenderName,
			&temp,
			&mail.BodyPreview,
			&mail.BodyType,
			&mail.BodyContent,
			&mail.Date)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(temp), &mail.Attachments)
		mails = append(mails,mail)
	}

	return mails
}

func GetPageEmails() []model.Mail{
	var mails []model.Mail
	rows, err := db.Query(model.GetMailsQuery)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&mail.Id,
			&mail.User,
			&mail.Subject,
			&mail.SenderEmail,
			&mail.SenderName,
			&temp,
			&mail.BodyPreview,
			&mail.BodyType,
			&mail.BodyContent,
			&mail.Date)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(temp), &mail.Attachments)
		mails = append(mails,mail)
	}

	return mails


}

func InsertEmail(mail model.Mail){
	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.InsertMailQuery)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	marshal, _ := json.Marshal(mail.Attachments)
	_, err := stmt.Exec(mail.Id,
		mail.User,
		mail.Subject,
		mail.SenderEmail,
		mail.SenderName,
		string(marshal),
		mail.BodyPreview,
		mail.BodyType,
		mail.BodyContent,
		mail.ToRecipient,
		mail.ToRecipientName,
		mail.Date.Format("2006-01-02 15:04:05"))
	tx.Commit()
	if err != nil{
		log.Printf("ERROR: %s",err)
	}

}

func SearchUserEmails(email string,searchKey string) []model.Mail {
	var mails []model.Mail

	searchKey = "%" + searchKey + "%"

	rows, err := db.Query(model.SearchUserMailsQuery,email,searchKey)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&mail.Id,
			&mail.User,
			&mail.Subject,
			&mail.SenderEmail,
			&mail.SenderName,
			&temp,
			&mail.BodyPreview,
			&mail.BodyType,
			&mail.BodyContent,
			&mail.ToRecipient,
			&mail.ToRecipientName,
			&mail.Date)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(temp), &mail.Attachments)
		mails = append(mails,mail)
	}

	return mails
}


func SearchEmails(searchKey string) []model.Mail {
	var mails []model.Mail

	searchKey = "%" + searchKey + "%"

	rows, err := db.Query(model.SearchEmailQuery,searchKey)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&mail.Id,
			&mail.User,
			&mail.Subject,
			&mail.SenderEmail,
			&mail.SenderName,
			&temp,
			&mail.BodyPreview,
			&mail.BodyType,
			&mail.BodyContent,
			&mail.ToRecipient,
			&mail.ToRecipientName,
			&mail.Date)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal([]byte(temp), &mail.Attachments)
		mails = append(mails,mail)
	}

	return mails
}


func GetEmail(id string) model.Mail {
	var temp string
	row := db.QueryRow(model.GetEmailQuery,id)
	mail := model.Mail{}
	err := row.Scan(&mail.Id,
		&mail.User,
		&mail.Subject,
		&mail.SenderEmail,
		&mail.SenderName,
		&temp,
		&mail.BodyPreview,
		&mail.BodyType,
		&mail.BodyContent,
		&mail.ToRecipient,
		&mail.ToRecipientName,
		&mail.Date)
	if err != nil {
		return mail
	}
	json.Unmarshal([]byte(temp), &mail.Attachments)
	return mail
}
