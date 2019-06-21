package main

import (
	"encoding/json"
	"fmt"
	"github.com/kettek/go-multipath"
	"path"
)

// Curriculum is the containing type for a set of courses.
type Curriculum struct {
	Name        string
	Description string
	CourseDirs  []string `json:"Courses"`
	Courses     []Course
}

// GetCourseByShortname returns a pointer to a course if it exists.
func (c *Curriculum) GetCourseByShortname(coursePath string) (index int, course *Course) {
	for i := range c.Courses {
		if c.Courses[i].Shortname == coursePath {
			course = &c.Courses[i]
			index = i
			return
		}
	}
	index = -1
	return
}

// CurriculumFromMultiPath attempts to unmarshal the json from the given file path into the curriculum.
func CurriculumFromMultiPath(multiPath multipath.Multipath, filePath string) (c Curriculum, err error) {
	var bytes []byte
	if bytes, err = multiPath.ReadFile(filePath); err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &c); err != nil {
		return
	}
	// Construct our Lesson Files.
	for _, courseDirname := range c.CourseDirs {
		var course Course
		if course, err = CourseFromDirMultiPath(multiPath, path.Join("curriculum", courseDirname)); err != nil {
			fmt.Print(err)
		}
		c.Courses = append(c.Courses, course)
	}
	return
}
