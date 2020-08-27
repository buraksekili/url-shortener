package shortener

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer) {
	help := "USAGE:\n " +
		"-e=<FROM>-<TO> string(new entry)\n        Enter new entry as <FROM>-<TO> (i.e, burak-http://www.google.com).\n        " +
		"Now, http://127.0.0.1/burak redirects to http://www.google.com\n        " +
		"NOTE: Do not use '/' as a prefix in <FROM> due to prevent compatibility issues with Windows OS.\n" +

		"\n-f=<FILENAME>, --fname=<FILENAME> string(filename operator)\n        OPTIONAL. The name of the json file that stores paths. (default 'paths.json')\n        " +
		"If <FILENAME> is not specified, it gives an error as 'unsupported operation'.\n" +

		"\n<ADDRESS> string\n        If .json file includes <ADDRESS>, <ADDRESS> can be used in redirection with this option.\n" +
		"\n-h, --help boolean\n        Prints this message\n"

	fmt.Fprintf(w, "%s\n", help)
}
