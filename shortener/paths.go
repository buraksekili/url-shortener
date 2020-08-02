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

func InitPathStruct() (map[string]string, error) {
	flag.Parse()

	if !(*redirect) {
		fmt.Println("wow")
		if len(*from) == 0 {
			return nil, fmt.Errorf("-f parameter is required!\nType -h or --help for further instructions.\n")
		}
		if len(*to) == 0 {
			return nil, fmt.Errorf("-t parameter is required!\nType -h or --help for further instructions.\n")
		}

		if !strings.HasPrefix(*to, "http") {
			return nil, fmt.Errorf("-t must starts with http(s)\n")
		}
		if !strings.HasPrefix(*from, "/") {
			return nil, fmt.Errorf("-f must starts with /\n")
		}
	}
	err := generateJSON(*filename)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	file, err := ioutil.ReadFile(*filename)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	var paths []Path
	err = json.Unmarshal(file, &paths)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	inpPath := Path{From: *from, To: *to}
	if !pathExists(&inpPath, &paths) {
		newPaths := append(paths, inpPath)
		byteArr, err := json.Marshal(newPaths)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("%s\n", err)
		}

		err = ioutil.WriteFile(*filename, byteArr, 0644)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("%s\n", err)
		}
		fmt.Printf("==== localhost:3000%s redirects to %s ====\n", inpPath.From, inpPath.To)
		return getURLS(newPaths), nil
	}
	fmt.Printf("==== localhost:3000%s redirects to %s ====\n", inpPath.From, inpPath.To)
	return getURLS(paths), nil
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
		err = InitJSON(fname)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitJSON(fname string) error {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	_, err = file.WriteString("[]")
	if err != nil {
		err = file.Close()
		if err != nil {
			fmt.Println("Json file couldn't closed!", err)
		}
		return fmt.Errorf("%s\n", err)
	}
	return nil
}
