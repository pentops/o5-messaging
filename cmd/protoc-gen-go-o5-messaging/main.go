package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/pentops/o5-messaging/internal/protocgen"
	"google.golang.org/protobuf/compiler/protogen"
)

const Version = "1.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-o5-messaging %v\n", Version)
		return
	}

	var flags flag.FlagSet
	extraHeadersString := flags.String("extra_headers", "", ": separated key:val:key:val pairs to add to gRPC headers")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(g *protogen.Plugin) error {

		extraHeaders := make([]protocgen.KeyValue, 0)

		if extraHeadersString != nil && *extraHeadersString != "" {
			parts := strings.Split(*extraHeadersString, ":")
			if len(parts)%2 != 0 {
				return fmt.Errorf("uneven number of extra_headers")
			}
			for x := 0; x < len(parts); x += 2 {
				extraHeaders = append(extraHeaders, protocgen.KeyValue{
					Key:  strings.TrimSpace(parts[x]),
					Code: []any{`"`, strings.TrimSpace(parts[x+1]), `"`},
				})
			}
		}

		return protocgen.Config{
			ExtraHeaders: extraHeaders,
		}.Generate(g)
	})
}
