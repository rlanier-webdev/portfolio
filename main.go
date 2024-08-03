package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Project struct {
	Name        string
	Description string
	URL         string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}

func parseTemplateFiles(filenames ...string) (*template.Template, error) {
	paths := make([]string, len(filenames))
	for i, file := range filenames {
		paths[i] = filepath.Join("templates", file)
	}
	return template.ParseFiles(paths...)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseTemplateFiles("base.html", "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template: %s\n", err.Error())
		return
	}
	data := struct {
		Title string
	}{
		Title: "Home Page",
	}
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %s\n", err.Error())
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseTemplateFiles("base.html", "project.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template: %s\n", err.Error())
		return
	}
	projects := []Project{
		{Name: "Project One", Description: "Description of project one", URL: "#"},
		{Name: "Project Two", Description: "Description of project two", URL: "#"},
	}
	data := struct {
		Title    string
		Projects []Project
	}{
		Title:    "Projects Page",
		Projects: projects,
	}
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %s\n", err.Error())
	}
}
