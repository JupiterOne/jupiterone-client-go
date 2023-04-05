package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	j1 "github.com/jupiterone/jupiterone-client-go/jupiterone"
	"github.com/jupiterone/jupiterone-client-go/jupiterone/domain"
)

type UploadPayload struct {
	ID        string `json:"_id"`
	InKEVList bool   `json:"inKEVList"`
}

type CISAKEVVulnerability struct {
	CveID             string `json:"cveID"`
	VendorProject     string `json:"vendorProject"`
	Product           string `json:"product"`
	VulnerabilityName string `json:"vulnerabilityName"`
	DateAdded         string `json:"dateAdded"`
	ShortDescription  string `json:"shortDescription"`
	RequiredAction    string `json:"requiredAction"`
	DueDate           string `json:"dueDate"`
	Notes             string `json:"notes"`
}

type CISAKEV struct {
	Title           string                 `json:"title"`
	CatalogVersion  string                 `json:"catalogVersion"`
	DateReleased    time.Time              `json:"dateReleased"`
	Count           int                    `json:"count"`
	Vulnerabilities []CISAKEVVulnerability `json:"vulnerabilities"`
}

func getEnvWithDefault(key string, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultVal
	}
	return value
}

func main() {
	var cisakevList CISAKEV
	cisaVulnerabilitiesByName := make(map[string]CISAKEVVulnerability)
	uploadPayloads := []interface{}{}

	log.Println("fetching vulnerabilities from CISA KEV list...")

	resp, err := http.Get("https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json")
	if err != nil {
		log.Fatalf("error requesting vulnerabilities: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading body: %v", err)
	}

	json.Unmarshal(body, &cisakevList)

	log.Println("fetched vulnerabilities from CISA KEV list.")

	for _, vulnerability := range cisakevList.Vulnerabilities {
		cisaVulnerabilitiesByName[vulnerability.CveID] = vulnerability
	}

	// Set configuration
	config := j1.Config{
		APIKey:    getEnvWithDefault("J1_API_TOKEN", ""),
		AccountID: getEnvWithDefault("J1_ACCOUNT", ""),
		Region:    getEnvWithDefault("J1_REGION", "us"),
	}

	client, err := j1.NewClient(&config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	log.Println("querying for vulnerabilities within JupiterOne...")
	results, err := client.Query.Query(j1.QueryInput{
		Query: "FIND cve",
	})
	if err != nil {
		log.Fatalf("failed query output: %v", err)
	}

	resultsAsList, err := client.Query.AsList(results)
	if err != nil {
		log.Fatalf("failed to convert results to list: %v", err)
	}

	log.Println("vulnerabilities received from JupiterOne...")

	for _, result := range resultsAsList.Data {
		cveName := result.Entity.DisplayName
		_, ok := cisaVulnerabilitiesByName[cveName]
		if !ok {
			continue
		}

		up := UploadPayload{
			ID:        result.Entity.ID,
			InKEVList: true,
		}

		uploadPayloads = append(uploadPayloads, up)
	}

	log.Printf("found %d vulnerabilities in CISA KEV list\n", len(uploadPayloads))

	if len(uploadPayloads) == 0 {
		log.Print("no vulnerabilities to upload, exiting")
		os.Exit(0)
	}

	stp := domain.StartParams{
		Source:   "api",
		SyncMode: "CREATE_OR_UPDATE",
	}
	syp := domain.SyncPayload{
		Entities: uploadPayloads,
	}

	log.Println("uploading new data into JupiterOne...")

	output, err := client.Synchronization.ProcessSyncJob(stp, syp)
	if err != nil {
		log.Fatalf("failed to process sync job: %v", err)
	}

	var jsonOutput map[string]interface{}

	b, err := json.Marshal(*output)
	if err != nil {
		log.Fatalf("failed to marshal output: %v", err)
	}

	err = json.Unmarshal(b, &jsonOutput)
	if err != nil {
		log.Fatal("failed to unmarshal output")
	}

	jsonString, err := json.MarshalIndent(jsonOutput, "", "  ")
	if err != nil {
		log.Printf("failed to marshal output: %v\n", err)
	}

	log.Println(string(jsonString))

	log.Print("program complete.")
}
