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
