package phonelookup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/malika/osint-master/config"
)

// PhoneInfo holds information about a phone number
type PhoneInfo struct {
	Number         string
	CountryCode    string
	Country        string
	Region         string
	Carrier        string
	LineType       string
	IsValid        bool
	OnWhatsApp     bool
	OnTelegram     bool
	OnSignal       bool
	OnViber        bool
	OnWeChat       bool
	OnLine         bool
	WhatsAppStatus string
	TelegramStatus string
	