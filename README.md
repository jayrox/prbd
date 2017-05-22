# prbd
Probes media for codec info and Analyzes for possible errors.

you will need `ffmpeg` and `ffprobe`. you can download these from here: http://ffmpeg.zeranoe.com/builds/  

# Build it yourself
setup your go path  
run `go build`

# Or download current build
https://github.com/jayrox/prbd/tree/master/build 

# Run
`prbd -dir /path/to/media/library`

----
# Output
Prbd will create a few files for each media file found.

`filename-probe.json` will contain codec info in json formatting.  
`filename-analyze.txt` will contain non-error info. typically this will be empty.  
`filename-analyze-error.txt` will be created if errors found.
