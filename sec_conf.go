package main

import (
	"fmt"
	"os"
	"secpaver/secconf"
)

func main() {
	if err := secconf.NewSecConfApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}
