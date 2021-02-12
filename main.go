package main

import (
	"fmt"
	"github.com/buraksekili/url-shortener/shortener"
	"github.com/fatih/color"
	"net/http"
)

func main() {

	url, err := shortener.CheckOp()
	if err != nil {
		fmt.Println(err)
		return
	}
	urlHandler := shortener.URLHandler(url, http.NewServeMux())
	fmt.Println("Live on port: ", 3000)
	if err := http.ListenAndServe(":3000", urlHandler); err != nil {
		fmt.Printf("%s %s\n", color.RedString("Error:"), err)
	}
}
