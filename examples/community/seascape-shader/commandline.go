package main

/* 
	This simply parses the command line arguments using the default golang
	package called 'flag'. This can be used as a simple template to parse
	command line arguments in other programs.
*/

import "log"
import "os"
import "flag"
import "fmt"

var (
	version       string
	race          bool
	debug         = os.Getenv("BUILDDEBUG") != ""
	filename      string
	width         int
	height        int
	timeout       = "120s"
	uDrift        float32
)

var customUsage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Println("")
	fmt.Println("EXAMPLE:")
	fmt.Println("")
	fmt.Println("seascape-shader --width 640 --height 480 --filename shaders/planetfall.glsl")
	fmt.Println("")
}

func parseFlags() {
	flag.StringVar  (&version,       "version",       "v0.1",                     "Set compiled in version string")
	flag.StringVar  (&filename,      "filename",      "shaders/seascape.glsl",    "path to GLSL file")
	flag.IntVar     (&width,         "width",         1024,                       "Width of the OpenGL Window")
	flag.IntVar     (&height,        "height",        768,                        "Height of the OpenGL Window")
	var tmp float64
	flag.Float64Var (&tmp,           "drift",         0.01,                       "Speed of the gradual camera drift")
	flag.BoolVar    (&race,          "race",          race,                       "Use race detector")

	// Set the output if something fails to stdout rather than stderr
	flag.CommandLine.SetOutput(os.Stdout)
	// flag.SetOutput(os.Stdout)

	flag.Usage = customUsage
	flag.Parse()

	if flag.Parsed() {
		log.Println("Parsed() worked. width=",width)
	} else {
		log.Println("Parsed() failed. width=",width)
	}

//	if err := flag.Parse(); err != nil {
//		log.Println("Example:",width)
//	}
		

	uDrift = float32(tmp)
	log.Println("width=",width)
	log.Println("height=",height)
	log.Println("uDrift=",uDrift)
}
