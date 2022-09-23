package main

import (
	"flag"
	"fmt"
	"os"

	klog "k8s.io/klog/v2"
)

// TODO: Use Cobra and Viper to setup cli
// TODO: investigate git-credentials-helper to see how I can mimic that
// TODO: Investigate ctx and cancel funcs to ensure running of background process
// TODO: Investigate ssh-agent to see how to handle stdout and stderr
// TODO: Investigate ssh-agent to see how the process runs in the background
// TODO: Enssure only ever running once, use a pid file?

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cli, err := ProvideCLI()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	cli.Execute()
}
