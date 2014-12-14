package cdbs

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jbarham/go-cdb"
	"io"
	"log"
	"os"
	"strings"
)

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func get_cdb_name(outprefix string, num_db int) string {
	return fmt.Sprintf("%s.%d.cdb", outprefix, num_db)
}

func makeCDB(outprefix string, num_db int, r *bufio.Reader) {
	outname := get_cdb_name(outprefix, num_db)
	log.Printf("Making %s", outname)
	outdb, err := os.OpenFile(outname, os.O_RDWR|os.O_CREATE, 0644)
	exitOnErr(err)

	exitOnErr(cdb.Make(outdb, r))
	exitOnErr(outdb.Sync())
	exitOnErr(outdb.Close())
	log.Printf("done")
}
func Output(r *bufio.Reader, outpath string) {
	var err error = nil
	var buf bytes.Buffer
	first_keys := []string{}
	buf_size := 0
	num_db := 0
	line, err := r.ReadString('\n') //first
	for err != io.EOF {
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		items := strings.SplitN(line[:len(line)-1], "\t", 2)
		if len(items) != 2 {
			log.Printf("skip an invalid line -> %s", line)
			line, err = r.ReadString('\n') //next
			continue
		}

		key := items[0]
		val := items[1]
		cdb_line := fmt.Sprintf("+%d,%d:%s->%s\n", len(key), len(val), key, val)
		cdb_line_byte := []byte(cdb_line)

		//if the buffer size will exceed 3.5GB, make DB before adding the new line
		if buf_size+len(cdb_line_byte) > 3.5*(1024*1024*1024) {
			r := bufio.NewReader(&buf)
			buf.WriteString("\n")
			makeCDB(outpath, num_db, r)
			num_db++

			//clear
			buf.Reset()
			buf_size = 0
			buf = bytes.Buffer{}
		}

		if buf_size == 0 {
			first_keys = append(first_keys, key)
		}
		buf_size += len(cdb_line_byte)
		buf.Write(cdb_line_byte)

		line, err = r.ReadString('\n') //next
	}

	rbuf := bufio.NewReader(&buf)
	buf.WriteString("\n")
	makeCDB(outpath, num_db, rbuf)

	//output keymap
	outf, err := os.Create(outpath + ".keymap")
	defer outf.Close()
	exitOnErr(err)
	w := bufio.NewWriter(outf)
	defer w.Flush()
	for idx, key := range first_keys {
		w.WriteString(key)
		w.WriteString(" ")
		w.WriteString(fmt.Sprintf("%d", idx))
		w.WriteString("\n")
	}

}
