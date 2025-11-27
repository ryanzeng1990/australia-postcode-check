package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	// Goquery is an excellent HTML parser, similar to jQuery or BeautifulSoup.
	// You will need to install it: go get github.com/PuerkitoBio/goquery
	"github.com/PuerkitoBio/goquery"
)

type PostcodeResult struct {
	Postcode string `json:"postcode"`
	Suburb   string `json:"suburb"`
	State    string `json:"state"`
}

// Base URL for the Australia Post postcode search.
const BASE_URL = "https://auspost.com.au/postcode/"

// --- Handlers ---

// postcodeHandler handles the /search API endpoint.
// It expects a 'keyword' query parameter.
func postcodeHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to ensure the client knows to expect JSON
	w.Header().Set("Content-Type", "application/json")

	// Extract the keyword from the URL query parameters
	keyword := r.URL.Query().Get("keyword")

	if keyword == "" {
		w.WriteHeader(http.StatusBadRequest)
		// Write a structured error response
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing 'keyword' parameter in the query string. Example: /search?keyword=sydney"})
		return
	}

	// Add a small delay to be polite to the server we are scraping (good practice)
	time.Sleep(500 * time.Millisecond)

	// Call the scraping function
	jsonOutput := searchPostcodes(keyword)

	// Note: searchPostcodes returns a JSON string, which we write directly.
	// If the output contains an error or no-result message, we set a 500 status.
	if strings.Contains(jsonOutput, `"error"`) || strings.Contains(jsonOutput, `"message"`) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(jsonOutput))
}

// --- Scraper Logic ---
// searchPostcodes fetches and scrapes the postcode data for a given keyword.
func searchPostcodes(keyword string) string {
	if keyword == "" {
		return `{"error": "Keyword cannot be empty."}`
	}

	// Construct the target URL.
	url := fmt.Sprintf("%s%s", BASE_URL, strings.ToLower(strings.TrimSpace(keyword)))
	log.Printf("Scraping target: %s", url)

	// 1. Make the HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to create request: %s"}`, err)
	}

	// Use a common user-agent header to mimic a regular browser visit
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to fetch the page: %s"}`, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf(`{"error": "Received non-OK HTTP status: %d"}`, resp.StatusCode)
	}

	// 2. Parse the HTML content using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to parse HTML: %s"}`, err)
	}

	resultsList := []PostcodeResult{}

	// --- IMPORTANT: TARGETING THE RESULTS TABLE ---
	// Selector found via inspection: <table class="resultsList fn_tableResultsList fn_tablePostcodeList"...
	const postcodeTableSelector = "table.fn_tablePostcodeList"

	dataContainerFound := false
	// Find all table rows (<tr>) within the selected table
	doc.Find(postcodeTableSelector + " tr").Each(func(i int, row *goquery.Selection) {
		dataContainerFound = true

		// Skip the header row (i=0)
		if i == 0 {
			return
		}

		// Find all table data cells (<td>) in the current row
		cols := row.Find("td")

		// Columns are: 0=Postcode, 1=Suburb, 2=Category
		if cols.Length() >= 2 {
			// Postcode is in the first column (index 0)
			postcodeText := strings.TrimSpace(cols.Eq(0).Text())

			// Suburb (with State) is in the second column (index 1)
			fullSuburbText := strings.TrimSpace(cols.Eq(1).Text())

			// Split the text by comma and space, e.g., "SYDNEY, NSW" -> ["SYDNEY", "NSW"]
			parts := strings.Split(fullSuburbText, ",")

			suburb := ""
			state := ""

			// The first part is always the Suburb
			if len(parts) >= 1 {
				suburb = strings.TrimSpace(parts[0])
			}
			// The second part is the State
			if len(parts) >= 2 {
				state = strings.TrimSpace(parts[1])
			}

			if postcodeText != "" && suburb != "" {
				resultsList = append(resultsList, PostcodeResult{
					Postcode: postcodeText,
					Suburb:   suburb,
					State:    state,
				})
			}
		}
	})

	// Log a warning if the selector fails, but allow the API to return a no-results message.
	if !dataContainerFound {
		log.Printf("Warning: Selector '%s' did not find any elements for keyword '%s'.", postcodeTableSelector, keyword)
	}

	if len(resultsList) == 0 {
		return fmt.Sprintf(`{"message": "No postcodes found for keyword '%s'. Please verify the CSS selectors."}`, keyword)
	}

	// 3. Return the data as a JSON string
	jsonOutput, err := json.MarshalIndent(resultsList, "", "    ")
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to marshal JSON: %s"}`, err)
	}

	return string(jsonOutput)
}

func main() {
	http.HandleFunc("/search", postcodeHandler)
	port := "8080"
	log.Printf("Starting postcode API server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
