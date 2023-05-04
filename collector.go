package main

import (
	"collector/impression"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Setup a request structure
type Request struct {
	User string `json:"user"`
	Url  string `json:"url"`
}

// Return the port to use. Use the PORT environment variable if it exists
// otherwise default to 8080.
func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func getFilename(user string) string {
	// TODO: Load the user file
	// TODO: Verify that this user is in the file
	// TODO: Return the filename
	var name string
	if user == "protectivemetrics" {
		name = "ORBY38BB"
	}
	return name
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Add the CORS and type headers
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// If this is an options request return just headers with a 204
	if r.Method == "OPTIONS" {
		w.WriteHeader(204)
		return
	}
	// Read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respond(400, "Bad JSON body", w)
		log.Printf("Unable to read the request body: %v", err)
		return
	}
	// Unmarshal the body into a variable
	var request Request
	err = json.Unmarshal(body, &request)
	// If there's an error return a 400 because of the bad request body
	if err != nil {
		respond(400, "Bad JSON body", w)
		log.Printf("Unable to unmarshal the request body: %v", err)
		return
	}
	// Grab the filename for this user
	user := getFilename(request.User)
	// Open the SQLite database
	db, err := sql.Open("sqlite3", user+".sqlite")
	if err != nil {
		respond(400, "Database error", w)
		log.Printf("Unable to open the database: %v", err)
		return
	}
	// Instantiate the database repository
	impressionRepository := impression.NewSQLiteRepository(db)
	// TODO: Create tables and indexes if they don't exist (could be slow)
	err = impressionRepository.Migrate()
	if err != nil {
		respond(400, "Database error", w)
		log.Printf("Unable to create tables and indexes: %v", err)
		return
	}
	// Write the impression
	data := impression.Impression{
		Url: request.Url,
	}
	err = impressionRepository.Create(data)
	if err != nil {
		respond(400, "Database error", w)
		log.Printf("Unable to record an impression: %v", err)
		return
	}
	// Respond 200 OK
	respond(200, "OK", w)
}

func livezHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func respond(code int, message string, w http.ResponseWriter) {
	w.WriteHeader(code)
	io.WriteString(w, message)
}

func main() {
	// Grab the port from the environment variable or set it to a default
	port := getServerPort()
	// Set the root handler and start the server
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/livez", livezHandler)
	fmt.Printf("Server started on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Unable to run server: %v", err)
	}
}
