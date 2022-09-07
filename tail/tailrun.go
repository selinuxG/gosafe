package tail

import (
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	path     = flag.String("path", "/var/log/secure", "此参数是tail模式下监测目录或文件")
	serverip = flag.String("serverip", "http://103.44.250.69:17171", "此参数是tail模式下指定server服务端的IP")
	dirfiles []string //递归目录的切片
	wg       sync.WaitGroup
)

func Run() {
	flag.Parse()
	paths := strings.Split(*path, ",") //基于“,“分隔,长度为1=监测单个
	if len(paths) == 1 {
		switch {
		case Dirs(*path) == true: //如果返回true说明是目录，则递归生成的目录的切片并发交给Tail模式处理是
			wg.Add(2)
			for _, v := range dirfiles {
				go Tailfile(v)
			}
		case Dirs(*path) == false:
			Tailfile(*path)
		}
	} else { //执行到这就说明不是一个文件，则循环切片验证是否目录交给Tail模式处理
		wg.Add(len(paths)) //开启线程
		for i := 0; i < len(paths); i++ {
			path := paths[i]
			switch {
			case Dirs(path) == true:
				for _, v := range dirfiles {
					go Tailfile(v)
				}
			case Dirs(path) == false:
				go Tailfile(path)
			}
		}
	}
	wg.Wait() //等待线程结束，否则主线程就直接结束程序退出。
}

func Tailfile(file string) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读,0为开头，2为末行
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,                                 // 监听新行，使用tail -f，这个参数非常重要
	}
	t, _ := tail.TailFile(file, config)
	for line := range t.Lines {
		if len(line.Text) == 0 {
			continue
		}
		//line.Text = strings.Split(line.Text, "\r")[0] //widnos下换行符包含/r
		line.Text = strings.Replace(line.Text, "\r", "", -1)
		line.Text = strings.Split(line.Text, "\n")[0]
		line.Text = strings.Replace(line.Text, ";", "", -1)    //有空格的数据传参测试有问题需转换
		line.Text = strings.Replace(line.Text, " ", `%20`, -1) //有空格的数据传参测试有问题需转换
		url := *serverip + "/warn?"
		url = fmt.Sprintf("%s&path=%s&text=%s", url, file, line.Text) //拼接url
		//log.Println(line.Text)
		//log.Println(url)
		Get(url)
	}
}

func Dirs(file string) bool {
	fileInfo, err := os.Stat(file)
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		log.Println("指定监控文件不存在,程序退出！")
		os.Exit(3)
	}
	if err != nil {
		fmt.Println("err :", err)
		return false
	}
	if fileInfo.IsDir() {
		WalkDir(file)
		return true
	}
	return false
}

//请求接口传参
func Get(url string) {
	// 超时时间：5秒
	client := &http.Client{Timeout: 3 * time.Second}
	_, err := client.Get(url)
	if err != nil {
		log.Println("err:", "发送数据失败")
	} else {
		log.Println("成功将监控日志发送到server...")
	}
}

//递归目录,开启监测
func WalkDir(filePath string) {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, v := range files {
			if v.IsDir() {
				newfile := filePath + "/" + v.Name()
				WalkDir(newfile)
				return
			} else {
				newfile := filePath + "/" + v.Name()
				dirfiles = append(dirfiles, newfile) //添加到全局切片中
			}
		}
	}
}
