package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
	"github.com/kettek/go-multipath"
)

// CourseMaster handles all requests to "/" and delineates accordingly.
type CourseMaster struct {
	multiPath      multipath.Multipath
	Curriculum     Curriculum
	LessonTemplate *template.Template
	ListTemplate   *template.Template
}

// ListTemplateData is the structure we pass to the list template.
type ListTemplateData struct {
	Name        string
	Title       string
	Description string
	Courses     []Course
}

// LessonTemplateData is the structure we pass to the lesson template.
type LessonTemplateData struct {
	CurriculumName string
	LessonIndex    int
	LessonsCount   int
	Lesson         *Lesson
}

func (h *CourseMaster) buildTemplates() (err error) {
	var bytes []byte

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
	}

	if bytes, err = h.multiPath.ReadFile(path.Join("templates", "lesson.gohtml")); err != nil {
		return
	}
	if h.LessonTemplate, err = template.New("lesson").Funcs(funcMap).Parse(string(bytes)); err != nil {
		return
	}

	if bytes, err = h.multiPath.ReadFile(path.Join("templates", "list.gohtml")); err != nil {
		return
	}
	if h.ListTemplate, err = template.New("list").Funcs(funcMap).Parse(string(bytes)); err != nil {
		return
	}

	return
}

// handleAPI sends a JSON response with a given course and lesson.
func (h *CourseMaster) handleAPI(w http.ResponseWriter, r *http.Request) {
	/*var err error
	var coursePath, lessonPath string

	coursePath, lessonPath, err = h.getCourseAndLessonStrings(r.URL.Path)*/
}

func (h *CourseMaster) handleHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var coursePath, lessonPath string
	var lessonIndex int

	coursePath, lessonPath, err = h.getCourseAndLessonStrings(r.URL.Path)
	// Redirect to list if an error occurred.
	if lessonPath != "" && err != nil {
		http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
		return
	}

	switch coursePath {
	case "":
		// Redirect to the welcome page if course is empty.
		http.Redirect(w, r, "/welcome/1", http.StatusTemporaryRedirect)
	case "list":
		if err = h.ListTemplate.Execute(w, ListTemplateData{
			Title:       "List of Courses",
			Name:		h.Curriculum.Name,
			Description: h.Curriculum.Description,
			Courses:     h.Curriculum.Courses,
		}); err != nil {
			fmt.Println(err)
		}
	default:
		// Check if the course exists and handle appropriately.
		if index, course := h.Curriculum.GetCourseByShortname(coursePath); course != nil {
			// Redirect to first lesson if just the course is accessed.
			if lessonPath == "" {
				http.Redirect(w, r, path.Join("/", coursePath, "1"), http.StatusTemporaryRedirect)
				return
			}
			// Get lessonIndex. Redirect to List on error.
			if lessonIndex, err = strconv.Atoi(lessonPath); err != nil {
				http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
				return
			}
			// Redirect to previous course if possible. If not, redirect to List.
			if lessonIndex <= 0 {
				if index > 0 {
					previousCourse := h.Curriculum.Courses[index-1]
					http.Redirect(w, r, path.Join("/", previousCourse.Shortname, strconv.Itoa(len(previousCourse.Lessons))), http.StatusTemporaryRedirect)
				} else {
					http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
				}
				return
			}
			// Redirect to next course if possible. If not, redirect to List.
			if lessonIndex > len(course.Lessons) {
				if index < len(h.Curriculum.Courses)-1 {
					nextCourse := h.Curriculum.Courses[index+1]
					http.Redirect(w, r, path.Join("/", nextCourse.Shortname, "1"), http.StatusTemporaryRedirect)
				} else {
					http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
				}
				return
			}
			// If we got to here that means we have a lesson ready.
			t := LessonTemplateData{
				CurriculumName: h.Curriculum.Name,
				LessonIndex:    lessonIndex,
				LessonsCount:   len(course.Lessons),
				Lesson:         &course.Lessons[lessonIndex-1],
			}
			if err = h.LessonTemplate.Execute(w, t); err != nil {
				fmt.Println(err)
			}
			return
		}
		// Otherwise redirect to the list if the course is invalid.
		http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
	}
}

func (h *CourseMaster) getCourseAndLessonStrings(s string) (course string, lesson string, err error) {
	parts := strings.Split(path.Clean("/"+s), "/")
	if len(parts) >= 2 {
		course = parts[1]
	}
	if len(parts) >= 3 {
		lesson = parts[2]
	}
	if len(parts) >= 4 {
		err = errors.New("Excessive length in passed string")
	}
	return
}

// ShiftPath is used to get the next segment in a path.
func (h *CourseMaster) ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)

	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}

	return p[1:i], p[i:]
}
