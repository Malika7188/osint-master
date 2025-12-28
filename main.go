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
