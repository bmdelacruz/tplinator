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

	params := map[string]interface{}{
		"hastwo":   false,
		"hasthree": false,
		"hasfour":  false,

		"uid":      "22233334234123",
		"secret":   "asodijfefeaofiasjdosads",
		"confcode": "sdasdetwet34twfesfsefsee",
	}

	count := 10000
	var totalDuration time.Duration
	for i := 0; i < count; i++ {
		_, duration := execute(pt, params)
		totalDuration += duration

		log.Println("exec", i, "took", duration)
	}
	if count > 0 {
		log.Println("exec took avg", time.Duration(int64(totalDuration)/int64(count)))
	}

	bytes, err := pt.Execute(params)
	if err != nil {
		panic(err)
	}
	gohtml.Condense = true
	fmt.Printf("%s\n", gohtml.FormatBytes(bytes))
}

func execute(pt *tplinator.PrecompiledTemplate, params map[string]interface{}) ([]byte, time.Duration) {
	execStart := time.Now()
	bytes, err := pt.Execute(params)
	if err != nil {
		panic(err)
	}
	return bytes, time.Since(execStart)
}
