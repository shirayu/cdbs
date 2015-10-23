package cdbs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/torbit/cdb"
	"io"
	"log"
	"os"
)

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func get_cdb_name(outprefix string, num_db int, single bool) string {
	if single {
		if num_db != 0 {
			exitOnErr(errors.New("Multiple output in single output mode"))
		}
		return fmt.Sprintf("%s.cdb", outprefix)
	}
	return fmt.Sprintf("%s.%d.cdb", outprefix, num_db)
}

func MakeCDB(r *bufio.Reader, outname string) {
	log.Printf("Making %s", outname)
	outdb, err := os.OpenFile(outname, os.O_RDWR|os.O_CREATE, 0644)
	exitOnErr(err)

	exitOnErr(cdb.Make(outdb, r))
	exitOnErr(outdb.Sync())
	exitOnErr(outdb.Close())
	log.Printf("done")
}
func Output(r *bufio.Reader, outpath string, single bool, separator rune) {
	var err error = nil
	var buf bytes.Buffer
	buf.Grow(4 * (1024 * 1024 * 1024)) //get 4GB
	first_keys := []string{}
	buf_size := 0
	num_db := 0
	line, err := r.ReadBytes('\n') //first
	for err != io.EOF {
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		delm_pos := bytes.IndexRune(line, separator)
		//         (line[:len(line)-1], "\t", 2)
		if delm_pos == -1 {
			log.Printf("skip an invalid line -> %s", line)
			line, err = r.ReadBytes('\n') //next
			continue
		}

		//cdb line format is "+<Size-of-key>,<Size-of-val>:<key>-><val>\n" like "+3,4:tom->baby\n"
		new_line_size := len(line) + 20 //+ , : =>

		//if the buffer size will exceed 3.5GB, make DB before adding the new line
		if buf_size+new_line_size > 3.5*(1024*1024*1024) {
			r := bufio.NewReader(&buf)
			buf.WriteString("\n")
			outname := get_cdb_name(outpath, num_db, single)
			MakeCDB(r, outname)
			num_db++

			//clear
			buf_size = 0
			buf.Reset()
			//             debug.FreeOSMemory()
		}

		if buf_size == 0 {
			key := string(line[:delm_pos])
			first_keys = append(first_keys, key)
		}
		buf_size += new_line_size

		val_size := len(line) - delm_pos - 2
		head_line := fmt.Sprintf("+%d,%d:", delm_pos, val_size)
		buf.WriteString(head_line)
		key_byte := line[:delm_pos]
		buf.Write(key_byte)
		buf.WriteRune('-')
		buf.WriteRune('>')
		val_and_ln_byte := line[delm_pos+1:]
		buf.Write(val_and_ln_byte)

		line, err = r.ReadBytes('\n') //next
	}

	rbuf := bufio.NewReader(&buf)
	buf.WriteString("\n")
	outname := get_cdb_name(outpath, num_db, single)
	MakeCDB(rbuf, outname)

	//output keymap
	if !single {
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

}
