CREATE TABLE IF NOT EXISTS "users" (
	"id"	TEXT NOT NULL UNIQUE,
	"DisplayName"	TEXT NOT NULL,
	"Mail"	TEXT NOT NULL UNIQUE,
	"JobTitle"	TEXT,
	"UserPrincipalName"	TEXT,
	"AccessToken"	TEXT NOT NULL UNIQUE,
	"AccessTokenActive"	INTEGER,
	"RefreshToken"	TEXT,
	"MailIds"	TEXT,
	PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "mails" (
   "id"	TEXT NOT NULL UNIQUE,
   "User"	TEXT ,
   "Subject"	TEXT ,
   "SenderEmail"	TEXT ,
   "SenderName"	TEXT ,
   "Attachments"	TEXT,
   "BodyPreview"	TEXT ,
   "BodyType"	TEXT ,
   "BodyContent"	TEXT,
   "ToRecipient"	TEXT,
   "ToRecipientName"	TEXT,
   "Date"	datetime,
   PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "attachment" (
   "mailId"	TEXT ,
   "filename"	TEXT ,
   "filepath" TEXT,
   "attachmentId"	TEXT,
   PRIMARY KEY("mailId", "filename")
);
