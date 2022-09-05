package info

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	apiserver = flag.String("apiserver", "http://103.44.250.69:17171/api?page=1&pageSize=100000", "此参数是info模式下调用的接口数据")
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
	createxml()
	start := time.Now() // 获取当前时间
	resp, err := http.Get(*apiserver)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var info []Warn
	//解析json数据到结构体
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f, err := excelize.OpenFile("severwarn.xlsx")
	if err != nil {
		log.Println("打开的时候报错了")
		fmt.Println(err)
	}
	//西想xml文件追加数据
	for i, v := range info {
		i += 2
		log.Printf("正在向severwarn.xlsx追加第%d条告警数据", i-1)
		if len(v.IP) == 0 {
			return
		}
		//写入接口数据
		intnumber, _ := strconv.Atoi(strconv.Itoa(i))
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), intnumber-1)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), v.IP)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), v.TYPE)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), v.DATE)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), v.COUNT)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), v.PATH)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), v.INFO)
	}
	//保存
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Println("该功能函数执行完成耗时：", elapsed)
}

//先通过excelize库创建一个手动创建的会报错
func createxml() {
	f := excelize.NewFile()
	//设置列宽度，A为序号调小一点
	f.SetColWidth("Sheet1", "A", "B", 5)
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
	if err := f.SaveAs("severwarn.xlsx"); err != nil {
		fmt.Println(err)
	}
}
