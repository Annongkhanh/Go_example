package mail

import (
	"testing"

	"github.com/Annongkhanh/Simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("..")

	require.NoError(t, err)
	require.NotNil(t, config)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "test email subject"
	content := `
	<h1>Hello</h1>
	<p>This is a test email from <a href = "t2826342@gmail.com"> An </a></p>
	`

	to := []string{"t2826342@gmail.com"}
	attachFiles := []string{"../app.env"}
	err = sender.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		attachFiles,
	)
	require.NoError(t, err)
}
