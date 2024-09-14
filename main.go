package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Domain struct {
	Name        string
	ServiceName string
	Status      string
}

var domains = []Domain{}
var monitoring = true
var apiToken string
var chatID string
var checkInterval = 5 * time.Minute // Check interval set to 5 minutes

// Load the API token and chat ID from api.conf
func loadAPIConfig() {
	file, err := os.Open("api.conf")
	if err != nil {
		fmt.Println("Error opening api.conf:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "api_token" {
				apiToken = value
			} else if key == "chat_id" {
				chatID = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading api.conf:", err)
	}
}

// Save the API token and chat ID to api.conf
func saveAPIConfig() {
	file, err := os.Create("api.conf")
	if err != nil {
		fmt.Println("Error creating api.conf:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "api_token=%s\n", apiToken)
	fmt.Fprintf(writer, "chat_id=%s\n", chatID)
	writer.Flush()
}

// Load the domains and their service names from domain.conf
func loadDomains() {
	file, err := os.Open("domain.conf")
	if err != nil {
		fmt.Println("Error opening domain.conf:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			// Split the domain and service name using "#"
			parts := strings.Split(line, "#")
			domainName := strings.TrimSpace(parts[0])
			serviceName := ""
			if len(parts) > 1 {
				serviceName = strings.TrimSpace(parts[1])
			}
			// Add the domain with its service name and default status "OK"
			domains = append(domains, Domain{Name: domainName, ServiceName: serviceName, Status: "OK"})
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading domain.conf:", err)
	}
}

// Save the domains and their service names to domain.conf
func saveDomains() {
	file, err := os.Create("domain.conf")
	if err != nil {
		fmt.Println("Error creating domain.conf:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, domain := range domains {
		fmt.Fprintf(writer, "%s #%s\n", domain.Name, domain.ServiceName)
	}
	writer.Flush()
}

// Function to send Telegram notifications
func sendTelegramMessage(message string) {
	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	req, _ := http.NewRequest("POST", telegramAPI, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()
}

// Function to check server status
func checkServerStatus(domain string) string {
	resp, err := http.Get(domain)
	if err != nil || resp.StatusCode != 200 {
		return "ERR"
	}
	return "OK"
}

// Monitoring function with interval checks
func monitorServers() {
	for {
		if monitoring {
			checkAllDomains() // Trigger domain checks
		}
		// Wait for the next check
		time.Sleep(checkInterval)
	}
}

// Function to check all domains and send notifications if status changes
func checkAllDomains() {
	for i, domain := range domains {
		previousStatus := domains[i].Status
		currentStatus := checkServerStatus(domain.Name)

		// If the status changes, send a Telegram notification
		if currentStatus != previousStatus {
			domains[i].Status = currentStatus
			message := fmt.Sprintf("Status changed for %s (%s): %s", domain.Name, domain.ServiceName, currentStatus)
			sendTelegramMessage(message)
		}
	}
}

// Render the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	data := map[string]interface{}{
		"domains":    domains,
		"monitoring": monitoring,
		"api_token":  apiToken,
		"chat_id":    chatID,
	}
	tmpl.Execute(w, data)
}

// Handle adding a new domain
func addDomainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		newDomain := r.FormValue("new_domain")
		serviceName := r.FormValue("service_name")
		// Add the new domain with a default status of "OK" and the provided service name
		domains = append(domains, Domain{Name: newDomain, ServiceName: serviceName, Status: "OK"})
		saveDomains() // Save domains to file
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handle deleting a domain
func deleteDomainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		domainToDelete := r.FormValue("domain")
		// Loop through domains and delete the one that matches
		for i, domain := range domains {
			if domain.Name == domainToDelete {
				domains = append(domains[:i], domains[i+1:]...)
				break
			}
		}
		saveDomains() // Save updated domains to file
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handle toggling the monitoring status
func toggleMonitoringHandler(w http.ResponseWriter, r *http.Request) {
	monitoring = !monitoring
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handle updating API settings and saving them to api.conf
func updateAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		apiToken = r.FormValue("telegram_bot_id")
		chatID = r.FormValue("chat_id")
		saveAPIConfig() // Save the updated token and chat ID
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handle immediate "Check Now" request
func checkNowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		checkAllDomains() // Trigger immediate check of all domains
		w.Write([]byte("Server check triggered"))
	}
}

func main() {
	// Load API token, chat ID, and domains from config files
	loadAPIConfig()
	loadDomains()

	// Set up handlers for each route
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add_domain", addDomainHandler)
	http.HandleFunc("/delete_domain", deleteDomainHandler)
	http.HandleFunc("/toggle_monitoring", toggleMonitoringHandler)
	http.HandleFunc("/update_api", updateAPIHandler)
	http.HandleFunc("/check_now", checkNowHandler)

	// Serve static files (for background image, CSS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server explicitly on port 5001
	fmt.Println("Server starting on port 5001...")
	go monitorServers() // Start monitoring servers in the background
	http.ListenAndServe(":5001", nil)
}
