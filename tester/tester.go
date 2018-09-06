package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bmdelacruz/tplinator"
	"github.com/yosssi/gohtml"
)

func main() {
	startTime := time.Now()
	templateFile, err := os.Open("tester.html")
	if err != nil {
		log.Fatalln(err)
	}
	duration := time.Since(startTime)
	fmt.Println("read file:", duration)

	startTime = time.Now()
	tpl, err := tplinator.Tplinate(templateFile)
	if err != nil {
		log.Fatalln(err)
	}
	duration = time.Since(startTime)
	fmt.Println("tplinated:", duration)

	pets := tplinator.RangeParams(
		tplinator.EvaluatorParams{
			"name":        "Larry",
			"description": "A wonderful dog",
			"imgurl":      "/images/dog.png",
		},
		tplinator.EvaluatorParams{
			"name":        "Perry",
			"description": "My best friend",
			"imgurl":      "/images/platypus.png",
		},
	)

	startTime = time.Now()
	str, err := tpl.RenderString(map[string]interface{}{
		"hasPets": len(pets) > 0,
		"pets":    pets,
	})
	if err != nil {
		log.Fatalln(err)
	}
	duration = time.Since(startTime)
	fmt.Println("rendered:", duration)

	gohtml.Condense = true
	fmt.Println(gohtml.Format(str))
}
