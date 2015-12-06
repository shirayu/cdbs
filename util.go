package cdbs

import (
	"bufio"
	"github.com/torbit/cdb"
	"io"
	"os"
	"sort"
	"strings"
)

//Cdbs handles several CDB files
type Cdbs struct {
	prefix string
	cdbs   []*cdb.Cdb
	names  []string
}

//NewCdbs returns Cdbs
func NewCdbs(prefix string) (*Cdbs, error) {
	self := new(Cdbs)
	self.prefix = prefix
	self.cdbs = []*cdb.Cdb{}
	self.names = []string{}

	ifname := prefix + ".keymap"
	inf, err := os.Open(ifname)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(inf)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.SplitN(line, " ", 2)
		if len(items) != 2 {
			continue
		}
		key := items[0]
		num := items[1]
		mycdb, err := cdb.Open(prefix + "." + num + ".cdb")
		if err != nil {
			return nil, err
		}
		//         defer mycdb.Close()
		self.cdbs = append(self.cdbs, mycdb)
		self.names = append(self.names, key)
	}
	return self, nil
}

//Get returns a value by binary search
func (cs *Cdbs) Get(key string) ([]byte, error) {
	index := sort.SearchStrings(cs.names, key)
	if index >= len(cs.names) {
		index--
	} else if cs.names[index] != key {
		index--
	}
	if index < 0 {
		index = 0
	}

	c := cs.cdbs[index]
	data, err := c.Bytes([]byte(key))
	return data, err
}

//BruteGet returns a value by linear search
func (cs *Cdbs) BruteGet(key string) ([]byte, error) {
	for _, c := range cs.cdbs {
		data, err := c.Bytes([]byte(key))
		if err == nil {
			return data, err
		}
	}
	return nil, io.EOF
}
