package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	port        = flag.String("port", "17171", "此参数server模式下启动API占用的端口")
	errcount    = 0
	Warnlist    = []Warn{}
	IPcount     = make(map[string]int)
	execlnumber = 1 //写入execl的初始值
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

func Run() {
	flag.Parse()
	r := gin.Default()
	r.Use(CORS())
	r.GET("/warn", func(c *gin.Context) {
		execlnumber += 1
		text := c.Query("text")
		path := c.Query("path")
		ip := c.ClientIP()
		log.Printf("---------------------------------------------第%d条异常记录😡", errcount)
		log.Println("服务器:", ip, " 异常来源:", path)
		fmt.Println(text, "\n")
		errcount++
		now := time.Now() //当前时间
		IPcount[ip]++
		warntype := ""
		if find := strings.Contains(path, "secure"); find {
			warntype += "SSH"
		}
		if find := strings.Contains(path, "ping"); find {
			warntype += "PING"
		}
		if find := strings.Contains(path, "warnhistory"); find {
			warntype += "History"
		}
		if find := strings.Contains(path, "warnmd5"); find {
			warntype += "MD5"
		}
		if find := strings.Contains(path, "warnnginx"); find {
			warntype += "WEB"
		}
		data := Warn{
			IP:    ip,
			COUNT: strconv.Itoa(IPcount[ip]),
			PATH:  path,
			INFO:  text,
			DATE:  now.Format("2006/1/02 15:04"),
			TYPE:  warntype,
		}
		//如果需要开启API则开启此条
		Warnlist = append(Warnlist, data)
	})

	r.GET("api", func(c *gin.Context) {
		page := c.Query("page")
		pageSize := c.Query("pageSize")
		if len(page) == 0 || len(pageSize) == 0 {
			all := len(Warnlist)
			c.String(200, strconv.Itoa(all))
			return
		}
		//string转换int
		intpage, _ := strconv.Atoi(page)
		intpageSize, _ := strconv.Atoi(pageSize)
		sliceStart, sliceEnd := SlicePage(intpage, intpageSize, len(Warnlist))
		bytes, _ := json.Marshal(Warnlist[sliceStart:sliceEnd])
		c.String(200, string(bytes))
	})

	r.Run(":" + *port)
}

//切片分页
func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 20
	}
	if pageSize > nums {
		return 0, nums
	}
	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

// 开启跨域函数
func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 允许 Origin 字段中的域发送请求
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		// 设置预验请求有效期为 86400 秒
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 设置允许请求的方法
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		// 设置允许请求的 Header
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length，Apitoken")
		// 设置拿到除基本字段外的其他字段，如上面的Apitoken, 这里通过引用Access-Control-Expose-Headers，进行配置，效果是一样的。
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
		// 配置是否可以带认证信息
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// OPTIONS请求返回200
		if context.Request.Method == "OPTIONS" {
			fmt.Println(context.Request.Header)
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
