package info

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	apiserver = flag.String("apiserver", "http://103.44.250.69:17171/api", "此参数是info模式下调用的接口数据")
	pageSize  = flag.String("pageSize", "200", "此参数是info模式下每页返回的数量")
	apichan   = flag.Int("apichan", 20, "此参数是info模式下控制的线程数量")
	infoname  = flag.String("xmlname", "serverwarn.xlsx", "此参数是info模式生成的xlsx文件名称")
	ii        = 1
	emo       = []string{"😀", "😁", "🤣", "😃", "😎", "😛", "🥱", "🤬", "😡", "❤"}
	rwLock    sync.RWMutex
)

// Warn json解析到结构体
type Warn struct {
	IP    string
	COUNT string
	PATH  string
	INFO  string
	DATE  string
	TYPE  string
}

func Runxlsx() {
	flag.Parse()
	createxml() //创建severwarn.xlsx
	log.Printf("启动info模式%s\n调用接口地址：%v 线程数为:%v API每页返回数为:%v 输出文件为:%s", emo[rand.Intn(len(emo))], *apiserver, *apichan, *pageSize, *infoname)
	start := time.Now() // 获取程序启动时间
	f, err := excelize.OpenFile(*infoname)
	if err != nil {
		log.Println("打开xlsx文件的时候时候报错了", err)
		return
	}
	pageCount := apicount(*apiserver) //总页数
	ch := make(chan struct{}, *apichan)
	wg := sync.WaitGroup{}

	for i := 1; i < pageCount+1; i++ {
		ch <- struct{}{}
		wg.Add(1)
		apiurl := fmt.Sprintf("%s?page=%d&pageSize=%s", *apiserver, i, *pageSize) //拼接APIurl
		go func(apiurl string) { //开启线程
			defer wg.Done() //函数退出时线程减1
			resp, err := http.Get(apiurl)
			if err != nil {
				log.Println("应该是并发太高请求频率太快链接失败了哦~", apiurl)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			var info []Warn                   //报警数据切片
			err = json.Unmarshal(body, &info) //解析json数据到结构体
			if err != nil {
				log.Println("解析json失败:", err)
				return
			}
			for _, v := range info {
				rwLock.Lock()
				ii += 1
				//log.Printf("正在向severwarn.xlsx追加第%d条告警数据...%s", ii-1, emo[rand.Intn(len(emo))])
				//写入接口数据
				intnumber, _ := strconv.Atoi(strconv.Itoa(ii))
				f.SetCellValue("Sheet1", "A"+strconv.Itoa(ii), intnumber-1)
				f.SetCellValue("Sheet1", "B"+strconv.Itoa(ii), v.IP)
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(ii), v.TYPE)
				f.SetCellValue("Sheet1", "D"+strconv.Itoa(ii), v.DATE)
				f.SetCellValue("Sheet1", "E"+strconv.Itoa(ii), v.COUNT)
				f.SetCellValue("Sheet1", "F"+strconv.Itoa(ii), v.PATH)
				f.SetCellValue("Sheet1", "G"+strconv.Itoa(ii), v.INFO)
				rwLock.Unlock()
			}
			<-ch
		}(apiurl)
	}
	wg.Wait()
	//保存数据
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}

	elapsed := time.Since(start)
	log.Printf("共计写入总数为:%d %s", elapsed, ii, emo[rand.Intn(len(emo))])
	log.Printf("该功能函数执行完成耗时%v", elapsed)
}

func apicount(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	intcount, err := strconv.Atoi(string(body))
	if err != nil {
		log.Printf("string转ini总数失败:%v", err)
	}
	intpageSIze, err := strconv.Atoi(*pageSize)
	pageCount := int(math.Ceil(float64(intcount) / float64(intpageSIze))) //总页数
	return pageCount

}

//先通过excelize库创建一个手动创建的会报错
func createxml() {
	f := excelize.NewFile()
	//设置列宽度，A为序号调小一点
	f.SetColWidth("Sheet1", "A", "B", 8)
	f.SetColWidth("Sheet1", "B", "E", 15)
	f.SetColWidth("Sheet1", "F", "F", 30)
	f.SetColWidth("Sheet1", "G", "G", 150)
	index := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "序号")
	f.SetCellValue("Sheet1", "B1", "异常IP")
	f.SetCellValue("Sheet1", "C1", "异常类型")
	f.SetCellValue("Sheet1", "D1", "记录时间")
	f.SetCellValue("Sheet1", "E1", "异常IP报警次数")
	f.SetCellValue("Sheet1", "F1", "来源文件")
	f.SetCellValue("Sheet1", "G1", "异常内容")
	f.SetActiveSheet(index)
	if err := f.SaveAs(*infoname); err != nil {
		fmt.Println(err)
	}
}
