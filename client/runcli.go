package client

import (
	_ "embed"
	"io/ioutil"
	"log"
	"os/exec"
)

//go:embed goclient.sh
var goclish string

func Runcli() {
	ioutil.WriteFile("clientgo.sh", []byte(goclish), 744) //写入执行脚本
	cmd := exec.Command("nohup", "sh", "/dhcp/clientgo.sh", "10", ">", "goclient.log", "2>&1", "&")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	log.Println("seccess")
}
