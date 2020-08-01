package main

import (
	"./shortener"
	"fmt"
	"net/http"
)

// TODO already existence links doesn't redirect
func main() {
	urls := shortener.InitPathStruct()
	urlHandler := shortener.URLHandler(urls, http.NewServeMux())
	fmt.Println("Live on port: ", 3000)
	if err := http.ListenAndServe(":3000", urlHandler); err != nil {
		fmt.Println("error:", err)
	}

}
