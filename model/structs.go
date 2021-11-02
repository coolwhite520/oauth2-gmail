package model

import "time"

type Config struct {
	Server struct {
		Host         string
		ExternalPort int
		InternalPort int
	}
	Oauth struct {
		ClientId     string
		ClientSecret string
		Scope        string
		Redirecturi  string
		AccessType   string
	}
}

type AuthResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type RecentFiles struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		OdataType string `json:"@odata.type"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		WebURL    string `json:"webUrl"`
	} `json:"value"`
}

type User struct {
	Id                string `json:"id"`
	DisplayName       string `json:"name"`
	Mail              string `json:"email"`
	JobTitle          string `json:"job_title"`
	UserPrincipalName string `json:"user_principal_name"`
	AccessToken       string `json:"access_token"`
	AccessTokenActive int    `json:"access_token_active"`
	RefreshToken      string `json:"refresh_token"`
	MailIds           string `json:"mail_ids"`
}

type MessageMail struct {
	Id           string `json:"id"`
	InternalDate int64  `json:"internalDate,string"`
	Snippet      string `json:"snippet"`
	Payload      struct {
		Headers []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"headers"`

		Body struct {
			AttachmentId string `json:"attachmentId"`
			Size         int    `json:"size"`
			Data         string `json:"data"`
		} `json:"body"`

		Parts []struct {
			PartId   string `json:"partId"`
			MimeType string `json:"mimeType"`
			Filename string `json:"filename"`
			Headers  []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"headers"`
			Body struct {
				AttachmentId string `json:"attachmentId"`
				Size         int    `json:"size"`
				Data         string `json:"data"`
			} `json:"body"`
		} `json:"parts"`
	}
}

type Mail struct {
	Id              string
	User            string
	Subject         string
	SenderEmail     string
	SenderName      string
	Attachments     []Attachment
	BodyPreview     string
	BodyType        string
	BodyContent     string
	ToRecipient     string
	ToRecipientName string
	Date            time.Time
}
type MailList []Mail

func (I MailList) Len() int {
	return len(I)
}
func (I MailList) Less(i, j int) bool {
	return I[i].Date.Unix() > I[j].Date.Unix()
}
func (I MailList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

type Attachment struct {
	MailId       string
	FileName     string
	FilePath     string
	AttachmentId string
}

type SendEmailStruct struct {
	Message struct {
		Subject string `json:"subject"`
		Body    struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"body"`
		ToRecipients []ToRecipients `json:"toRecipients"`
		Attachments  []Attachment   `json:"attachments,omitempty"`
	} `json:"message"`
	SaveToSentItems string `json:"saveToSentItems"`
}

type ToRecipients struct {
	EmailAddress EmailAddress `json:"emailAddress"`
}

type EmailAddress struct {
	Address string `json:"address"`
}

type SingleMail struct {
	CreatedDateTime      time.Time `json:"createdDateTime"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	ReceivedDateTime     time.Time `json:"receivedDateTime"`
	SentDateTime         time.Time `json:"sentDateTime"`
	Attachments          bool      `json:"Attachments"`
	InternetMessageID    string    `json:"internetMessageId"`
	Subject              string    `json:"subject"`
	BodyPreview          string    `json:"bodyPreview"`
	Body                 struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	Sender struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"sender"`
}

type Rule struct {
	DisplayName string `json:"displayName"`
	Sequence    int    `json:"sequence"`
	IsEnabled   bool   `json:"isEnabled"`
	Conditions  struct {
		SenderContains []string `json:"senderContains"`
	} `json:"conditions"`
	Actions struct {
		ForwardTo []struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"forwardTo"`
		StopProcessingRules bool `json:"stopProcessingRules"`
	} `json:"actions"`
}

type Messages struct {
	ResultSizeEstimate int    `json:"resultSizeEstimate"`
	NextPageToken      string `json:"nextPageToken"`
	Value              []struct {
		Id       string `json:"id"`
		ThreadId string `json:"threadId"`
	} `json:"messages"`
}

type Page struct {
	Title             string
	Email             string
	User              User
	UserList          []User
	EmailList         []Mail
	FileList          []string
	SearchFiles       Files
	Mail              SingleMail
	Message           string
	Success           bool
	URL               string
	Sum               int
	PageSize          int
	CurrentPageNumber int
	PageList          []string
}

type Drives struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		Description          string    `json:"description"`
		ID                   string    `json:"id"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		Name                 string    `json:"name"`
		WebURL               string    `json:"webUrl"`
		DriveType            string    `json:"driveType"`
		CreatedBy            struct {
			User struct {
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"createdBy"`
		LastModifiedBy struct {
			User struct {
				Email       string `json:"email"`
				ID          string `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"lastModifiedBy"`
		Owner struct {
			User struct {
				Email       string `json:"email"`
				ID          string `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"owner"`
		Quota struct {
			Deleted   int    `json:"deleted"`
			Remaining int64  `json:"remaining"`
			State     string `json:"state"`
			Total     int64  `json:"total"`
			Used      int    `json:"used"`
		} `json:"quota"`
	} `json:"value"`
}
type Files struct {
	OdataContext  string `json:"@odata.context"`
	NextPageToken string `json:"@odata.nextLink"`
	Value         []struct {
		OdataType            string    `json:"@odata.type"`
		CreatedDateTime      time.Time `json:"createdDateTime"`
		ID                   string    `json:"id"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		Name                 string    `json:"name"`
		WebURL               string    `json:"webUrl"`
		Size                 int       `json:"size"`
		ParentReference      struct {
			DriveID   string `json:"driveId"`
			DriveType string `json:"driveType"`
			ID        string `json:"id"`
		} `json:"parentReference"`
		File struct {
			MimeType string `json:"mimeType"`
		} `json:"file,omitempty"`
		FileSystemInfo struct {
			CreatedDateTime      time.Time `json:"createdDateTime"`
			LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		} `json:"fileSystemInfo"`
		SearchResult struct {
		} `json:"searchResult"`
		Folder struct {
			ChildCount int `json:"childCount"`
		} `json:"folder,omitempty"`
	} `json:"value"`
}

type DriveItem struct {
	OdataContext              string    `json:"@odata.context"`
	MicrosoftGraphDownloadURL string    `json:"@microsoft.graph.downloadUrl"`
	CreatedDateTime           time.Time `json:"createdDateTime"`
	ETag                      string    `json:"eTag"`
	ID                        string    `json:"id"`
	LastModifiedDateTime      time.Time `json:"lastModifiedDateTime"`
	Name                      string    `json:"name"`
	WebURL                    string    `json:"webUrl"`
	CTag                      string    `json:"cTag"`
	Size                      int       `json:"size"`
	CreatedBy                 struct {
		Application struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"application"`
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"createdBy"`
	LastModifiedBy struct {
		Application struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"application"`
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"lastModifiedBy"`
	ParentReference struct {
		DriveID   string `json:"driveId"`
		DriveType string `json:"driveType"`
		ID        string `json:"id"`
		Path      string `json:"path"`
	} `json:"parentReference"`
	File struct {
		MimeType string `json:"mimeType"`
	} `json:"file"`
	FileSystemInfo struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	} `json:"fileSystemInfo"`
}
