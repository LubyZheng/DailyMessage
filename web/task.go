package web

import (
	"DailyMessage/api"
	"DailyMessage/email"
	"DailyMessage/template"
	"fmt"
	env "github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

var mtx sync.Mutex
var nyc *time.Location
var cJob *cron.Cron
var JobList map[string]cron.EntryID

type HTML_Content struct {
	Content string
}

func init() {
	JobList = make(map[string]cron.EntryID)
	nyc, _ = time.LoadLocation("Asia/Shanghai")
	cJob = cron.New(cron.WithLocation(nyc))
}

func SetTasks(task Task) {
	err := env.Overload()
	if err != nil {
		log.Fatalf("Load .env file error: %s", err)
	}
	scheduleTime := fmt.Sprintf("%s %s * * *", task.Time.Minute, task.Time.Hour)
	// 如果有一个相同的job但是时间不同，先删除之前的job再添加
	RemoveTask(task)
	ID, err := cJob.AddFunc(scheduleTime, func() {
		taskJob(task)
	})
	JobList[task.EmailAddress] = ID
	if err != nil {
		return
	}
	cJob.Start()
}

func RemoveTask(task Task) {
	if _, ok := JobList[task.EmailAddress]; ok {
		cJob.Remove(JobList[task.EmailAddress])
		delete(JobList, task.EmailAddress)
		if len(JobList) == 0 {
			cJob.Stop()
		}
	} else {
		return
	}
}

func taskJob(task Task) {
	mtx.Lock()
	apiFuncs := make(map[string]func() (interface{}, error))
	if task.Weather == true {
		apiFuncs["weather"] = func() (interface{}, error) {
			return api.GetWeather(task.Location.Country, task.Location.Province, task.Location.City)
		}
	}
	if task.Covid == true {
		apiFuncs["covid"] = func() (interface{}, error) {
			return api.GetCovidStatistics(task.Location.Province, task.Location.City)
		}
	}
	if task.News == true {
		apiFuncs["news"] = func() (interface{}, error) {
			return api.GetNews()
		}
	}
	wg := sync.WaitGroup{}
	info := make(map[string]interface{})
	var err error //error暂不处理
	for name, apiFunc := range apiFuncs {
		wg.Add(1)
		go func(key string, fn func() (interface{}, error)) {
			defer wg.Done()
			info[key], err = fn()
		}(name, apiFunc)
	}
	wg.Wait()
	mtx.Unlock()
	info["message"] = AllContent[task.EmailAddress] // 邮件内容
	if _, ok := AllContent[task.EmailAddress]; ok {
		if AllContent[task.EmailAddress].Content != "" {
			info["message"] = AllContent[task.EmailAddress]
			delete(AllContent, task.EmailAddress) //获取完当天的message后就从map中删除，避免第二天获取相同的message
		}
	} else {
		info["message"] = HTML_Content{
			"There is no message today",
		}
	}
	if info["weather"] != nil {
		weather_temp := (info["weather"]).(api.Weather)
		wea_cont := fmt.Sprintf(`<div style="text-align: center;font-size: 18px;"><img width="15%%" src="%s">%s
		</div>
		<br>
		<div style="padding: 0;width: 100%%;">
		<div><span style="color: #6e6e6e">温度：</span>%s</div>
		<div><span style="color: #6e6e6e">湿度：</span>%s</div>
		<div><span style="color: #6e6e6e">风向：</span>%s</div>
		<div><span style="color: #6e6e6e">空气：</span>%s</div>
		<div><span style="color: #6e6e6e">提示：</span>%s</div>
		</div>
		<br>`,
			weather_temp.ImgURL,
			weather_temp.Weather,
			weather_temp.Temp,
			weather_temp.Humidity,
			weather_temp.Wind,
			weather_temp.Air,
			weather_temp.Note)
		info["weather"] = HTML_Content{
			wea_cont,
		}
	} else {
		info["weather"] = HTML_Content{
			"Weather switch is off",
		}
	}
	if info["covid"] != nil {
		covid_temp := (info["covid"]).(api.Covid)
		covid_cont := fmt.Sprintf(`    
		<div>
			<div><span style="color: #6e6e6e">数据更新时间：</span>%s</div>
			<div><span style="color: #6e6e6e">现存确诊人数：</span>%s</div>
			<div><span style="color: #6e6e6e">累计确诊人数：</span>%s</div>
			<div><span style="color: #6e6e6e">疑似感染人数：</span>%s</div>
			<div><span style="color: #6e6e6e">治愈人数：</span>%s</div>
			<div><span style="color: #6e6e6e">死亡人数：</span>%s</div>
		</div>`,
			covid_temp.UpdateTime,
			covid_temp.ExistingInfectedPopulation,
			covid_temp.TotalInfectedPopulation,
			covid_temp.SuspectedInfectedPopulation,
			covid_temp.CuredPopulation,
			covid_temp.DeadPopulation)
		info["covid"] = HTML_Content{
			covid_cont,
		}
	} else {
		info["covid"] = HTML_Content{
			"Covid switch is off",
		}
	}
	if info["news"] != nil {
		var news_cont string
		for i := 0; i < 10; i++ {
			news_cont += fmt.Sprintf("<div><a href=\"%s\">%d.%s</a></div>\n",
				(info["news"]).([]api.News)[i].URL, i+1, (info["news"]).([]api.News)[i].Title)
		}
		info["news"] = HTML_Content{
			news_cont,
		}
	} else {
		info["news"] = HTML_Content{
			"News switch is off",
		}
	}
	email.SendMail(template.GenerateHTML(template.HTML, info), task.EmailAddress)
}
