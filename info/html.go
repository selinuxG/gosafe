package info

import (
	"github.com/go-echarts/go-echarts/charts"
	"log"
	"os"
)

func Runhtml(sshfa map[string]interface{},
	sshac map[string]interface{},
	nginx map[string]interface{},
	hist map[string]interface{}) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.TitleOpts{Title: "登陆失败IP:"})
	pie.Add("登陆失败IP:", sshfa,
		charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"},
		charts.PieOpts{Radius: []string{"40%", "75%"}},
	)

	sshacp := charts.NewPie()
	sshacp.SetGlobalOptions(charts.TitleOpts{Title: "登陆成功IP:"})
	sshacp.Add("登陆成功IP:", sshac,
		charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"},
		charts.PieOpts{Radius: []string{"40%", "75%"}},
	)
	nginxht := charts.NewPie()
	nginxht.SetGlobalOptions(charts.TitleOpts{Title: "WEB访问监控:"})
	nginxht.Add("状态码:", nginx,
		charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"},
		charts.PieOpts{Radius: []string{"40%", "75%"}},
	)

	history := charts.NewPie()
	history.SetGlobalOptions(charts.TitleOpts{Title: "异常命令次数:"})
	history.Add("异常命令次数:", hist,
		charts.LabelTextOpts{Show: true, Formatter: "{b}: {c}"},
		charts.PieOpts{Radius: []string{"40%", "75%"}})

	f, err := os.Create("gasafe.html")
	if err != nil {
		log.Println(err)
	}
	sshacp.Render(f)
	pie.Render(f)
	nginxht.Render(f)
	history.Render(f)

}
