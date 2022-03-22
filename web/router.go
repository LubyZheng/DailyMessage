package web

import (
	"DailyMessage/email"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

var AllTask map[string]Task
var AllContent map[string]Content

func init() {
	AllTask = make(map[string]Task)
	AllContent = make(map[string]Content)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", login)
	r.POST("/updateInfo", updateInfo)
	r.POST("/startTask", startTask)
	r.POST("/getContent", getContent)
	r.DELETE("/deleteTask", deleteTask)
	return r
}

func login(ctx *gin.Context) {
	account := User{}
	err := ctx.BindJSON(&account)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	target := map[string]string{
		"MAIL_USERNAME": account.Account,
		"MAIL_PASSWORD": account.Password,
	}
	err = OverWriteEnvFile(target)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	response(ctx, true)
}

func updateInfo(ctx *gin.Context) {
	userInfo := UserInfo{}
	err := ctx.BindJSON(&userInfo)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	target := map[string]string{
		"MAIL_USERNAME":  userInfo.Account,
		"MAIL_PASSWORD":  userInfo.Password,
		"MAIL_SIGNATURE": userInfo.Signature,
	}
	err = OverWriteEnvFile(target)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	response(ctx, true)
}

func startTask(ctx *gin.Context) {
	task := Task{}
	err := ctx.BindJSON(&task)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	if _, ok := AllTask[task.EmailAddress]; ok {
		delete(AllTask, task.EmailAddress)
	}
	AllTask[task.EmailAddress] = task
	var t []Task
	for key := range AllTask {
		t = append(t, AllTask[key])
	}
	b, err := json.Marshal(t)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	target := map[string]string{
		"MAIL_TO": string(b),
	}
	err = OverWriteEnvFile(target)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	SetTasks(task)
	response(ctx, true)
}

func deleteTask(ctx *gin.Context) {
	task := Task{}
	err := ctx.BindJSON(&task)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	delete(AllTask, task.EmailAddress)
	var t []Task
	for key := range AllTask {
		t = append(t, AllTask[key])
	}
	b, err := json.Marshal(t)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	target := map[string]string{
		"MAIL_TO": string(b),
	}
	err = OverWriteEnvFile(target)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	RemoveTask(task)
	response(ctx, true)
}

func getContent(ctx *gin.Context) {
	content := Content{}
	err := ctx.BindJSON(&content)
	if err != nil {
		response(ctx, false)
		log.Println(err.Error())
		return
	}
	content_format := strings.ReplaceAll(email.EMAIL_FORMAT, "{{ReceiverName}}", content.ReceiverName)
	content_format = strings.ReplaceAll(content_format, "{{Content}}", content.Content)
	content_format = strings.ReplaceAll(content_format, "{{Signature}}", os.Getenv("MAIL_SIGNATURE"))
	AllContent[content.EmailAddress] = Content{
		Content: content_format,
	}
	response(ctx, true)
}

func response(ctx *gin.Context, success bool) {
	if success == true {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "true",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "false",
		})
	}
}
