package main

import (
	"fmt"
	"net/http"
	"flag"
	"strconv"
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

	// Default run-time variables
	address := ""
	port := 8888
	sharepath := path.Join(filepath.Dir(filepath.Dir(binaryPath)), "share", "go-learn-it")
	datapath := ""
	curriculum := "curriculum.json"

	// Merge environment with run-time
	if addressEnv, exists := os.LookupEnv("GLI_ADDRESS"); exists {
		address = addressEnv
	}
	if portEnv, exists := os.LookupEnv("GLI_PORT"); exists {
		port, _ = strconv.Atoi(portEnv)
	}
	if sharepathEnv, exists := os.LookupEnv("GLI_SHAREPATH"); exists {
		sharepath = sharepathEnv
	}
	if datapathEnv, exists := os.LookupEnv("GLI_DATAPATH"); exists {
		datapath = datapathEnv
	}
	if curriculumEnv, exists := os.LookupEnv("GLI_CURRICULUM"); exists {
		curriculum = curriculumEnv
	}

	// Override run-time/environment with flags
	flag.StringVar(&address, "address", address, "HTTP listen address - GLI_ADDRESS")
	flag.IntVar(&port, "port", port, "HTTP listen port - GLI_PORT")
	flag.StringVar(&sharepath, "sharepath", sharepath, "Path to the built-in go-learn-it static and template data - GLI_SHAREPATH")
	flag.StringVar(&datapath, "datapath", datapath, "Path to external go-learn-it static and template data - GLI_DATAPATH")
	flag.StringVar(&curriculum, "curriculum", curriculum, "Filename of the curriculum data file located in either basepath or datapath - GLI_CURRICULUM")

	flag.Parse()

	// Add sharepath to the multipath if it is set.
	if len(sharepath) > 0 {
		courseMaster.multiPath.AddPath(sharepath, multipath.FirstPriority)
		staticMultiPath.AddPath(path.Join(sharepath, "static"), multipath.FirstPriority)
	}
	// Add datapath to the multipath if it is set.
	if len(datapath) > 0 {
		courseMaster.multiPath.AddPath(datapath, multipath.FirstPriority)
		staticMultiPath.AddPath(path.Join(datapath, "static"), multipath.FirstPriority)
	}

	// Okay, we're good to go!
	fmt.Println("Starting go-learn-it")

	// Setup the CourseMaster.
	if courseMaster.Curriculum, err = CurriculumFromMultiPath(courseMaster.multiPath, curriculum); err != nil {
		fmt.Print(err)
		return
	}
	if err := courseMaster.buildTemplates(); err != nil {
		fmt.Print(err)
		return
	}
	// Setup our HTTP server.
	fs := multipathFileSystem{http.Dir(""), &staticMultiPath}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(fs)))
	http.HandleFunc("/api/", courseMaster.handleAPI)
	http.HandleFunc("/", courseMaster.handleHTTP)

	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
	fmt.Printf("Now listening on %s:%d\n", address, port)
}

type multipathFileSystem struct {
	http.FileSystem
	mpath *multipath.Multipath
}

// Open provides the multipath opening for our multipath HTTP FileSystem
func (fs multipathFileSystem) Open(name string) (http.File, error) {
	file, err := fs.mpath.Open(name)
	return file, err
}