package utils

import (
	"educlypro/config"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
	"time"
)

func SendResetPasswordEmail(to, username string, token string) error {
	frontendURL := config.Cfg.Server.FrontendURL
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	data := struct {
		AppName   string
		Username  string
		Year      int
		ResetLink string
	}{
		AppName:   "EduclyPro",
		Username:  username,
		Year:      time.Now().Year(),
		ResetLink: resetLink,
	}

	// Load both the layout and the specific template
	layoutPath := filepath.Join("templates", "emails", "layout.html")
	tmplPath := filepath.Join("templates", "emails", "reset-password.html")

	tmpl, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	var body bytes.Buffer
	// Execute the "layout" template (defined in layout.html)
	if err := tmpl.ExecuteTemplate(&body, "layout", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	subject := "Reset Your Password"
	return SendEmail(to, subject, body.String())
}

func SendWelcomeEmail(to, username string) error {
	data := struct {
		Username     string
		DashboardURL string
	}{
		Username:     username,
		DashboardURL: config.Cfg.Server.FrontendURL,
	}

	layoutPath := filepath.Join("templates", "emails", "layout.html")
	tmplPath := filepath.Join("templates", "emails", "welcome.html")

	tmpl, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "layout", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	subject := "Welcome to EduclyPro!"
	return SendEmail(to, subject, body.String())
}

func SendAdminRegistrationEmail(to, username, userEmail string) error {
	data := struct {
		Username string
		Email    string
	}{
		Username: username,
		Email:    userEmail,
	}

	layoutPath := filepath.Join("templates", "emails", "layout.html")
	tmplPath := filepath.Join("templates", "emails", "user-registered.html")

	tmpl, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "layout", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	subject := "New User Registered"
	return SendEmail(to, subject, body.String())
}

func SendEmail(to, subject, body string) error {
	cfg := config.Cfg.Mail

	log.Printf("DEBUG: Attempting to send email to %s via %s:%d", to, cfg.Host, cfg.Port)

	// Standard SMTP message
	message := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s", cfg.FromName, cfg.FromAddress, to, subject, body)

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	err := smtp.SendMail(addr, auth, cfg.FromAddress, []string{to}, []byte(message))
	if err != nil {
		log.Printf("DEBUG: SMTP SendMail error: %v", err)
		return err
	}

	log.Printf("DEBUG: Email sent successfully to %s", to)
	return nil
}
