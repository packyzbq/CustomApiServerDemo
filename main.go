package main

import (
	"CustomApiServerDemo/pkg/cmd/server"
	"flag"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"os"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewCustomServerOptions(os.Stdout, os.Stderr)
	cmd := server.NewCommandStartCustomServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
