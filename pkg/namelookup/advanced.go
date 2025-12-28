package namelookup

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedSearchByName performs enhanced people search with additional sources
func AdvancedSearchByName(fullName string) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced people search\n")
	// result.WriteString("This mode provides additional search sources\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard search first
	startTime := time.Now()
	standardResult, err := SearchByName(fullName)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	// Parse name
	firstName, lastName := parseName(fullName)
	fullNameEncoded := strings.ReplaceAll(fullName, " ", "%20")

	result.WriteString("Professional Networks:\n")
	result.WriteString(fmt.Sprintf("  - LinkedIn: https://www.linkedin.com/search/results/all/?keywords=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Indeed Resume: https://www.indeed.com/resumes?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - AngelList: https://angel.co/search?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Crunchbase: https://www.crunchbase.com/discover/people?q=%s\n", fullNameEncoded))

	result.WriteString("\nSocial Media Deep Search:\n")
	result.WriteString(fmt.Sprintf("  - Facebook People: https://www.facebook.com/search/people/?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Twitter Advanced: https://twitter.com/search?q=%s&f=user\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Instagram: https://www.instagram.com/explore/tags/%s/\n", strings.ToLower(strings.ReplaceAll(fullName, " ", ""))))
	result.WriteString(fmt.Sprintf("  - TikTok: https://www.tiktok.com/search/user?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Reddit: https://www.reddit.com/search/?q=%s\n", fullNameEncoded))

	result.WriteString("\nPeople Search Services:\n")
	result.WriteString(fmt.Sprintf("  - Whitepages: https://www.whitepages.com/name/%s-%s\n", firstName, lastName))
	result.WriteString(fmt.Sprintf("  - TruePeopleSearch: https://www.truepeoplesearch.com/results?name=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - FastPeopleSearch: https://www.fastpeoplesearch.com/name/%s-%s\n", firstName, lastName))
	result.WriteString(fmt.Sprintf("  - Spokeo: https://www.spokeo.com/%s-%s\n", firstName, lastName))
	result.WriteString(fmt.Sprintf("  - Pipl: https://pipl.com/search/?q=%s\n", fullNameEncoded))

	result.WriteString("\nPublic Records:\n")
	result.WriteString(fmt.Sprintf("  - Court Records: Search county clerk websites\n"))
	result.WriteString(fmt.Sprintf("  - Property Records: Search county assessor websites\n"))
	result.WriteString(fmt.Sprintf("  - Voter Registration: Search state voter databases\n"))
	result.WriteString(fmt.Sprintf("  - Business Filings: Search Secretary of State websites\n"))

	result.WriteString("\nProfessional & Academic:\n")
	result.WriteString(fmt.Sprintf("  - Google Scholar: https://scholar.google.com/scholar?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - ResearchGate: https://www.researchgate.net/search.Search.html?query=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - ORCID: https://orcid.org/orcid-search/search?searchQuery=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Academia.edu: https://www.academia.edu/search?q=%s\n", fullNameEncoded))

	result.WriteString("\nContent & Profiles:\n")
	result.WriteString(fmt.Sprintf("  - GitHub: https://github.com/search?q=%s&type=users\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Stack Overflow: https://stackoverflow.com/users?search=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Medium: https://medium.com/search/people?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - YouTube: https://www.youtube.com/results?search_query=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Vimeo: https://vimeo.com/search?q=%s\n", fullNameEncoded))

	result.WriteString("\nUsername Variations to Try:\n")
	result.WriteString(fmt.Sprintf("  - %s%s\n", strings.ToLower(firstName), strings.ToLower(lastName)))
	result.WriteString(fmt.Sprintf("  - %s.%s\n", strings.ToLower(firstName), strings.ToLower(lastName)))
	result.WriteString(fmt.Sprintf("  - %s_%s\n", strings.ToLower(firstName), strings.ToLower(lastName)))
	if len(firstName) > 0 && len(lastName) > 0 {
		result.WriteString(fmt.Sprintf("  - %c%s\n", strings.ToLower(firstName)[0], strings.ToLower(lastName)))
	}

	result.WriteString("\nData Breach Checking:\n")
	result.WriteString("  Check if name appears in data breaches:\n")
	result.WriteString(fmt.Sprintf("  - HIBP: https://haveibeenpwned.com/\n"))
	result.WriteString(fmt.Sprintf("  - DeHashed: https://dehashed.com/\n"))

	result.WriteString("\n⚠️  IMPORTANT NOTES:\n")
	result.WriteString("  - Always verify information from multiple sources\n")
	result.WriteString("  - Respect privacy laws and ethical boundaries\n")
	result.WriteString("  - Some services require paid subscriptions\n")
	result.WriteString("  - Public records vary by jurisdiction\n")
