package email

type EmailContent struct {
	To                   []string
	Subject              string
	Message              *string
	HtmlTemplateFileName *string
	TextTemplateFileName *string
	TemplateValueMap     map[string]string
}

type SMTPConfig struct {
	UserName   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Override   bool   `mapstructure:"override"`
	OverrideTo string `mapstructure:"overrideTo"`
}

type emailType string

const (
	NIL  emailType = "_"
	HTML emailType = "HTML"
	TEXT emailType = "TEXT"
)
