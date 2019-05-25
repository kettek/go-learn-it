package main

import (
	"fmt"
	"net/http"
	"flag"
	"os"
	"path"
	"path/filepath"
	"github.com/kettek/go-multipath"
)

func main() {
	var err error
	var courseMaster CourseMaster
	var staticMultiPath multipath.Multipath

	// Get installation directory
	binaryPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Print(err)
		return
	}
	sharePath := path.Join(filepath.Dir(filepath.Dir(binaryPath)), "share", "go-learn-it")

	sharepathPtr := flag.String("sharepath", sharePath, "Path to the built-in go-learn-it static and template data")
	datapathPtr := flag.String("datapath", "", "Path to external go-learn-it static and template data")
	curriculumPtr := flag.String("curriculum", "curriculum.json", "Path to the curriculum data file located in either basepath or datapath")

	flag.Parse()

	if len(*sharepathPtr) > 0 {
		courseMaster.multiPath.AddPath(*sharepathPtr, multipath.FirstPriority)
		staticMultiPath.AddPath(path.Join(*sharepathPtr, "static"), multipath.FirstPriority)
	}
	if len(*datapathPtr) > 0 {
		courseMaster.multiPath.AddPath(*datapathPtr, multipath.FirstPriority)
		staticMultiPath.AddPath(path.Join(*datapathPtr, "static"), multipath.FirstPriority)
	}

	//
	fmt.Println("Starting go-learn-it")

	// Setup our HTTP server.
	fs := multipathFileSystem{http.Dir(""), &staticMultiPath}
	if courseMaster.Curriculum, err = CurriculumFromMultiPath(courseMaster.multiPath, *curriculumPtr); err != nil {
		fmt.Print(err)
		return
	}
	if err := courseMaster.buildTemplates(); err != nil {
		fmt.Print(err)
		return
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))
	http.HandleFunc("/api/", courseMaster.handleAPI)
	http.HandleFunc("/", courseMaster.handleHTTP)
	http.ListenAndServe(":8888", nil)
}

type multipathFileSystem struct {
	http.FileSystem
	mpath *multipath.Multipath
}

func (fs multipathFileSystem) Open(name string) (http.File, error) {
	file, err := fs.mpath.Open(name)
	return file, err
}