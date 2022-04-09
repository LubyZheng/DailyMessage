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
        {{weather.Content}}
    <hr>

    <div><h3 style="text-align: center;">COVID</h3></div>
        {{covid.Content}}
    <br>
    <hr>

    <div><h3 style="text-align: center;">NEWS</h3></div>
    <div>
        {{news.Content}}
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
