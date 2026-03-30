package writer

// MailMessage is the fully-rendered email ready for serialisation.
type MailMessage struct {
	From        string
	To          []string // one or more recipients
	Cc          []string // optional
	Bcc         []string // optional
	Subject     string
	Body        string // HTML content
	Attachments []Attachment
}

// Attachment holds a single file to be embedded in multipart/mixed.
type Attachment struct {
	Filename    string // display name in email
	ContentType string // e.g. "application/pdf"
	Data        []byte // raw file bytes (Base64 encoded by EMLWriter)
}

// Writer is the pluggable output interface.
type Writer interface {
	Write(msg MailMessage, month string) error
}
