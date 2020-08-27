package shortener

import (
	"strings"
)

type Path struct {
	From string `json:"from, omitempty"`
	To   string `json:"to, omitempty"`
}

type Op interface{}

// HelpOp describes printing helper commands.
type HelpOp struct{}

// FilenameOp describes name of the JSON file
type FilenameOp struct {
	Name string
}

// FindOp finds input site name in the JSON.
type FindOp struct {
	target   string
	filename string
}

// EntryOp describes the site which will be redirected.
type EntryOp struct {
	From     string
	Target   string
	FileName string
}

// UnknownOp describes the unsupported flags.
type UnknownOp struct {
	args []string
}

func ParseFlags(argv []string) Op {
	// If no args indicated, show help operator
	if len(argv) == 0 {
		return HelpOp{}
	}
	if len(argv) == 1 {
		f := strings.TrimSpace(argv[0])
		// Handle Help operator
		if f == "--help" || f == "-h" {
			return HelpOp{}
		}

		if strings.HasPrefix(f, "-e") {
			// if -e flags' length is not 3 at least, input can be assumed as an invalid.
			if len(f) < 3 {
				return UnknownOp{args: argv}
			}

			if currf := f[3:]; len(currf) > 1 {
				if strings.Contains(currf, "-") {
					currfArr := strings.Split(currf, "-")
					// if len(argv) == 1, we cannot have additional filename. So, paths.json
					// is used as FileName field of EntryOp.
					return EntryOp{Target: currfArr[1], From: currfArr[0], FileName: "paths.json"}
				}
			}
			return EntryOp{FileName: "paths.json"}
		}

		// Rest of the flags starting with '-' should be unrecognized flag.
		if strings.HasPrefix(f, "-") {
			return UnknownOp{args: argv}
		}

		// If there is no flag left, i need to find variable f in .json file.
		return FindOp{target: f, filename: "paths.json"}
	}

	if len(argv) == 2 {

		f1 := argv[0]
		f2 := argv[1]

		if !strings.HasPrefix(f1, "-") || !strings.HasPrefix(f2, "-") {
			// we are in find mode
			fname, searching := parseFindFlags(f1, f2)
			if len(fname) == 0 || len(searching) == 0 {
				return UnknownOp{args: argv}
			}
			return FindOp{target: searching, filename: fname}
		}

		var fname string
		var newPath Path
		if strings.HasPrefix(f1, "-e") && strings.HasPrefix(f2, "-f") || strings.HasPrefix(f2, "--fname") {
			newPath, fname = parseEntryArgs(f1, f2)
			return EntryOp{From: newPath.From, Target: newPath.To, FileName: fname}
		} else if strings.HasPrefix(f2, "-e") && strings.HasPrefix(f1, "-f") || strings.HasPrefix(f1, "--fname") {
			newPath, fname = parseEntryArgs(f2, f1)
			return EntryOp{From: newPath.From, Target: newPath.To, FileName: fname}
		}

		return UnknownOp{args: argv}
	}
	// TODO handle too many arguments
	return UnknownOp{}
}

// parseFindFlags parses flags for FindOp flags.
// One of the arguments must belong to FindOp.
func parseFindFlags(f1 string, f2 string) (string, string) {
	var fnameFlag string
	var findFlag string

	// If one of the arguments(f1, f2) starts with -f or --fname, it must be FileNameOp flag
	if strings.HasPrefix(f1, "-f") || strings.HasPrefix(f1, "--fname") {
		fnameFlag = f1
		findFlag = f2
	} else if strings.HasPrefix(f2, "-") || strings.HasPrefix(f2, "--fname") {
		fnameFlag = f2
		findFlag = f1
	} else {
		return "", ""
	}
	if strings.Contains(fnameFlag, "=") {
		fnameArr := strings.Split(fnameFlag, "=")
		if len(fnameArr) != 2 {
			return "", ""
		}
		return fnameArr[1], findFlag
	}
	return "", ""
}

// parseArgs parses flags to obtain their values.
// First argument must be entry flag string and the second one must be filename flag.
func parseEntryArgs(ef string, ff string) (Path, string) {
	if strings.Contains(ef, "=") {
		entries := strings.Split(ef, "=")
		if len(entries) != 2 {
			return Path{}, ""
		}
		entriesInfo := strings.Split(entries[1], "-")
		from := entriesInfo[0]
		target := entriesInfo[1]

		fname := getFname(ff)
		if len(fname) == 0 {
			return Path{}, ""
		}

		newPath := Path{From: from, To: target}
		return newPath, fname
	}
	return Path{}, ""
}

func getFname(ff string) string {
	var fArr []string
	if !strings.Contains(ff, "=") {
		return ""
	}

	fArr = strings.Split(ff, "=")
	if len(fArr) != 2 {
		return ""
	}
	return fArr[1]
}
