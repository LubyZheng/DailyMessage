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

func SetTasks(task Task) {
	err := env.Overload()
	if err != nil {
		log.Fatalf("Load .env file error: %s", err)
	}
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	cJob := cron.New(cron.WithLocation(nyc))
	scheduleTime := fmt.Sprintf("%s %s * * *", task.Time.Minute, task.Time.Hour)
	_, err = cJob.AddFunc(scheduleTime, func() {
		taskJob(task)
	})
	if err != nil {
		return
	}
	cJob.Start()
	select {}
}

func taskJob(task Task) {
	apiFuncs := map[string]func() (interface{}, error){
		"weather": func() (interface{}, error) {
			return api.GetWeather(task.Location.Country, task.Location.Province, task.Location.City)
		},
		"covid": func() (interface{}, error) {
			return api.GetCovidStatistics(task.Location.Province, task.Location.City)
		},
		"news": func() (interface{}, error) {
			return api.GetNews()
		},
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
	if info["news"] != nil {
		for i := 0; i < 10; i++ {
			info[fmt.Sprintf("news%d", i+1)] = api.News{
				Title: (info["news"]).([]api.News)[i].Title, URL: (info["news"]).([]api.News)[i].URL,
			}
		}
		delete(info, "news")
	}
	email.SendMail(template.GenerateHTML(template.HTML, info), task.EmailAddress)
}
