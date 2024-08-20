package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Project struct {
	Name        string `json:"name"`
	Stack       string `json:"stack"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImagePath   string `json:"image"`
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT environment variable is not set
	}

	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}
// Custom function to extract the file extension
func ext(path string) string {
	return filepath.Ext(path)
}

func parseTemplateFiles(filenames ...string) (*template.Template, error) {
	// Create a FuncMap with your custom function
	funcMap := template.FuncMap{
		"ext": ext,
	}

	paths := make([]string, len(filenames))
	for i, file := range filenames {
		paths[i] = filepath.Join("templates", file)
	}
	// Parse the templates and apply the FuncMap
	return template.New("").Funcs(funcMap).ParseFiles(paths...)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseTemplateFiles("base.html", "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template: %s\n", err.Error())
		return
	}
	projects, err := loadProjects("projects.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error loading projects: %s\n", err.Error())
		return
	}
	data := struct {
		Title    string
		Projects []Project
	}{
		Title:    "RaShunda Williams Dev Portfolio | Remote",
		Projects: projects,
	}
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %s\n", err.Error())
	}
}

func loadProjects(filename string) ([]Project, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var projects []Project
	if err := json.NewDecoder(file).Decode(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}
