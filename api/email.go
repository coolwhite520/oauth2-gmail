package api

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"oauth2-gmail/database"
	"oauth2-gmail/model"
	"oauth2-gmail/xlsxtool"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// SendEmail will send an email using the api
func SendEmail(user model.User, email model.SendEmailStruct) (string, int) {

	data, _ := json.Marshal(email)
	resp, code := CallAPIMethod("POST", model.EmailEndPointRoot, "/me/sendMail", user.AccessToken, "", []byte(data), "application/json")
	return resp, code
	//log.Printf("E-mail to : %s responded with status code: %d", email.Message.ToRecipients[0].EmailAddress.Address, respcode)
}

func ViewEmailById(user model.User, emailID string) model.SingleMail {
	dbMail := database.GetEmail(emailID)
	mail := model.SingleMail{}
	mail.BodyPreview = dbMail.BodyPreview
	mail.Subject = dbMail.Subject
	mail.Body.Content = dbMail.BodyContent
	mail.Body.ContentType = dbMail.BodyType
	mail.ReceivedDateTime = dbMail.Date
	return mail
}
func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// 打包成zip文件
func Zip(src_dir string, zip_file_name string) {
	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)
	// 创建：zip文件
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()
	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	// 遍历路径信息
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {
		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}


func GenerateMailHtml(user model.User, mails []model.Mail) (string, error) {
	folderDir := fmt.Sprintf("./exports/%s/%s", user.UserPrincipalName, fmt.Sprintf("%s", time.Now().Format("2006_01_02_15_04_05")))
	folderHtmlDir := fmt.Sprintf("%s/html", folderDir)

	if _, err := os.Stat(folderHtmlDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(folderHtmlDir, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
	}

	for _, v := range mails {
		filename := fmt.Sprintf("%s/%s.html", folderHtmlDir,  v.Id)
		ioutil.WriteFile(filename, []byte(v.BodyContent), 0644)
		for _, a := range v.Attachments {
			mailAttachDir := fmt.Sprintf("%s/attachment/%s", folderDir, a.MailId)
			os.MkdirAll(mailAttachDir, os.ModePerm)
			dest := fmt.Sprintf("%s/%s", mailAttachDir, a.FileName)
			copyFile(a.FilePath, dest)
		}
	}
	page := model.Page{
		Title:             "",
		Email:             "",
		User:              model.User{},
		UserList:          nil,
		EmailList:         mails,
		FileList:          nil,
		SearchFiles:       model.Files{},
		Mail:              model.SingleMail{},
		Message:           "",
		Success:           false,
		URL:               "",
		Sum:               0,
		PageSize:          0,
		CurrentPageNumber: 0,
		PageList:          nil,
	}
	tpl, err := template.ParseFiles("templates/exportcatalogue.html")
	if err != nil {
		log.Fatal(err)
	}
	indexFileName := fmt.Sprintf("%s/index.html", folderDir)
	file, err := os.Create(indexFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	tpl.Execute(file,  page)
	// 压缩
	zipFile := fmt.Sprintf("%s.zip", folderDir)
	Zip(folderDir, zipFile)
	os.RemoveAll(folderDir)
	return zipFile, nil
}


func GenerateMailExcel(user model.User, mails []model.Mail) (bool, error) {
	t1 := time.Now()
	defer func() {
		fmt.Println(time.Since(t1))
	}()
	t := make([]string, 0)
	t = append(t, "邮件ID")
	t = append(t, "发送方姓名")
	t = append(t, "发送方MAIL")
	t = append(t, "接收方姓名")
	t = append(t, "接收方MAIL")
	t = append(t, "邮件内容")
	t = append(t, "附件")
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet")
	if err != nil {
		return false,err
	}
	titleRow := sheet.AddRow()
	xlsRow := xlsxtool.NewRow(titleRow, t)
	err = xlsRow.SetRowTitle()
	if err != nil {
		return false,err
	}
	for _, v := range mails {
		currentRow := sheet.AddRow()
		tmp := make([]string, 0)
		tmp = append(tmp, v.Id)
		tmp = append(tmp, v.SenderName)
		tmp = append(tmp, v.SenderEmail)
		tmp = append(tmp, v.ToRecipientName)
		tmp = append(tmp, v.ToRecipient)
		tmp = append(tmp, v.BodyContent)
		marshal, err := json.Marshal(v.Attachments)
		if err != nil {
			tmp = append(tmp, string(marshal))
		} else {
			tmp =  append(tmp, "")
		}

		xlsRow := xlsxtool.NewRow(currentRow, tmp)
		err = xlsRow.GenerateRow()
		if err != nil {
			return false, err
		}
	}

	folderDir := fmt.Sprintf("./exports/%s", user.UserPrincipalName)
	if _, err := os.Stat(folderDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(folderDir, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
	}
	filename := fmt.Sprintf("%s/%s.xlsx", folderDir,  time.Now().Format("2006_01_02_15_04_05"))
	err = file.Save(filename)
	if err != nil {
		return false, err
	}
	return true, nil
}




func GetEmailAllIds(user model.User) []string {
	var ids []string
	nextPageToken := ""
	for {
		values := url.Values{}
		messages := model.Messages{}
		values.Add("maxResults", "500")
		if nextPageToken != "" {
			values.Add("pageToken", nextPageToken)
		}
		messagesResponse, _ := CallAPIMethod("GET", model.EmailEndPointRoot, "/users/me/messages", user.AccessToken, values.Encode(), nil, "")
		json.Unmarshal([]byte(messagesResponse), &messages)
		for _, item := range messages.Value {
			ids = append(ids, item.Id)
		}
		if messages.NextPageToken != "" {
			nextPageToken = messages.NextPageToken
		} else {
			break
		}
	}
	return ids
}

func DownAttachmentFile(user model.User, mailId string, attachmentId string, filename string) (model.Attachment, error) {
	var a model.Attachment
	folderDir := fmt.Sprintf("./attachment/%s", mailId)
	if _, err := os.Stat(folderDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(folderDir, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println(err)
		}
	}
	downFile := fmt.Sprintf("%s/%s", folderDir, filepath.Base(filename))
	out, err := os.Create(downFile)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	interFacePath := fmt.Sprintf("/users/me/messages/%s/attachments/%s", mailId, attachmentId)
	attachmentContent, _ := CallAPIMethod("GET", model.EmailEndPointRoot, interFacePath, user.AccessToken, "", nil, "")
	var attachment map[string]string
	json.Unmarshal([]byte(attachmentContent), &attachment)
	decodeBytes, err := base64.URLEncoding.DecodeString(attachment["data"])
	if err != nil {
		return a, err
	}
	out.Write(decodeBytes)
	a = model.Attachment{
		MailId:       mailId,
		FileName:     filename,
		FilePath:     downFile,
		AttachmentId: attachmentId,
	}
	database.InsertAttachment(a)
	return a, nil
}

func FetchOneMail(user model.User, id string) model.Mail {
	interFacePath := fmt.Sprintf("/users/me/messages/%s", id)
	mail, _ := CallAPIMethod("GET", model.EmailEndPointRoot, interFacePath, user.AccessToken, "", nil, "")
	m := model.MessageMail{}
	json.Unmarshal([]byte(mail), &m)
	var subject = ""
	var SenderEmail = ""
	var SenderName = ""
	var ToRecipient = ""
	var ToRecipientName = ""
	var BodyType = ""
	var BodyContent = ""
	for _, obj := range m.Payload.Headers {
		if obj.Name == "Subject" {
			subject = obj.Value
		} else if obj.Name == "From" {
			arr := strings.Split(obj.Value, " ")
			arrLen := len(arr)
			for i := 0; i < arrLen-1; i++ {
				SenderName += " "
				SenderName += arr[i]
			}
			SenderEmail = arr[arrLen-1]
		} else if obj.Name == "To" {
			ToRecipient = obj.Value
			ToRecipientName = obj.Value
		}
	}
	var Attachments []model.Attachment
	if len(m.Payload.Parts) > 0 {
		for _, item := range m.Payload.Parts {
			if item.MimeType == "text/plain" {
				decodeBytes, err := base64.URLEncoding.DecodeString(item.Body.Data)
				if err != nil {
					log.Fatalln(err)
				}
				BodyType = "text"
				BodyContent = string(decodeBytes)
			} else if item.MimeType == "text/html" {
				decodeBytes, err := base64.URLEncoding.DecodeString(item.Body.Data)
				if err != nil {
					log.Fatalln(err)
				}
				BodyType = "html"
				BodyContent = string(decodeBytes)
			} else if item.Body.AttachmentId != "" {
				a, err := DownAttachmentFile(user, id, item.Body.AttachmentId, item.Filename)
				if err == nil {
					Attachments = append(Attachments, a)
				}

			}
		}
	}
	return model.Mail{
		Id:              m.Id,
		User:            user.UserPrincipalName,
		Subject:         subject,
		SenderEmail:     SenderEmail,
		SenderName:      SenderName,
		Attachments:     Attachments,
		BodyPreview:     m.Snippet,
		BodyType:        BodyType,
		BodyContent:     BodyContent,
		ToRecipient:     ToRecipient,
		ToRecipientName: ToRecipientName,
		Date:            time.Unix(m.InternalDate/1000, 0),
	}
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func GetPageEmails(user model.User, pageNumber int, pageSize int) []model.Mail {
	values := url.Values{}
	values.Add("maxResults", string(pageSize))
	var dbMails []model.Mail
	ids := strings.Split(user.MailIds, ",")
	var idsSlice []string
	if len(ids) < pageNumber*pageSize {
		idsSlice = ids[(pageNumber-1)*pageSize:]
	} else {
		idsSlice = ids[(pageNumber-1)*pageSize : pageNumber*pageSize]
	}
	var wg sync.WaitGroup
	for _, id := range idsSlice {
		wg.Add(1)
		go func(id string) {
			email := database.GetEmail(id)
			// 如果邮件不存在于DB中，那么拉取它
			if email.Id == "" {
				email = FetchOneMail(user, id)
				database.InsertEmail(email)
			}
			// 判断邮件中的附件文件是否存在
			for _, v := range email.Attachments {
				if ok, err := PathExists(v.FilePath); !ok && err == nil{
					DownAttachmentFile(user, v.MailId, v.AttachmentId, v.FileName)
				}
			}
			dbMails = append(dbMails, email)
			wg.Done()
		}(id)
	}
	wg.Wait()
	sort.Sort(model.MailList(dbMails))
	return dbMails

}
