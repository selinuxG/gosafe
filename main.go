package main

import (
	"flag"
	"gosafe/client"
	"gosafe/history"
	"gosafe/info"
	"gosafe/server"
	"gosafe/tail"
)

var (
	run = flag.String("run", "server", "此参数启动模式:history、server、tail、info")
)

func main() {
	flag.Parse()
	switch *run {
	case "server":
		server.Run()
	case "tail":
		tail.Run()
	case "history":
		history.Run()
	case "info":
		info.Runxlsx()
	case "client":
		client.Runcli()
	}
}
