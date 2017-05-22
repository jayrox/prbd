# prbd
Probes media for codec info and Analyzes for possible errors.


setup your go path  
run `go build`

then `prbd -dir /path/to/media/library`

will create a files files for each media file.

`filename-probe.json` will contain codec info in json formatting.  
`filename-analyze.txt` will contain non-error info. typically this will be empty.  
`filename-analyze-error.txt` will contain error info. hopefully, this will be empty.
