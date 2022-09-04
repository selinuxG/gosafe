package history

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	warcmd = []string{
		"ifconfig",
		"passwd",
		"shadow",
		"su",
		"find",
		"ping",
		"useradd",
		"nohup",
		"netstat",
		"cron",
	}
	cmdpath  = flag.String("cmdfile", "history.log", "此参数是history模式下指定监测命令文件")
	warnpath = flag.String("warnfile", "warnhistory.log", "此参数是history模式下指定异常命令写入到的监测命令文件")
)

func Run() {
	flag.Parse()

	files, err := ioutil.ReadFile(*cmdpath)
	if err != nil {
		log.Println(err)
	}

	warnfile, err := os.OpenFile(*warnpath, os.O_APPEND|os.O_WRONLY, 0766)
	if err != nil {
		log.Println(err)
		return
	}
	defer warnfile.Close()
	write := bufio.NewWriter(warnfile)

	for _, v := range strings.Split(string(files), "\n") {
		if Warncmd(v) {
			//timeStr := time.Now().Format("2006-01-02_15:04:05")
			warcmd := "存在可疑命令:" + v
			write.WriteString(warcmd + "\n")
		}
	}
	write.Flush()
}

func Warncmd(cmd string) bool {
	for _, v := range warcmd {
		if find := strings.Contains(cmd, v); find {
			return true
		}
	}
	return false
}
