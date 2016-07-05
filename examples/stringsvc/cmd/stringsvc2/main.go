package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"

	httptrans "github.com/go-kit/kit/transport/http"

	"github.com/ayiga/go-kit-middlewarer/encoding"
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
	flag.Usage = usage
	flag.Parse()
	a := flag.Args()

	var mode = ""
	if len(a) > 0 {
		mode = a[0]
	}

	switch mode {
	case "server":
		var svc StringService
		options := []httptrans.ServerOption{httptrans.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
			w.WriteHeader(500)
			encoding.JSON(1).EncodeResponse()(ctx, w, err)
		}),
		}
		trans.HTTPServersForEndpointsWithOptions(svc, []trans.ServerLayer{}, options)
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
