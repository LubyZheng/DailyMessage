# EE5415 Simple Server

### 效果图

![](https://github.com/LubyZheng/DailyMessage/raw/master/pic/email.png)

### 服务端暂时只支持QQ邮箱发送，需要先在QQ邮箱开通IMAP service

开启路径: 设置->账户

![](https://github.com/LubyZheng/DailyMessage/raw/master/pic/imap.png)



### 进程示意图，MAC用户使用根目录下的mac.exe, WIN用户用win.exe

![](https://github.com/LubyZheng/DailyMessage/raw/master/pic/exe.png)



### ENV环境文件

服务端默认端口号是8080，可以自行修改

![](https://github.com/LubyZheng/DailyMessage/raw/master/pic/env.png)



### 路由

#### 127.0.0.1:8080/login

邮箱账户和IMAP密码会填入env文件的MAIL_USERNAME和MAIL_PASSWORD

```html
{
  "account": "your email account",    
  "password": "your IMAP password"   
}
```



#### 127.0.0.1:8080/updateInfo

更新env文件

邮箱账户,IMAP密码和签名会填入env文件的MAIL_USERNAME，MAIL_PASSWORD和MAIL_SIGNATURE

```html
{
  "account": "your email account",   
  "password": "your IMAP password"    
  "signature": "XXX"                 
}
```



#### 127.0.0.1:8080/startTask

暂时只支持一个receiver，如果两个的话会报错，后续优化，不影响调试和客户端代码，填入env文件的MAIL_TO

```html
{
    "receiverName": "XXX",              
    "emailAddress": "receiver's email", 
    "location": {
        "country": "china",            #例子，某些省份或城市暂时不支持，例如香港，后续有必要再优化
        "province": "guangdong",
        "city": "shenzhen"
    },
    "time": {
        "hour": "19",              #注意此处必须是两位数，例如早上九点应该为09，分钟同理
        "minute": "00"
    },
    "weather": true,               #bool类型，后台暂时没有添加判断机制，后续优化
    "news": true,
    "covid": true
}
```



#### 127.0.0.1:8080/getContent

```html
{
    "receiverName": "XXX",
    "emailAddress": "receiver's email",
    "content": "XXX"             
}
```



#### 127.0.0.1:8080/deleteTask

逻辑暂时还没完善，后续优化，不影响调试和客户端代码

```html
{
    "receiverName": XXX",
    "emailAddress": "receiver's email"
}
```



### Server响应

```
{
    "success": "true 成功/ false 失败"
}
```



reference: 

https://github.com/BarryYan/daily-warm
