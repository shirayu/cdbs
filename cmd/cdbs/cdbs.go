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
	Help   bool   `short:"h" long:"help" description:"Show this help message"`
	Input  string `short:"i" long:"input" default:"-"`
	Output string `short:"o" long:"output"`
	Single bool   `long:"single" description:"Only output a single CDB file" default:"false"`
	Log    bool   `long:"log" description:"Enable logging" default:"false"`
}

func main() {
	opts := cmdOptions{}
	optparser := flags.NewParser(&opts, flags.Default)
	optparser.Name = "cdbs"
	optparser.Usage = "-i input -o output [OPTIONS]"
	optparser.Parse()

	//show help
	if len(os.Args) == 1 {
		optparser.WriteHelp(os.Stdout)
		os.Exit(0)
	}
	for _, arg := range os.Args {
		if arg == "-h" {
			os.Exit(0)
		}
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
	if opts.Single {
		cdbs.MakeCDB(r, opts.Output)
	} else {
		cdbs.Output(r, opts.Output)
	}
}
