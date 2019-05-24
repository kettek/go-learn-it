package main

import (
	"encoding/json"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"log"
	"path"
	"strconv"
)

// Course is our type for holding a collection of Lessons.
type Course struct {
	Name              string
	Shortname         string
	Description       string
	Language          string
	LanguageExtension string
	LessonTitles      []string `json:"Lessons"`
	Lessons           []Lesson
}

// CourseFromDirPath attempts to unmarshal the json from the given path into the course.
func CourseFromDirPath(dirPath string) (c Course, err error) {
	c.Shortname = path.Base(dirPath)
	// Unmarshal our JSON into our Course.
	var bytes []byte
	jsonPath := dirPath + ".json"
	if bytes, err = ioutil.ReadFile(jsonPath); err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &c); err != nil {
		return
	}
	// Construct our Lesson Files.
	for i, lessonTitle := range c.LessonTitles {
		var markdownSource, sourceCode []byte
		markdownPath := path.Join(dirPath, strconv.Itoa(i+1)) + ".md"
		sourcePath := path.Join(dirPath, strconv.Itoa(i+1)) + "." + c.LanguageExtension
		log.Printf("Attempting to read: %s %s\n", markdownPath, sourcePath)

		if markdownSource, err = ioutil.ReadFile(markdownPath); err != nil {
			log.Printf("Could not read markdown file for lesson(%d): %s\n", i+1, lessonTitle)
			continue
		}
		if sourceCode, err = ioutil.ReadFile(sourcePath); err != nil {
			log.Printf("Could not read source file for lesson(%d): %s\n", i+1, lessonTitle)
			continue
		}

		c.Lessons = append(c.Lessons, Lesson{
			Title:       lessonTitle,
			HTMLContent: string(markdown.ToHTML(markdownSource, nil, nil)),
			SourceCode:  string(sourceCode),
			Language:    c.Language,
		})
	}
	return
}
