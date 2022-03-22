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

type news_instance struct {
	Content string
}

func init() {
	JobList = make(map[string]cron.EntryID)
	nyc, _ = time.LoadLocation("Asia/Shanghai")
	cJob = cron.New(cron.WithLocation(nyc))
}

func SetTasks(task Task) {
	mtx.Lock()
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
	mtx.Unlock()
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
	info["message"] = AllContent[task.EmailAddress] // 邮件内容
	var err error                                   //error暂不处理
	for name, apiFunc := range apiFuncs {
		wg.Add(1)
		go func(key string, fn func() (interface{}, error)) {
			defer wg.Done()
			info[key], err = fn()
		}(name, apiFunc)
	}
	wg.Wait()
	//将News数组分割成若干结构体，否则调用GenerateHTML接口会panic
	//hard code，赶时间，不是特别好
	var news_content string
	if info["news"] != nil {
		for i := 0; i < 10; i++ {
			news_content += fmt.Sprintf("<div><a href=\"%s\">%d.%s</a></div>\n",
				(info["news"]).([]api.News)[i].URL, i+1, (info["news"]).([]api.News)[i].Title)
		}
		info["news"] = news_instance{
			news_content,
		}
	} else {
		info["news"] = news_instance{
			"News switch is off",
		}
	}
	email.SendMail(template.GenerateHTML(template.HTML, info), task.EmailAddress)
}
