package main

import (
	"fmt"
	"net/http"
)

func main() {
	var err error
	var courseMaster CourseMaster
	fmt.Println("Starting go-learn-it")

	// Setup our HTTP server.
	fs := http.FileServer(http.Dir("static/"))
	if courseMaster.Curriculum, err = CurriculumFromFilePath("./curriculum.json"); err != nil {
		fmt.Print(err)
		return
	}
	if err := courseMaster.buildTemplates(); err != nil {
		fmt.Print(err)
		return
	}
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/api/", courseMaster.handleAPI)
	http.HandleFunc("/", courseMaster.handleHTTP)
	http.ListenAndServe(":8888", nil)
}