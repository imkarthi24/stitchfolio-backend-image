package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	// Test case 1: Valid email with HTML template
	// Note: This test is commented out as it requires actual SMTP server
	// config := &config.SMTPConfig{
	// 	Host:     "smtp.gmail.com",
	// 	Port:     587,
	// 	UserName: "test@example.com",
	// 	Password: "password",
	// }

	// message := "Test message"
	// htmlTemplate := "test_template.html"
	// mailContent := EmailContent{
	// 	To:                   "recipient@example.com",
	// 	Subject:             "Test Subject",
	// 	Message:             &message,
	// 	HtmlTemplateFileName: &htmlTemplate,
	// 	TemplateValueMap:    map[string]string{"key": "value"},
	// }

	// err := SendEmail(config, mailContent)
	// require.Nil(t, err)
}

func TestBuildEmailBody(t *testing.T) {
	// Test case 1: HTML template
	message := "Test message"
	htmlTemplate := "test_template.html"
	mailContent := EmailContent{
		Message:              &message,
		HtmlTemplateFileName: &htmlTemplate,
		TemplateValueMap:     map[string]string{"key": "value"},
	}

	mailType, content, err := BuildEmailBody(mailContent)
	require.NotNil(t, err) // Will fail as template file doesn't exist
	require.Equal(t, NIL, mailType)
	require.Nil(t, content)

	// Test case 2: Text template
	textTemplate := "test_template.txt"
	mailContent = EmailContent{
		Message:              &message,
		TextTemplateFileName: &textTemplate,
		TemplateValueMap:     map[string]string{"key": "value"},
	}

	mailType, content, err = BuildEmailBody(mailContent)
	require.NotNil(t, err) // Will fail as template file doesn't exist
	require.Equal(t, NIL, mailType)
	require.Nil(t, content)

	// Test case 3: No template, just message
	mailContent = EmailContent{
		Message: &message,
	}

	mailType, content, err = BuildEmailBody(mailContent)
	require.Nil(t, err)
	require.Equal(t, NIL, mailType)
	require.Nil(t, content)
}

func TestReadContentFromFile(t *testing.T) {
	// Test case 1: Non-existent file
	content, err := readContentFromFile("non_existent_file.txt")
	require.NotNil(t, err)
	require.Nil(t, content)

	// Test case 2: Valid file (would need to create a test file first)
	// This test would require creating a test file in the workspace
	// content, err = readContentFromFile("test_file.txt")
	// require.Nil(t, err)
	// require.NotNil(t, content)
}
