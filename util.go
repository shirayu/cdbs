package cdbs

import (
	"bufio"
	"github.com/torbit/cdb"
	"io"
	"os"
	"sort"
	"strings"
)

type Cdbs struct {
	prefix string
	cdbs   []*cdb.Cdb
	names  []string
}

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

func (self *Cdbs) Get(key string) ([]byte, error) {
	index := sort.SearchStrings(self.names, key)
	if index >= len(self.names) {
		index--
	} else if self.names[index] != key {
		index--
	}
	if index < 0 {
		index = 0
	}

	c := self.cdbs[index]
	data, err := c.Bytes([]byte(key))
	return data, err
}

func (self *Cdbs) BruteGet(key string) ([]byte, error) {
	for _, c := range self.cdbs {
		data, err := c.Bytes([]byte(key))
		if err == nil {
			return data, err
		}
	}
	return nil, io.EOF
}
