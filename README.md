
# cdbs

[![Build Status](https://travis-ci.org/shirayu/cdbs.svg?branch=master)](https://travis-ci.org/shirayu/cdbs)
[![Coverage Status](https://coveralls.io/repos/shirayu/cdbs/badge.png)](https://coveralls.io/r/shirayu/cdbs)
[![LGPLv3](https://img.shields.io/badge/license-LGPLv3-blue.svg)](LGPLv3)
[![BSD](https://img.shields.io/badge/license-BSD-blue.svg)](BSD)


## What's this

- This tool converts the input key and value pairs into several CDB files
- Input SHOULD be sorted by the key, otherwise you can not look up by using several cdbs with ``Get(key string)``
    - For unsorted keys, ``BruteGet(key string)`` can be used. Because this searches all CDB files, the efficiency will be down.

## Usage
```
Usage of cdbs:
  -i, -input :  Input file name. - or no designation means STDIN.
  -o, -output: Output file name suffix.
  -t, --separator: Separator of keys and values (deault: "\t")
  -z, --compress:   Compress values in gzip format (deault: false)
  --single:     Only output a single CDB file (deault: false)
  --log:        Enable logging (deault: false)
```

## INSTALL

```
go get github.com/shirayu/cdbs/cmd/cdbs
```

- To install library
```
go get github.com/shirayu/cdbs
```

## Acknowledgement

I developed this program as a part of the research project 
["Establishment of Knowledge-Intensive Structural Natural Language Processing and Construction of Knowledge Infrastructure"](http://nlp.ist.i.kyoto-u.ac.jp/CREST/?en)
in [Kyoto University](http://www.kyoto-u.ac.jp/en)
supported by [CREST, JST](http://www.jst.go.jp/kisoken/crest/en/).


## Licence

- (c) Yuta Hayashibe 2014
- Released under any of the following licences
    - Lesser GNU General Public License 3.0 (see the file LGPL)
    - New BSD License (3-clause BSD License) (see the file BSD)

