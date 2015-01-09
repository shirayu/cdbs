
# cdbs


## What's this

- This tool converts the input key and value pairs into several CDB files
- Input SHOULD be sorted by the key, otherwise you can not look up by using several cdbs

## Usage
```
Usage of cdbs:
  -i, -input :  Input file name. - or no designation means STDIN.
  -o, -output: Output file name suffix.
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

