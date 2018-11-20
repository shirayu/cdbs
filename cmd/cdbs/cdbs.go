package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/shirayu/cdbs"
)

func getInputFile(ifname string) (inf *os.File, err error) {

	if ifname == "-" {
		inf = os.Stdin
	} else {
		inf, err = os.Open(ifname)
		if err != nil {
			return nil, err
		}
	}
	return inf, nil
}

type cmdOptions struct {
	Input     string `short:"i" long:"input" default:"-"`
	Output    string `short:"o" long:"output"`
	Separator string `short:"t" long:"separator" description:"Separator of keys and values" default:"	"`
	Compress  bool   `short:"z" long:"compress" description:"Compress values in gzip format"`
	Single    bool   `long:"single" description:"Only output a single CDB file"`
	Log       bool   `long:"log" description:"Enable logging"`
	Limit     int    `long:"lmit" description:"Limit size of one CDB file" default:"4089446400"`

	Version bool `short:"v" long:"version" description:"Show version"`
}

//Version of this program
var Version = "Unknown version"

//VersionDate is the commit date of this version
var VersionDate = ""

func main() {
	opts := cmdOptions{}
	optparser := flags.NewParser(&opts, flags.Default)
	optparser.Name = "cdbs"
	optparser.Usage = "-i input -o output [OPTIONS]"
	_, err := optparser.Parse()

	if len(os.Args) == 1 {
		optparser.WriteHelp(os.Stdout)
		os.Exit(0)
	} else if err != nil {
		if e, ok := err.(*flags.Error); !ok {
			log.Fatalf("Expected flags.Error, but got %T", err)
		} else if e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("%s\n", err)
	} else if opts.Version {
		if len(VersionDate) != 0 {
			fmt.Fprintf(os.Stderr, "%s: %s on %s\n", optparser.Name, Version, VersionDate)
		} else {
			fmt.Fprintf(os.Stderr, "%s: %s\n", optparser.Name, Version)
		}
		os.Exit(0)
	}

	runes := []rune(opts.Separator)
	if len(runes) != 1 {
		log.Printf("The length of separator is not 1")
		os.Exit(1)
	}

	if opts.Input == "" || opts.Output == "" {
		log.Printf("Input or output is not given")
		os.Exit(1)
	}

	if opts.Log == false {
		log.SetOutput(ioutil.Discard)
	}

	//open
	inf, err := getInputFile(opts.Input)
	defer inf.Close()
	if err != nil {
		log.Printf("Error when opening %s", opts.Input)
		os.Exit(1)
	}

	//operate
	r := bufio.NewReader(inf)
	cdbs.Output(r, opts.Output, opts.Single, runes[0], opts.Compress, opts.Limit)
}
