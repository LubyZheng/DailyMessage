package template

import (
	"fmt"
	"reflect"
	"strings"
)

// HTML for email template
const HTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Daily Message</title>
</head>
<body>
<div style="max-width: 375px; margin: 20px auto;color:#444; font-size: 16px;">
    <div><h2 style="text-align: center">Daily Message</h2></div>
    <div style="margin-top: 10px;line-height: 1;">{{message.Content}}</div>
    <br>
    <hr>
    <!--    <h3 style="text-align: center">{{weather.City}}</h3>-->
    <div style="text-align: center;font-size: 18px;"><img width="15%" src="{{weather.ImgURL}}"> {{weather.Weather}}
    </div>
    <br>
    <div style="padding: 0;width: 100%;">
        <div><span style="color: #6e6e6e">温度：</span>{{weather.Temp}}</div>
        <div><span style="color: #6e6e6e">湿度：</span>{{weather.Humidity}}</div>
        <div><span style="color: #6e6e6e">风向：</span>{{weather.Wind}}</div>
        <div><span style="color: #6e6e6e">空气：</span>{{weather.Air}}</div>
        <div><span style="color: #6e6e6e">提示：</span>{{weather.Note}}</div>
    </div>
    <br>
    <hr>

    <div><h3 style="text-align: center;">Covid</h3></div>
    <div>
        <div><span style="color: #6e6e6e">数据更新时间：</span>{{covid.UpdateTime}}</div>
        <div><span style="color: #6e6e6e">现存确诊人数：</span>{{covid.ExistingInfectedPopulation}}</div>
        <div><span style="color: #6e6e6e">累计确诊人数：</span>{{covid.TotalInfectedPopulation}}</div>
        <div><span style="color: #6e6e6e">疑似感染人数：</span>{{covid.SuspectedInfectedPopulation}}</div>
        <div><span style="color: #6e6e6e">治愈人数：</span>{{covid.CuredPopulation}}</div>
        <div><span style="color: #6e6e6e">死亡人数：</span>{{covid.DeadPopulation}}</div>
    </div>
    <br>
    <hr>

    <div><h3 style="text-align: center;">News</h3></div>
    <div>
        <div><a href="{{news1.URL}}">1.{{news1.Title}}</a></div>
        <div><a href="{{news2.URL}}">2.{{news2.Title}}</a></div>
        <div><a href="{{news3.URL}}">3.{{news3.Title}}</a></div>
        <div><a href="{{news4.URL}}">4.{{news4.Title}}</a></div>
        <div><a href="{{news5.URL}}">5.{{news5.Title}}</a></div>
        <div><a href="{{news6.URL}}">6.{{news6.Title}}</a></div>
        <div><a href="{{news7.URL}}">7.{{news7.Title}}</a></div>
        <div><a href="{{news8.URL}}">8.{{news8.Title}}</a></div>
        <div><a href="{{news9.URL}}">9.{{news9.Title}}</a></div>
        <div><a href="{{news10.URL}}">10.{{news10.Title}}</a></div>
    </div>
    <hr>

</div>
<br><br>
</body>
</html>
`

func GenerateHTML(html string, info map[string]interface{}) string {
	for key, data := range info {
		rDataKey := reflect.TypeOf(data)
		rDataVal := reflect.ValueOf(data)
		fieldNum := rDataKey.NumField()
		for i := 0; i < fieldNum; i++ {
			fName := rDataKey.Field(i).Name
			rValue := rDataVal.Field(i)

			var fValue string
			switch rValue.Interface().(type) {
			case string:
				fValue = rValue.String()
			case []string:
				fValue = strings.Join(rValue.Interface().([]string), "<br>")
			}

			mark := fmt.Sprintf("{{%s.%s}}", key, fName)
			html = strings.ReplaceAll(html, mark, fValue)
		}
	}
	return html
}