package pdfgen

import (
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF creates a professional PDF report from OSINT results
func GeneratePDF(filename, title, content string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)

	// Add page
	pdf.AddPage()

	// Add header with branding
	addHeader(pdf, title)

	// Add metadata
	addMetadata(pdf)

	// Add content
	addContent(pdf, content)

	// Add footer
	addFooter(pdf)

	// Save PDF
	return pdf.OutputFileAndClose(filename)
}

// addHeader adds a professional header to the PDF
func addHeader(pdf *gofpdf.Fpdf, title string) {
	// Title background
	pdf.SetFillColor(41, 128, 185) // Blue background
	pdf.Rect(0, 0, 210, 40, "F")

	// Title text
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetY(10)
	pdf.CellFormat(0, 10, "OSINT Master Report", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.SetY(25)
	pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")

	// Reset text color
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(15)
}

// addMetadata adds report metadata
func addMetadata(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(100, 100, 100)

	// Report info
	pdf.CellFormat(0, 6, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Tool: OSINT Master v1.0.0", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Educational & Authorized Use Only", "", 1, "L", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)

	// Separator line
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(8)
}

// addContent adds the main content to the PDF
func addContent(pdf *gofpdf.Fpdf, content string) {
	pdf.SetFont("Arial", "", 11)

	// Split content into lines
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// Check for special formatting
		if strings.HasPrefix(line, "===") || strings.HasPrefix(line, "---") {
			// Separator
			pdf.Ln(2)
			pdf.SetDrawColor(220, 220, 220)
			pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
			pdf.Ln(2)
			continue
		}

		// Handle warnings and alerts
		if strings.Contains(line, "⚠️") || strings.Contains(line, "WARNING") {
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(231, 76, 60) // Red
			pdf.MultiCell(0, 6, cleanText(line), "", "L", false)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("Arial", "", 11)
			continue
		}

		// Handle success indicators
		if strings.Contains(line, "✓") || strings.Contains(line, "FOUND") {
			pdf.SetTextColor(39, 174, 96) // Green
			pdf.MultiCell(0, 6, cleanText(line), "", "L", false)
			pdf.SetTextColor(0, 0, 0)
			continue
		}

		// Handle section headers (lines with colons at the end)
		if strings.HasSuffix(strings.TrimSpace(line), ":") && len(line) < 50 {
			pdf.SetFont("Arial", "B", 12)
			pdf.SetTextColor(41, 128, 185) // Blue
			pdf.Ln(3)
			pdf.MultiCell(0, 6, cleanText(line), "", "L", false)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetFont("Arial", "", 11)
			continue
		}

		// Handle links (URLs)
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") || strings.Contains(line, "URL:") {
			pdf.SetTextColor(52, 152, 219) // Link blue
			pdf.MultiCell(0, 6, cleanText(line), "", "L", false)
			pdf.SetTextColor(0, 0, 0)
			continue
		}

		// Regular text
		if strings.TrimSpace(line) != "" {
			pdf.MultiCell(0, 6, cleanText(line), "", "L", false)
		} else {
			pdf.Ln(2)
		}
	}
}

// addFooter adds a professional footer
func addFooter(pdf *gofpdf.Fpdf) {
	pdf.SetY(-25)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(3)

	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(128, 128, 128)
	pdf.CellFormat(0, 5, "OSINT Master - Educational OSINT Tool", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, "Always obtain permission before gathering information", "", 1, "C", false, 0, "")

	pdf.SetY(-10)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
}

// cleanText removes special characters that might not render in PDF
func cleanText(text string) string {
	// Replace emojis and special chars with text equivalents
	replacements := map[string]string{
		"✓": "[OK]",
		"✗": "[X]",
		"⚠️":  "[WARNING]",
		"🔗": "[LINK]",
		"📡": "[API]",
		"💡": "[TIP]",
	}

	for emoji, replacement := range replacements {
		text = strings.ReplaceAll(text, emoji, replacement)
	}

	return text
}

// GenerateEmailPDF creates a PDF specifically for email lookup results
func GenerateEmailPDF(filename, email, content string) error {
	title := fmt.Sprintf("Email Lookup: %s", email)
	return GeneratePDF(filename, title, content)
}

// GeneratePhonePDF creates a PDF specifically for phone lookup results
func GeneratePhonePDF(filename, phone, content string) error {
	title := fmt.Sprintf("Phone Lookup: %s", phone)
	return GeneratePDF(filename, title, content)
}

// GenerateUsernamePDF creates a PDF specifically for username search results
func GenerateUsernamePDF(filename, username, content string) error {
	title := fmt.Sprintf("Username Search: @%s", username)
	return GeneratePDF(filename, title, content)
}

// GenerateIPPDF creates a PDF specifically for IP lookup results
func GenerateIPPDF(filename, ip, content string) error {
	title := fmt.Sprintf("IP Lookup: %s", ip)
	return GeneratePDF(filename, title, content)
}

// GenerateDomainPDF creates a PDF specifically for domain enumeration results
func GenerateDomainPDF(filename, domain, content string) error {
	title := fmt.Sprintf("Domain Enumeration: %s", domain)
	return GeneratePDF(filename, title, content)
}

// GenerateNamePDF creates a PDF specifically for name lookup results
func GenerateNamePDF(filename, name, content string) error {
	title := fmt.Sprintf("Name Lookup: %s", name)
	return GeneratePDF(filename, title, content)
}
