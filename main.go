package main

import (
	"fmt"
	"os"
)

// TODO: Use Cobra and Viper to setup cli
// TODO: investigate git-credentials-helper to see how I can mimic that
// TODO: Investigate ctx and cancel funcs to ensure running of background process
// TODO: Use go-daemon to run in background as daemon
// TODO: Should wait until initial auth before daemon and being quiet, so it can error on terminal load like gpg-agent

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
