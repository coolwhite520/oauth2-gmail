package database

import (
	"log"
	"oauth2-gmail/model"
)

func SearchAttachmentByMailId(mailId string) []model.Attachment {
	var attachments []model.Attachment
	rows, err := db.Query("select * from attachment where mailId=?", mailId)
	attachment := model.Attachment{}
	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&attachment.MailId, &attachment.FileName, &attachment.FilePath, &attachment.AttachmentId)
		if err != nil {
			log.Fatal(err)
		}
		attachments = append(attachments, attachment)
	}
	return attachments
}

func InsertAttachment(attachment model.Attachment)  {
	tx, _ := db.Begin()
	stmt, err := tx.Prepare(model.InsertAttachmentQuery)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(attachment.MailId, attachment.FileName, attachment.FilePath, attachment.AttachmentId)
	tx.Commit()
	if err != nil{
		log.Printf("ERROR: %s",err)
	}
}
