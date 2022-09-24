package main

import (
	"fmt"
	"os"
)

// TODO: Use Cobra and Viper to setup cli
// TODO: investigate git-credentials-helper to see how I can mimic that
// TODO: Investigate ctx and cancel funcs to ensure running of background process
// TODO: Investigate ssh-agent to see how to handle stdout and stderr
// TODO: Investigate ssh-agent to see how the process runs in the background
// TODO: Enssure only ever running once, use a pid file?

func main() {
	cli, err := ProvideCLI()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
