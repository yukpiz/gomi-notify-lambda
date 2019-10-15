package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	debug = flag.Bool("d", false, "Console debug mode")
)

func Handler(ctx context.Context) error {
	fmt.Fprintf(os.Stdout, "START ===>")
	return nil
}

func main() {
	flag.Parse()
	if !*debug {
		lambda.Start(Handler)
	}
	if err := Handler(context.Background()); err != nil {
		panic(err)
	}
}
