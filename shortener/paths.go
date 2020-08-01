package shortener

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var filename = flag.String("n", "paths.json", "OPTIONAL. The name of the json file that stores paths.")
var redirect = flag.Bool("r", false, "OPTIONAL. If -r is true, you can use redirection without specifying new url.")
var from = flag.String("f", "", "REQUIRED in edit mode. The name of the site that will be shortened. It needs to include http or https part of the url. (From site -f to site -t)")
var to = flag.String("t", "", "REQUIRED in edit mode. The name of the target website. It needs to include http or https part of the url. (From site -f to site -t)")

type Path struct {
	From string `json:"from, omitempty"`
	To   string `json:"to, omitempty"`
}

func InitPathStruct() map[string]string {
	flag.Parse()
	if *redirect {
		fmt.Println("REDIRECTION MODE!")
		return nil
	}
	if len(*from) == 0 {
		fmt.Println("-f parameter is required!\nType -h or --help for further instructions.")
		return nil
	}
	if len(*to) == 0 {
		fmt.Println("-t parameter is required!\nType -h or --help for further instructions.")
		return nil
	}

	if !strings.HasPrefix(*to, "http") {
		fmt.Println("-t must starts with http(s)")
		return nil
	}
	if !strings.HasPrefix(*from, "/") {
		fmt.Println("-t must starts with /")
		return nil
	}

	err := generateJSON(*filename)
	if err != nil {
		fmt.Println("generateJSON |", err)
		return nil
	}

	file, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println("ReadFile |", err)
		return nil
	}

	var paths []Path
	err = json.Unmarshal(file, &paths)
	if err != nil {
		fmt.Println("Unmarshal |", err)
		return nil
	}

	inpPath := Path{From: *from, To: *to}
	if !pathExists(&inpPath, &paths) {
		newPaths := append(paths, inpPath)
		byteArr, err := json.Marshal(newPaths)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		err = ioutil.WriteFile(*filename, byteArr, 0644)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Printf("==== localhost:3000%s redirects to %s ====\n", inpPath.From, inpPath.To)
		return getURLS(newPaths)
	}
	fmt.Printf("==== localhost:3000%s redirects to %s ====\n", inpPath.From, inpPath.To)
	return getURLS(paths)
}

func pathExists(path *Path, paths *[]Path) bool {
	for _, v := range *paths {
		if (*path).From == v.From {
			return true
		}
	}
	return false
}

func getURLS(paths []Path) map[string]string {
	urls := make(map[string]string)
	for _, s := range paths {
		if _, exists := urls[s.From]; !exists {
			urls[s.From] = s.To
		}
	}
	return urls
}

func generateJSON(fname string) error {
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		_, err := os.Create(fname)
		if err != nil {
			fmt.Println(err)
			return err
		}
		initJSON(fname)
	}
	return nil
}

func initJSON(fname string) {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("Json file couldn't open.", err)
		return
	}

	_, err = file.WriteString("[]")
	if err != nil {
		fmt.Println("Json file couldn't initialized: ", err)
		err = file.Close()
		if err != nil {
			fmt.Println("Json file couldn't closed!", err)
		}
		return
	}
}
