package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// WaitDuration Seconds to wait between cycles
const WaitDuration = 3 * time.Second

func main() {

	// Initial directory scan
	oldInfo, err := scanDirectory("/data")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Old (%d total)\n", len(oldInfo))
	for filename, info := range oldInfo {
		fmt.Printf("file: %s, modified: %s, size: %d\n", filename, info.ModTime(), info.Size())
	}
	fmt.Print("\n")

	// Wait some amount of time before secondary scan
	time.Sleep(WaitDuration)

	newInfo, err := scanDirectory("/data")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New (%d total)\n", len(newInfo))
	for filename, info := range newInfo {
		fmt.Printf("file: %s, modified: %s, size: %d\n", filename, info.ModTime(), info.Size())
	}
	fmt.Print("\n")

	unchanged, changed := compareScans(oldInfo, newInfo)

	fmt.Printf("Unchanged (%d total)\n", len(unchanged))
	fmt.Printf("%#v\n", unchanged)
	for filename, info := range unchanged {
		fmt.Printf("file: %s, modified: %s, size: %d\n", filename, info.ModTime(), info.Size())
	}
	fmt.Print("\n")

	fmt.Printf("Changed (%d total)\n", len(changed))
	fmt.Printf("%#v\n", changed)
	for filename, info := range changed {
		fmt.Printf("file: %s, modified: %s, size: %d\n", filename, info.ModTime(), info.Size())
	}
	fmt.Print("\n")
}

func compareScans(old map[string]os.FileInfo, new map[string]os.FileInfo) (map[string]os.FileInfo, map[string]os.FileInfo) {

	// Union of all filenames. Length is likely incorrect but a decent guess at a beginning size
	allFilenames := make(map[string]bool, len(old))

	for filename := range old {
		allFilenames[filename] = true
	}

	for filename := range new {
		allFilenames[filename] = true
	}

	unchanged := make(map[string]os.FileInfo)
	changed := make(map[string]os.FileInfo)

	for filename := range allFilenames {
		newInfo, newOk := new[filename]
		oldInfo, oldOk := old[filename]
		if newOk && oldOk {
			if fileInfoEqual(oldInfo, newInfo) {
				unchanged[filename] = oldInfo
			} else {
				changed[filename] = oldInfo // TODO should this be newInfo?
			}
		} else if newOk {
			changed[filename] = newInfo
		} else if oldOk {
			changed[filename] = oldInfo
		}
	}

	return unchanged, changed
}

func fileInfoEqual(a os.FileInfo, b os.FileInfo) bool {
	if a.Name() != b.Name() {
		return false
	}
	if a.Size() != b.Size() {
		return false
	}
	if a.Mode() != b.Mode() {
		return false
	}
	if a.ModTime() != b.ModTime() {
		return false
	}
	return true
}

func scanDirectory(path string) (map[string]os.FileInfo, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// Read all files in dir into fileList
	fileList, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}

	fd.Close()

	// Loop over fileList and create map keyed off of full path and value of fileinfo
	fileInfo := make(map[string]os.FileInfo)
	for i := 0; i < len(fileList); i++ {
		d := fileList[i]
		filename := filepath.Join(path, d.Name())
		info, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}
		fileInfo[filename] = info
	}

	return fileInfo, nil
}
