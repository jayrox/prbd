package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	flagdebug = flag.Bool("d", true, "show debug output")
	flagdir   = flag.String("dir", "cwd", "directory to scan. default is current working directory (cwd)")
	pf        prbdFlags
	mediaext  = []string{".avi", ".mkv", ".mp4"}
)

func main() {
	flag.Parse()
	// Print the logo :P
	printLogo()

	// Root folder to scan
	fpSAbs, _ := filepath.Abs(flagString(flagdir))
	pf.Dir = fpSAbs
	if flagString(flagdir) == "cwd" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		check(err)
		pf.Dir = dir
	}

	fmt.Printf("Scanning directory: %s\n", pf.Dir)
	fmt.Println("_____________________")

	i := folderWalk(pf.Dir)
	if i < 1 {
		fmt.Println("No media found.")
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
	cmd := exec.Command("./ffprobe", "-v", "error", "-show_format", "-show_streams", filename, "-print_format", "json")

	r, err := cmd.StdoutPipe()
	check(err)

	r2, err := cmd.StderrPipe()
	check(err)

	err = cmd.Start()
	check(err)

	fname := filename + "-probe.json"
	body, err := ioutil.ReadAll(r)
	check(err)
	writeToFile(fname, body)

	fname2 := filename + "-probe-error.json"
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
	cmd := exec.Command("./ffmpeg", "-v", "error", "-i", filename, "-f", "null", "-")

	r, err := cmd.StdoutPipe()
	check(err)

	r2, err := cmd.StderrPipe()
	check(err)

	err = cmd.Start()
	check(err)

	fname := filename + "-analyze.txt"
	body, err := ioutil.ReadAll(r)
	check(err)
	writeToFile(fname, body)

	fname2 := filename + "-analyze-error.txt"
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
	f, err := os.Create(path)
	check(err)

	defer f.Close()

	n, err := f.Write(data)
	check(err)
	fmt.Printf("wrote %d bytes\n", n)

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
	if *flagdebug {
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
