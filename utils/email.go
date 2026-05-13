package utils

import (
	"fmt"
	"log/slog"

	"cms-backend/config"

	"github.com/resend/resend-go/v2"
)

func SendEmail(to, subject, body string) error {
	client := resend.NewClient(config.SMTP_PASSWORD)
	params := &resend.SendEmailRequest{
		From:    config.SMTP_FROM,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}
	_, err := client.Emails.Send(params)
	return err
}

func SendProjectWelcomeEmail(toEmail, clientName, projectTitle, magicLink string) {
	subject := fmt.Sprintf("Your Project %q is Now Live", projectTitle)
	logoHTML := ""
	if config.COMPANY_LOGO_URL != "" {
		logoHTML = fmt.Sprintf(`<img src="%s" alt="%s" style="max-width:180px">`, config.COMPANY_LOGO_URL, config.COMPANY_NAME)
	}
	body := `<!DOCTYPE html><html><body style="font-family:sans-serif;background:#f3f4f6;padding:40px 20px">
<div style="max-width:600px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden">
<div style="background:#1e1e20;padding:28px 32px;text-align:center">` + logoHTML + `</div>
<div style="padding:40px">
<p style="font-size:24px;font-weight:700;margin:0 0 12px">Hello ` + clientName + `!</p>
<p style="color:#6b7280">Your project <strong>` + projectTitle + `</strong> is now live. Track progress below.</p>
<div style="text-align:center;margin:28px 0">
<a href="` + magicLink + `" style="background:#0c4196;color:#fff;text-decoration:none;padding:14px 32px;border-radius:8px;font-weight:600">View My Project</a>
</div>
</div>
<div style="padding:24px 40px;background:#f9fafb;text-align:center;font-size:13px;color:#9ca3af">
<strong>` + config.COMPANY_NAME + `</strong> | ` + config.COMPANY_PHONE + ` | ` + config.COMPANY_EMAIL + `
</div></div></body></html>`

	go func() {
		if err := SendEmail(toEmail, subject, body); err != nil {
			slog.Error("failed to send welcome email", "to", toEmail, "err", err)
		}
	}()
}
