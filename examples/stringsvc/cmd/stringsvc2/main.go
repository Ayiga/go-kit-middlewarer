package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	trans "github.com/ayiga/go-kit-middlewarer/examples/stringsvc/transport/http"
)

type arguments struct {
	httpPort string
}

var args arguments

func init() {
	flag.StringVar(&args.httpPort, "httpPort", ":9000", "Specifies which port to listen for requests on")
}

func usage() {
	fmt.Printf("%s: server [-httpPort :port|-httpPort=:port]\n%s: client cmd arg [-httpPort :port|-httpPort=:port]\n", os.Args[0], os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	a := flag.Args()

	var mode = ""
	if len(a) > 0 {
		mode = a[0]
	}

	switch mode {
	case "server":
		var svc StringService
		trans.HTTPServersForEndpoints(svc)
		http.ListenAndServe(args.httpPort, nil)
	case "client":
		if len(a) < 3 {
			usage()
			return
		}

		client := trans.NewHTTPClient("127.0.0.1" + args.httpPort)

		cmd := a[1]
		arg := a[2]
		switch cmd {
		case "uppercase":
			str, err := client.Uppercase(arg)
			fmt.Printf("\t\"%s\".Uppercase: \"%s\", err: %s\n", arg, str, err)
		case "count":
			count := client.Count(arg)
			fmt.Printf("\t\"%s\".Count: %d\n", arg, count)
		}

	default:
		usage()
	}

}
