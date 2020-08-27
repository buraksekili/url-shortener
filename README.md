#url-shortener

## Usage
```
    go build main.go
    ./main.go -e=gh-http://www.github.com
```
```
    go run main.go -e=gh-http://www.github.com -f="custom.json"
```
## Available flags
```
USAGE:
-e=<FROM>-<TO> string(new entry)
        Enter new entry as <FROM>-<TO> (i.e, burak-http://www.google.com).
        Now, http://127.0.0.1/burak redirects to http://www.google.com
        NOTE: Do not use '/' as a prefix in <FROM> due to prevent compatibility issues with Windows OS.

-f=<FILENAME>, --fname=<FILENAME> string(filename operator)
        OPTIONAL. The name of the json file that stores paths. (default 'paths.json')
        If <FILENAME> is not specified, it gives an error as 'unsupported operation'.

<ADDRESS> string
        If .json file includes <ADDRESS>, <ADDRESS> can be used in redirection with this option.

-h, --help boolean
        Prints this message
```