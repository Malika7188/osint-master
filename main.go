package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/malika/osint-master/config"
	"github.com/malika/osint-master/internal/output"
	"github.com/malika/osint-master/pkg/domain"
	"github.com/malika/osint-master/pkg/emaillookup"
	"github.com/malika/osint-master/pkg/iplookup"
	"github.com/malika/osint-master/pkg/namelookup"
	"github.com/malika/osint-master/pkg/pdfgen"
	"github.com/malika/osint-master/pkg/phonelookup"
	"github.com/malika/osint-master/pkg/username"
	"github.com/malika/osint-master/pkg/webserver"
)

const version = "1.0.0"

func main() {
	// Define command-line flags
	nameFlag := flag.String("n", "", "Search information by full name")
	ipFlag := flag.String("i", "", "Search information by IP address")
	usernameFlag := flag.String("u", "", "Search information by username")
	domainFlag := flag.String("d", "", "Enumerate subdomains and check for takeover risks")
	emailFlag := flag.String("e", "", "Search information by email address")
	phoneFlag := flag.String("p", "", "Search information by phone number")
	outputFlag := flag.String("o", "", "File name to save output")
	pdfFlag := flag.String("pdf", "", "Generate PDF report (specify filename)")
	webFlag := flag.String("web", "", "Start web GUI server (specify port, e.g., 8080)")
	advancedFlag := flag.Bool("advanced", false, "Use advanced mode (browser automation - slower but more accurate)")
	setupConfigFlag := flag.Bool("setup-config", false, "Create sample config file for API keys")
	helpFlag := flag.Bool("help", false, "Display help information")

	flag.Parse()

	// Handle setup-config command
	if *setupConfigFlag {
		if err := config.CreateSampleEnvFile(); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle web server mode
	if *webFlag != "" {
		if err := webserver.StartServer(*webFlag); err != nil {
			fmt.Printf("Error starting web server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Show help if requested or no flags provided
	if *helpFlag || flag.NFlag() == 0 {
		showHelp()
		return
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Validate that at least one search flag is provided
	if *nameFlag == "" && *ipFlag == "" && *usernameFlag == "" && *domainFlag == "" && *emailFlag == "" && *phoneFlag == "" {
		fmt.Println("Error: Please provide at least one search option (-n, -i, -u, -d, -e, or -p)")
		fmt.Println("Use --help for more information")
		os.Exit(1)
	}

	// Process based on the flag provided
	var result string
	var err error

	if *nameFlag != "" {
		fmt.Printf("Searching for: %s\n", *nameFlag)
		result, err = namelookup.SearchByName(*nameFlag)
	} else if *ipFlag != "" {
		fmt.Printf("Looking up IP: %s\n", *ipFlag)
		result, err = iplookup.LookupIP(*ipFlag)
	} else if *usernameFlag != "" {
		fmt.Printf("Searching for username: %s\n", *usernameFlag)

		// Use advanced mode if flag is set
		if *advancedFlag {
			result, err = username.AdvancedSearchUsername(*usernameFlag)
		} else {
			result, err = username.SearchUsername(*usernameFlag)
		}
	} else if *domainFlag != "" {
		fmt.Printf("Enumerating domain: %s\n", *domainFlag)
		result, err = domain.EnumerateDomain(*domainFlag)
	} else if *emailFlag != "" {
		fmt.Printf("Looking up email: %s\n", *emailFlag)
		result, err = emaillookup.LookupEmailWithConfig(*emailFlag, cfg.HIBPAPIKey)
	} else if *phoneFlag != "" {
		fmt.Printf("Looking up phone: %s\n", *phoneFlag)
		result, err = phonelookup.LookupPhone(*phoneFlag)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Println(result)

	// Save to file if output flag is provided
	if *outputFlag != "" {
		err = output.SaveToFile(*outputFlag, result)
		if err != nil {
			fmt.Printf("Error saving to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Data saved in %s\n", *outputFlag)
	}

	// Generate PDF if pdf flag is provided
	if *pdfFlag != "" {
		var pdfErr error
		if *emailFlag != "" {
			pdfErr = pdfgen.GenerateEmailPDF(*pdfFlag, *emailFlag, result)
		} else if *phoneFlag != "" {
			pdfErr = pdfgen.GeneratePhonePDF(*pdfFlag, *phoneFlag, result)
		} else if *usernameFlag != "" {
			pdfErr = pdfgen.GenerateUsernamePDF(*pdfFlag, *usernameFlag, result)
		} else if *ipFlag != "" {
			pdfErr = pdfgen.GenerateIPPDF(*pdfFlag, *ipFlag, result)
		} else if *domainFlag != "" {
			pdfErr = pdfgen.GenerateDomainPDF(*pdfFlag, *domainFlag, result)
		} else if *nameFlag != "" {
			pdfErr = pdfgen.GenerateNamePDF(*pdfFlag, *nameFlag, result)
		}

		if pdfErr != nil {
			fmt.Printf("Error generating PDF: %v\n", pdfErr)
			os.Exit(1)
		}
		fmt.Printf("PDF report generated: %s\n", *pdfFlag)
	}
}

func showHelp() {
	fmt.Println("\nWelcome to osintmaster multi-function Tool")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("OPTIONS:")
	fmt.Println("    -n  \"Full Name\"        Search information by full name")
	fmt.Println("    -i  \"IP Address\"       Search information by IP address")
	fmt.Println("    -u  \"Username\"         Search information by username")
	fmt.Println("    -d  \"Domain\"           Enumerate subdomains and check for takeover risks")
	fmt.Println("    -e  \"Email\"            Search information by email address")
	fmt.Println("    -p  \"Phone Number\"     Search information by phone number")
	fmt.Println("    -o  \"FileName\"         File name to save output")
	fmt.Println("    --pdf \"FileName.pdf\"   Generate professional PDF report")
	fmt.Println("    --web \"8080\"           Start web GUI server on specified port")
	fmt.Println("    --advanced             Use advanced mode with browser automation (slower)")
	fmt.Println("    --setup-config         Create sample API configuration file")
	fmt.Println("    --help                 Display this help message")
	fmt.Println("\nEXAMPLES:")
	fmt.Println("    osintmaster -n \"John Doe\" -o result.txt")
	fmt.Println("    osintmaster -i 8.8.8.8 -o ip_info.txt")
	fmt.Println("    osintmaster -u \"@username\" -o user_search.txt")
	fmt.Println("    osintmaster -u \"@username\" --advanced -o user_search.txt  (Advanced mode)")
	fmt.Println("    osintmaster -d \"example.com\" -o domain_info.txt")
	fmt.Println("    osintmaster -e \"email@example.com\" -o email_info.txt")
	fmt.Println("    osintmaster -e \"email@example.com\" --pdf report.pdf      (PDF report)")
	fmt.Println("    osintmaster -p \"+1234567890\" -o phone_info.txt")
	fmt.Println("    osintmaster --web 8080                                    (Start web GUI)")
	fmt.Println("\nCONFIGURATION:")
	fmt.Println("    osintmaster --setup-config         Create API config file")
	fmt.Println("    Config file location: ~/.osintmaster/.env")
	fmt.Println("    Add API keys to unlock full features (HIBP, phone lookup, etc.)")
	fmt.Println("\nETHICAL NOTICE:")
	fmt.Println("    This tool is for EDUCATIONAL PURPOSES ONLY.")
	fmt.Println("    Always obtain permission before gathering information.")
	fmt.Println("    Respect privacy and comply with all applicable laws.")
	fmt.Println("\nADVANCED MODE WARNING:")
	fmt.Println("    --advanced flag uses browser automation to bypass basic bot detection.")
	fmt.Println("    Use ONLY for:")
	fmt.Println("      ✓ Educational learning")
	fmt.Println("      ✓ Authorized penetration testing (with written permission)")
	fmt.Println("      ✓ Testing on YOUR OWN systems")
	fmt.Println("DO NOT use for unauthorized access or ToS violations.")
}
