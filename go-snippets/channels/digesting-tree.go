// +build OMIT

package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error { // HL
		if err != nil {
			// fmt.Print("1::::::::::::::")
			return err
		}
		if !info.Mode().IsRegular() {
			// fmt.Print("2::::::::::::::")
			return nil
		}
		data, err := ioutil.ReadFile(path) // HL
		if err != nil {
			//fmt.Print("3::::::::::::::")
			return nil
		}
		m[path] = md5.Sum(data) // HL
		return nil
	})
	if err != nil {
		return nil, err
		fmt.Print("4::::::::::::::")
	}
	return m, nil
}

func main() {
	fmt.Println("Starting MD5 digest...")
	timeBegin := time.Now()
	// Calculate the MD5 sum of all files under the specified directory,
	// then print the results sorted by path name.
	m, err := MD5All(os.Args[1]) // HL
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}

	durationDigest := time.Since(timeBegin)

	timeSortAndPrint := time.Now()

	sort.Strings(paths) // HL
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}

	durationSort := time.Since(timeSortAndPrint)

	fmt.Printf("Digested %v files in %v ns and printed in %v", len(paths), durationDigest, durationSort)
}
