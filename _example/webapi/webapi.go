package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/shirayu/cdbs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type cmdOptions struct {
	Help  bool   `short:"h" long:"help" description:"Show this help message"`
	Input string `short:"i" long:"input" required:"true"`
	Port  int    `short:"p" long:"port" default:"8000" description:"Port number to serve"`
	NoLog bool   `long:"nolog" description:"Disable logging" default:"false"`
}

func accessLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s\t%s\t%s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

type Server struct {
	cdbs *cdbs.Cdbs
}

func NewServer(cdbs *cdbs.Cdbs) *Server {
	self := new(Server)
	self.cdbs = cdbs
	return self
}
func (self *Server) Handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) < 2 {
		fmt.Fprintf(w, "Format error")
		return
	}
	id := strings.Join(paths[1:], "/")
	for k, _ := range r.Form {
		id = id + "?" + k
	}
	log.Printf("\t->[%s]", id)

	data, err := self.cdbs.BruteGet(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{\"error\": \"%s\"}", err)
	} else {
		w.Write(data)
	}
}

func operation(opts *cmdOptions) {
	mycdbs, err := cdbs.NewCdbs(opts.Input)
	if err != nil {
		log.Fatal(err)
	}
	server := NewServer(mycdbs)
	http.HandleFunc("/", server.Handler)

	//serve
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	address := fmt.Sprintf(":%d", opts.Port)
	log.Printf("Started serving at [%s] with [%s]", address, opts.Input)
	err = http.ListenAndServe(address, accessLog(http.DefaultServeMux))
	log.Printf("Error: %s", err)

}

func main() {
	opts := cmdOptions{}
	optparser := flags.NewParser(&opts, flags.Default)
	optparser.Name = ""
	optparser.Usage = "-i input [OPTIONS]"
	_, err := optparser.Parse()

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
	if err != nil {
		os.Exit(1)
	}

	if opts.NoLog {
		log.SetOutput(ioutil.Discard)
	}

	operation(&opts)
}
