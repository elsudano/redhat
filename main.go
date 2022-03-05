package main

import (
	"flag"
	"fmt"

	"github.com/elsudano/redhat/redhat"
)

func main() {
	url := flag.String("url", "", "You need put the URL from download the file")
	fix := flag.Bool("fix", false, "By default it's false, but if you want to see the implementation with teh correct JSON change at true")
	flag.Parse()

	if *url != "" && *fix {
		fmt.Printf("%+v\n", redhat.JsonImplementation(url))
	} else if *url != "" {
		fmt.Printf(redhat.DefaultImplementation(url))
	} else {
		flag.Usage()
	}
}
