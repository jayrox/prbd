package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	flagDebug  = flag.Bool("d", true, "show debug output")
	flagDir    = flag.String("dir", "cwd", "directory to scan. default is current working directory (cwd)")
	flagLog    = flag.String("log", "cwd", "directory to write the log. default is current working directory (cwd)")
	flagSonarr = flag.String("s", "false", "used when called from sonarr. gets the file name from the enviroment variables.")
	flagRadarr = flag.String("r", "false", "used when called from radarr. gets the file name from the enviroment variables.")
	pf         prbdFlags
	mediaext   = []string{".avi", ".divx", ".m4v", ".mkv", ".mp4"}
)

func main() {
	flag.Parse()
	// Print the logo :P
	printLogo()

	// Root folder to scan
	fpSAbs, _ := filepath.Abs(flagString(flagDir))
	pf.Dir = fpSAbs
	if flagString(flagDir) == "cwd" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		check(err)
		pf.Dir = dir
	}

	// Root folder to write logs
	fpLAbs, _ := filepath.Abs(flagString(flagLog))
	pf.Log = fpLAbs
	if flagString(flagLog) == "cwd" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		check(err)
		pf.Log = dir
	}

	if flagString(flagSonarr) == "true" {
		pf.File = os.Getenv("sonarr_episodefile_path")
		pf.Dir = ""
	}

	if flagString(flagRadarr) == "true" {
		pf.File = os.Getenv("radarr_moviefile_path")
		pf.Dir = ""
	}

	if pf.File == "" {
		fmt.Printf("Scanning directory: %s\n", pf.Dir)
		fmt.Println("_____________________")

		i := folderWalk(pf.Dir)
		if i < 1 {
			fmt.Println("No media found.")
		}
	} else {
		if _, err := os.Stat(pf.File); os.IsNotExist(err) {
			fmt.Println(pf.File + "- does not exist.")
			pf.File = strings.Replace(pf.File, ".avi", ".mp4", -1)
			pf.File = strings.Replace(pf.File, ".divx", ".mp4", -1)
			pf.File = strings.Replace(pf.File, ".m4v", ".mp4", -1)
			pf.File = strings.Replace(pf.File, ".mkv", ".mp4", -1)
			fmt.Println("Trying: " + pf.File)
		}

		ok := false

		ok = mediaProbe(pf.File)
		if ok == false {
			fmt.Println("Probe Failed - skipping media")
		}

		ok = mediaTest(pf.File)
		if ok == false {
			fmt.Println("Test Failed - skipping media")
		}

	}
}

func flagString(fs *string) string {
	return fmt.Sprint(*fs)
}

func flagInt(fi *int64) int64 {
	return int64(*fi)
}

func flagBool(fb *bool) bool {
	return bool(*fb)
}

func folderWalk(path string) (i int64) {
	i = 0
	var err = filepath.Walk(path, func(path string, _ os.FileInfo, _ error) error {
		for _, x := range mediaext {
			if filepath.Ext(path) == x {
				var ok bool

				t := time.Now()
				fmt.Println(t.Format("2006-01-02:15:04:05"))

				ok = mediaProbe(path)
				if ok == false {
					fmt.Println("Probe Failed - skipping media")
					fmt.Println("_________")
					continue
				}

				ok = mediaTest(path)
				if ok == false {
					fmt.Println("Test Failed - skipping media")
					fmt.Println("_________")
					continue
				}
				fmt.Println("_________ _________ _________")
			}
		}
		return nil
	})
	if err != nil {
		printDebug("er: %+v\n", err)
	}
	return
}

func mediaProbe(path string) (ok bool) {
	ok = true
	printDebug("Probing: %s\n", path)
	err := probe(path)
	check(err)

	return
}

func mediaTest(path string) (ok bool) {
	ok = true
	printDebug("Testing: %s\n", path)
	err := analyze(path)
	check(err)

	return
}

func probe(filename string) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)
	ffprobe := filepath.Join(dir, "ffprobe")
	cmd := exec.Command(ffprobe, "-v", "error", "-show_format", "-show_streams", filename, "-print_format", "json")

	r, err := cmd.StdoutPipe()
	check(err)

	r2, err := cmd.StderrPipe()
	check(err)

	err = cmd.Start()
	check(err)

	fname := strings.Trim(filename, " ") + "-probe.json"
	body, err := ioutil.ReadAll(r)
	check(err)
	writeToFile(fname, body)

	fname2 := strings.Trim(filename, " ") + "-probe-error.json"
	body2, err := ioutil.ReadAll(r2)
	check(err)
	if len(body2) > 0 {
		writeToFile(fname2, body2)
	}

	err = cmd.Wait()
	check(err)

	return nil
}

func analyze(filename string) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)
	ffmpeg := filepath.Join(dir, "ffmpeg")
	cmd := exec.Command(ffmpeg, "-v", "error", "-i", filename, "-f", "null", "-")

	r, err := cmd.StdoutPipe()
	check(err)

	r2, err := cmd.StderrPipe()
	check(err)

	err = cmd.Start()
	check(err)

	fname := strings.Trim(filename, " ") + "-analyze.txt"
	body, err := ioutil.ReadAll(r)
	check(err)
	writeToFile(fname, body)

	fname2 := strings.Trim(filename, " ") + "-analyze-error.txt"
	body2, err := ioutil.ReadAll(r2)
	check(err)
	if len(body2) > 0 {
		writeToFile(fname2, body2)
	}

	err = cmd.Wait()
	check(err)

	return nil
}

func writeToFile(path string, data []byte) {
	if filepath.Dir(path) == filepath.Dir(pf.Log) {
		// do natta thing
	} else {
		path = filepath.Join(pf.Log, filepath.Base(path))
	}
	fmt.Println(path)

	fileName := filepath.Base(path)
	fileName = strings.Replace(fileName, "\"", "", -1)
	fileName = strings.Replace(fileName, ":", "", -1)
	fileName = strings.Replace(fileName, "*", "", -1)
	fileName = strings.Replace(fileName, "?", "", -1)
	fileName = strings.Replace(fileName, "<", "", -1)
	fileName = strings.Replace(fileName, ">", "", -1)
	fileName = strings.Replace(fileName, "|", "", -1)
	fileName = strings.Trim(fileName, " ")

	filePath := filepath.Dir(path)
	path = filepath.Join(filePath, fileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	check(err)

	defer f.Close()

	_, err = f.Write(data) // _ was n
	check(err)
	//fmt.Printf("wrote %d bytes\n", n)

	f.Sync()
}

// Check err
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Only print debug output if the debug flag is true
func printDebug(format string, vars ...interface{}) {
	if *flagDebug {
		if vars[0] == nil {
			fmt.Println(format)
			return
		}
		fmt.Printf(format, vars...)
	}
}

// Hold flag data
type prbdFlags struct {
	Dir   string
	Debug bool
	File  string
	Log   string
}

// Print the logo, obviously
func printLogo() {
	fmt.Println("██████╗ ██████╗ ██████╗ ██████╗")
	fmt.Println("██╔══██╗██╔══██╗██╔══██╗██╔══██╗")
	fmt.Println("██████╔╝██████╔╝██████╔╝██║  ██║")
	fmt.Println("██╔═══╝ ██╔══██╗██╔══██╗██║  ██║")
	fmt.Println("██║     ██║  ██║██████╔╝██████╔╝")
	fmt.Println("╚═╝     ╚═╝  ╚═╝╚═════╝ ╚═════╝ probed")
	fmt.Println("")
}
