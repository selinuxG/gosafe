package info

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	serverpath = flag.String("infofile", "server.log", "此参数是info模式下分析的日志文件路径")
	infolog    = flag.String("infolog", "false", "此参数是info模式下是否开启详细日志输出信息")
	//ssh成功登陆
	sshaclist = []string{}
	sshacip   = []string{}
	//ssh失败登陆
	sshfalist = []string{}
	sshfaip   = []string{}
	//nginx404记录
	nginx404     int
	nginx404list = []string{}
	//nginx302记录
	nginx302     int
	nginx302list = []string{}
	//可以命令记录
	warnhistory     int
	warnhistorylist = []string{}
)

func Run() {
	flag.Parse()
	f, err := ioutil.ReadFile(*serverpath)
	if err != nil {
		if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
			log.Println("指定监控文件不存在,程序退出！")
		} else {
			log.Println(err)
		}
		os.Exit(3)
	}
	//SSH登陆成功的记录，写入到切片
	for _, data := range strings.Split(string(f), "\n") {
		if sshac := strings.Contains(data, "Ac"); sshac {
			if GIN := strings.Contains(data, "GIN"); GIN {
			} else {
				sshaclist = append(sshaclist, data)
				re := regexp.MustCompile(`[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}`)
				sshac := re.FindAllString(data, -1)
				for _, v := range sshac {
					sshacip = append(sshacip, v)
				}
			}
		}
		//SSH登陆失败的记录，写入到切片
		if sshac := strings.Contains(data, "Fa"); sshac {
			if GIN := strings.Contains(data, "GIN"); GIN {
			} else {
				sshfalist = append(sshfalist, data)
				re := regexp.MustCompile(`[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}`)
				sshac := re.FindAllString(data, -1)
				for _, v := range sshac {
					sshfaip = append(sshfaip, v)
				}
			}
		}
		//NGINX 404记录
		if four0four := strings.Contains(data, "404"); four0four {
			if GIN := strings.Contains(data, "GIN"); GIN {
			} else {
				if four0four := strings.Contains(data, "状态码为:404"); four0four {
					nginx404 += 1
					nginx404list = append(nginx404list, data)
				}
			}
		}
		//NGINX 302记录
		if three0three := strings.Contains(data, "302"); three0three {
			if GIN := strings.Contains(data, "GIN"); GIN {
			} else {
				if four0four := strings.Contains(data, "状态码为:302"); four0four {
					nginx302 += 1
					nginx302list = append(nginx302list, data)
				}
			}
		}
		//可疑历史命令记录
		if warhi := strings.Contains(data, "存在可疑命令"); warhi {
			if GIN := strings.Contains(data, "GIN"); GIN {
			} else {
				warnhistory += 1
				warnhistorylist = append(warnhistorylist, data)
			}
		}

	}

	log.Println("--------登陆成功IP统计")
	for s, v := range calSliceCount(sshacip) {
		log.Printf("IP：%s 登陆次数总计:%d", s, v)
	}

	log.Println("--------登陆失败IP统计")
	for s, v := range calSliceCount(sshfaip) {
		log.Printf("IP：%s 失败登陆次数总计:%d", s, v)
	}
	log.Println("--------可疑历史命令统计")
	log.Printf("共检测到可疑命令:%d条数据", warnhistory)
	log.Println("--------nginx请求统计")
	log.Printf("共检测到404状态:%d条数据", nginx404)
	log.Printf("共检测到302状态:%d条数据", nginx302)

	sshacc := make(map[string]interface{})
	sshfac := make(map[string]interface{})
	nginx := make(map[string]interface{})
	nginx["404"] = nginx404
	nginx["302"] = nginx302
	history := make(map[string]interface{})
	//history["可疑执行命令次数"] = warnhistory
	for s, v := range calSliceCount(sshacip) {
		sshacc[s] = v
	}
	for s, v := range calSliceCount(sshfaip) {
		sshfac[s] = v
	}
	for s, v := range calSliceCount(warnhistorylist) {
		history[s] = v
	}
	//log.Println(history)
	//生成html图表
	Runhtml(sshfac, sshacc, nginx, history)
	//html.Run(sshfac, sshacc, nginx, history)

	if *infolog == "true" {
		log.Println("--------系统登陆成功日志信息")
		for _, v := range sshaclist {
			log.Printf(v)
		}
		log.Println("--------登陆登陆失败日志信息")
		for _, v := range sshfalist {
			log.Printf(v)
		}
		log.Println("--------登陆可疑命令日志信息")
		for _, v := range warnhistorylist {
			log.Printf(v)
		}
		log.Println("--------Nginx404日志信息")
		for _, v := range nginx404list {
			log.Printf(v)
		}
		log.Println("--------Nginx302日志信息")
		for _, v := range nginx302list {
			log.Printf(v)
		}

	}
}

//string切片去重计数
func calSliceCount(slices []string) map[string]int {
	m := map[string]int{}
	for _, ele := range slices {
		m[ele]++
	}
	return m
}
