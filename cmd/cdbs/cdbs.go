package main

import (
	"bufio"
	"github.com/jessevdk/go-flags"
	"github.com/shirayu/cdbs"
	"io/ioutil"
	"log"
	"os"
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
	Help      bool   `short:"h" long:"help" description:"Show this help message"`
	Input     string `short:"i" long:"input" default:"-"`
	Output    string `short:"o" long:"output"`
	Separator string `short:"t" long:"separator" description:"Separator of keys and values" default:"	"`
	Single    bool   `long:"single" description:"Only output a single CDB file" default:"false"`
	Log       bool   `long:"log" description:"Enable logging" default:"false"`
}

func main() {
	opts := cmdOptions{}
	optparser := flags.NewParser(&opts, flags.Default)
	optparser.Name = "cdbs"
	optparser.Usage = "-i input -o output [OPTIONS]"
	_, err := optparser.Parse()
	if err != nil {
		for _, arg := range os.Args {
			if arg == "-h" {
				os.Exit(0)
			}
		}
		os.Exit(1)
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
	cdbs.Output(r, opts.Output, opts.Single, runes[0])
}
