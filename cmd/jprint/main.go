package main

import (
	"fmt"
	"github.com/wdsgyj/jclass"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	err := filepath.Walk(os.Args[1], walk)
	if err != nil {
		log.Fatalln(err)
	}
}

func walk(path string, info os.FileInfo, err error) error {
	if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".class") {
		classFile, err := jclass.NewClassFileFromPath(path)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("//", path)
		fmt.Println(classFile)
		fmt.Println()
	}
	return nil
}
