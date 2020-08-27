package shortener

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
)

func readJSON(fname string) ([]Path, error) {
	err := generateJSON(fname)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	var paths []Path
	err = json.Unmarshal(file, &paths)
	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}
	return paths, nil
}

func CheckOp() (map[string]string, error) {
	var op Op
	op = ParseFlags(os.Args[1:])

	switch v := op.(type) {
	case HelpOp:
		printHelp(os.Stdout)
	case FilenameOp:
		if !strings.HasSuffix(v.Name, "json") {
			fmt.Printf("%s Expected <FILENAME> type is *.json, but got '%s'\n", color.RedString("Error:"), v.Name)
			printHelp(os.Stdout)
			os.Exit(1)
		}
		var paths []Path
		paths, err := readJSON(v.Name)
		if err != nil {
			fmt.Printf("Error in readJSON %s\n", err)
			os.Exit(1)
		}
		byteArr, err := json.Marshal(paths)
		err = ioutil.WriteFile(v.Name, byteArr, 0644)
		if err != nil {
			fmt.Println("error in ioutil.Writefile:", err)
			os.Exit(1)
		}
	case FindOp:
		target := v.target
		var fname string = "paths.json"

		// If filename is specified via flags, update fname
		if len(v.filename) > 0 {
			fname = v.filename
		}
		// Usage of '/' as a prefix disabled because it causes errors in Windows.
		if strings.HasPrefix(target, "/") || strings.Contains(target, ":/") {
			fmt.Printf(color.RedString("Do not use '/' as a prefix.\n"))
			printHelp(os.Stdout)
			return nil, fmt.Errorf(color.RedString("Do not use '/' as a prefix."))
		}

		target = "/" + target
		var paths []Path
		paths, err := readJSON(fname)
		if err != nil {
			fmt.Printf("%s %s\n", color.RedString("Error: "), err)
		}
		if exists, val := pathExists(target, &paths); exists {
			fmt.Printf("%s\t%s exists in %s as:%s\n", color.GreenString("DONE!"), target, fname, val)
			urls := make(map[string]string)
			urls[target] = val
			return urls, nil
		} else {
			return nil, fmt.Errorf("%s\n", color.MagentaString(target+" cannot found in "+fname))
		}
	case EntryOp:
		fmt.Printf("from: %s\ttarget: %s\tfile: %s\n", v.From, v.Target, v.FileName)
		filename := v.FileName
		err := generateJSON(filename)
		if err != nil {
			return nil, fmt.Errorf("%s\n", err)
		}

		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("%s\n", err)
		}

		var paths []Path
		err = json.Unmarshal(file, &paths)
		if err != nil {
			return nil, fmt.Errorf("%s\n", err)
		}

		newFrom := "/" + v.From
		if exists, _ := pathExists(newFrom, &paths); !exists {
			newPaths := append(paths, Path{From: newFrom, To: v.Target})
			byteArr, err := json.Marshal(newPaths)
			if err != nil {
				fmt.Println(err)
				return nil, fmt.Errorf("%s\n", err)
			}

			err = ioutil.WriteFile(filename, byteArr, 0644)
			if err != nil {
				fmt.Println(err)
				return nil, fmt.Errorf("%s\n", err)
			}
			return getURLS(newPaths), nil
		}
		return nil, fmt.Errorf("lol")

	case UnknownOp:
		fmt.Printf("%s unsupported operation: %s\n", color.RedString("Error:"), v.args)
		if v.args[0] == "-f" || v.args[0] == "--fname" {
			fmt.Printf("%s %s requires <FILENAME>.\n", color.RedString("Note!"), v.args[0])
		}
		printHelp(os.Stdout)
		os.Exit(1)
	default:
		fmt.Printf("Operation %T is not handled\n", op)
	}
	return nil, nil
}

func pathExists(path string, paths *[]Path) (bool, string) {
	for _, v := range *paths {
		fmt.Printf("From: %s\tTo: %s\n", v.From, v.To)
		if path == v.From {
			return true, v.To
		}
	}
	return false, ""
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
