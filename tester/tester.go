package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/yosssi/gohtml"

	"github.com/bmdelacruz/tplinator"
)

func main() {
	docBytes, err := ioutil.ReadFile("tester.html")
	if err != nil {
		panic(err)
	}

	tplinateStart := time.Now()
	pt, err := tplinator.Tplinate(bytes.NewReader(docBytes))
	if err != nil {
		panic(err)
	}
	tplinateEnd := time.Since(tplinateStart)
	log.Println("tplinate took", tplinateEnd)

	count := 100
	var totalDuration time.Duration
	for i := 0; i < count; i++ {
		_, duration := execute(pt)
		totalDuration += duration

		log.Println("exec", i, "took", duration)
	}
	log.Println("exec took avg", time.Duration(int64(totalDuration)/int64(count)))

	bytes, err := pt.Execute(map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	gohtml.Condense = true
	fmt.Printf("%s\n", gohtml.FormatBytes(bytes))
}

func execute(pt *tplinator.PrecompiledTemplate) ([]byte, time.Duration) {
	execStart := time.Now()
	bytes, err := pt.Execute(map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	return bytes, time.Since(execStart)
}
