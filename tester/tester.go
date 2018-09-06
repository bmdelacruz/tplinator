package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yosssi/gohtml"

	"github.com/bmdelacruz/tplinator"
)

func main() {
	templateFile, err := os.Open("tester.html")
	if err != nil {
		log.Fatalln(err)
	}
	tpl, err := tplinator.Tplinate(templateFile)
	if err != nil {
		log.Fatalln(err)
	}
	str, err := tpl.RenderString(map[string]interface{}{
		"hasPets": true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	gohtml.Condense = true
	fmt.Println(gohtml.Format(str))
}
