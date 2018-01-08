# prbd
Probes media for codec info and Analyzes for possible errors.

you will need `ffmpeg` and `ffprobe`. you can download these from here: http://ffmpeg.zeranoe.com/builds/  

# Build it yourself
setup your go path  
run `go build`

# Or download current build
https://github.com/jayrox/prbd/tree/master/build 

# Run
`prbd -dir /path/to/media/library` # Will recursively scan specified directory

`prbd -r` # Will scan the path in the `radarr_moviefile_path` env variable, for use as a Connected Script in Radarr

`prbd -s` # Will scan the path in the `sonarr_episodefile_path` env variable, for use as a Connected Script in Sonarr

add a `-log /path/to/reports` and prbd will output the reports to a log directory. useful for building reports and quickly finding media with specific codecs or codec levels.

----
# Output
Prbd will create a few files for each media file found.

`filename-probe.json` will contain codec info in json formatting.  
`filename-analyze.txt` will contain non-error info. typically this will be empty.  
`filename-analyze-error.txt` will be created if errors found.
