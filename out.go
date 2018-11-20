package cdbs

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/torbit/cdb"
)

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCdbName(outprefix string, numDB int, single bool) string {
	if single {
		if numDB != 0 {
			exitOnErr(errors.New("Multiple output in single output mode"))
		}
		return fmt.Sprintf("%s.cdb", outprefix)
	}
	return fmt.Sprintf("%s.%d.cdb", outprefix, numDB)
}

//MakeCDB creates single CDB file with given buffer
func MakeCDB(r *bufio.Reader, outname string) {
	log.Printf("Making %s", outname)
	outdb, err := os.OpenFile(outname, os.O_RDWR|os.O_CREATE, 0644)
	exitOnErr(err)

	exitOnErr(cdb.Make(outdb, r))
	exitOnErr(outdb.Sync())
	exitOnErr(outdb.Close())
	log.Printf("done")
}

//Output crests CDB files
func Output(r *bufio.Reader, outpath string, single bool, separator rune, compress bool, limit int) {
	var err error
	var buf bytes.Buffer
	buf.Grow(4 * (1024 * 1024 * 1024)) //get 4GB
	firstKeys := []string{}
	const minimumCdbSize = 2046
	cdbSize := minimumCdbSize
	numDB := 0
	line, err := r.ReadBytes('\n') //first
	for err != io.EOF {
		if err != io.EOF && err != nil {
			log.Fatal(err)
		}
		delmPos := bytes.IndexRune(line, separator)
		//         (line[:len(line)-1], "\t", 2)
		if delmPos == -1 {
			log.Printf("skip an invalid line -> %s", line)
			line, err = r.ReadBytes('\n') //next
			continue
		}

		//Get value expression
		var valByte []byte
		var valSize int
		if compress {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			if _, err := gz.Write(line[delmPos+1 : len(line)-1]); err != nil {
				log.Fatal(err)
			}
			if err := gz.Flush(); err != nil {
				log.Fatal(err)
			}
			if err := gz.Close(); err != nil {
				log.Fatal(err)
			}
			valByte = b.Bytes()
			valSize = len(valByte)
		} else {
			valByte = line[delmPos+1 : len(line)-1]
			valSize = len(line) - delmPos - 2
		}
		// 24 bytes per record, plus the space for keys and data
		newEntrySize := delmPos + valSize + 24

		//if the CDB size will exceed the given file size, make DB before adding the new line
		if cdbSize+newEntrySize > limit {
			r := bufio.NewReader(&buf)
			buf.WriteString("\n")
			outname := getCdbName(outpath, numDB, single)
			MakeCDB(r, outname)
			numDB++

			//clear
			cdbSize = minimumCdbSize
			buf.Reset()
			//             debug.FreeOSMemory()
		}

		if cdbSize == minimumCdbSize {
			key := string(line[:delmPos])
			firstKeys = append(firstKeys, key)
		}
		cdbSize += newEntrySize

		headLine := fmt.Sprintf("+%d,%d:", delmPos, valSize)
		buf.WriteString(headLine)
		keyByte := line[:delmPos]
		buf.Write(keyByte)
		buf.WriteRune('-')
		buf.WriteRune('>')
		buf.Write(valByte)
		buf.WriteRune('\n')

		line, err = r.ReadBytes('\n') //next
	}

	rbuf := bufio.NewReader(&buf)
	buf.WriteString("\n")
	outname := getCdbName(outpath, numDB, single)
	MakeCDB(rbuf, outname)

	//output keymap
	if !single {
		outf, err := os.Create(outpath + ".keymap")
		defer outf.Close()
		exitOnErr(err)
		w := bufio.NewWriter(outf)
		defer w.Flush()
		for idx, key := range firstKeys {
			w.WriteString(key)
			w.WriteString(" ")
			w.WriteString(fmt.Sprintf("%d", idx))
			w.WriteString("\n")
		}
	}

}
