package api

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	WEATHER_URL = "https://tianqi.moji.com/weather"
	COVID_URL   = "https://lab.isaaclin.cn/nCoV/api/area?latest=1&provinceEng="
	NEWS_URL    = "http://v.juhe.cn/toutiao/index?type=guonei&key=d268884b9b07c0eb9d6093dc54116018"
)

func GetWeather(country, province, city string) (Weather, error) {
	url := WEATHER_URL + fmt.Sprintf("/%s/%s/%s",
		strings.ToLower(country), strings.ToLower(province), strings.ToLower(city)) //小写
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return Weather{}, fmt.Errorf("[Weather API error]"+
			"status code:%d, error information:%s", resp.StatusCode, err.Error())
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	imgURL, _ := doc.Find(".forecast .days li img").Attr("src")
	humidityDesc := strings.Split(doc.Find(".wea_about span").Text(), " ")
	humidity := "unknown"
	if len(humidityDesc) >= 2 {
		humidity = humidityDesc[1]
	}
	return Weather{
		City:     doc.Find("#search .search_default em").Text(),
		ImgURL:   imgURL,
		Temp:     doc.Find(".forecast .days li:nth-child(3)").Eq(0).Text(),
		Weather:  doc.Find(".forecast .days li:nth-child(2)").Eq(0).Text(),
		Air:      doc.Find(".wea_alert em").Text(),
		Humidity: humidity,
		Wind:     doc.Find(".wea_about em").Text(),
		Note:     strings.ReplaceAll(doc.Find(".wea_tips em").Text(), "。", ""),
	}, nil
}

func GetCovidStatistics(province, city string) (Covid, error) {
	province = strings.ToUpper(province[:1]) + strings.ToLower(province[1:]) //首字母大写，接口调用需要
	url := COVID_URL + province
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return Covid{}, fmt.Errorf("[Covid API error]"+
			"status code:%d, error information:%s", resp.StatusCode, err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Covid{}, fmt.Errorf("[Covid API error] error information:%s", err.Error())
	}
	defer resp.Body.Close()
	var CovidJson CovidResponse
	err = json.Unmarshal(body, &CovidJson)
	if err != nil {
		return Covid{}, fmt.Errorf("[Covid API error] error information:%s", err.Error())
	}
	if CovidJson.Success == false {
		return Covid{}, fmt.Errorf("[Covid API error] error information:unknown error")
	}
	pos := 0
	for ; pos < len(CovidJson.Results[0].Cities); pos++ {
		cityEngName := strings.ToUpper(city[:1]) + strings.ToLower(city[1:])
		if CovidJson.Results[0].Cities[pos].CityEnglishName == cityEngName {
			break
		}
	}
	if pos == len(CovidJson.Results[0].Cities) {
		return Covid{}, fmt.Errorf("[Covid API error] This city has no covid data")
	}
	return Covid{
		ExistingInfectedPopulation:  strconv.Itoa(CovidJson.Results[0].Cities[pos].ExistingInfectedPopulation),
		TotalInfectedPopulation:     strconv.Itoa(CovidJson.Results[0].Cities[pos].TotalInfectedPopulation),
		SuspectedInfectedPopulation: strconv.Itoa(CovidJson.Results[0].Cities[pos].SuspectedInfectedPopulation),
		CuredPopulation:             strconv.Itoa(CovidJson.Results[0].Cities[pos].CuredPopulation),
		DeadPopulation:              strconv.Itoa(CovidJson.Results[0].Cities[pos].DeadPopulation),
		UpdateTime: time.Unix(0, CovidJson.Results[0].UpdateTime*int64(time.Millisecond)).
			Format("2006-01-02 15:04:05"),
	}, nil
}

func GetNews() ([]News, error) {
	resp, err := http.Get(NEWS_URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[News API error]"+
			"status code:%d, error information:%s", resp.StatusCode, err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[News API error] error information:%s", err.Error())
	}
	defer resp.Body.Close()
	var NewsJson NewsResponse
	err = json.Unmarshal(body, &NewsJson)
	if err != nil {
		return nil, fmt.Errorf("[News API error] error information:%s", err.Error())
	}
	if NewsJson.Success != "success!" {
		//不成功的原因是当天接口的免费次数用完了，先用hard code代替，之后优化
		return []News{
			{"理解了这些，你就明白实现社会面动态清零要走几步？", "https://mini.eastday.com/mobile/220314143506041207595.html"},
			{"包钢500万吨球团脱硫项目发生火灾致7人死亡 事故原因正在调查", "https://mini.eastday.com/mobile/220314143504637837928.html"},
			{"天暖情更暖 双休假日公交员工日行一善暖人心", "https://mini.eastday.com/mobile/220314143453027545312.html"},
			{"95后和同事在双休日成“战友”，快递小哥变身“大白”，这个区的年轻人站到疫情防控志愿服务一线", "https://mini.eastday.com/mobile/220314143451640536224.html"},
			{"这个保温杯里泡着价值77万元的东西……", "https://mini.eastday.com/mobile/220314143437690922937.html"},
			{"为救助小动物 杭州这个学计算机专业的男孩自学考取兽医证", "https://mini.eastday.com/mobile/220314143012559857674.html"},
			{"不要出租它，否则成“帮凶", "https://mini.eastday.com/mobile/220314143007390127836.html"},
			{"汛期到来前逐户检修 北新桥上千个院落", "https://mini.eastday.com/mobile/220314143005926962955.html"},
			{"止损93万！他智扮演受害人“爸爸”四句话把骗子气“吐血”", "https://mini.eastday.com/mobile/220314143004281309717.html"},
			{"沪上春苔「图」", "https://mini.eastday.com/mobile/220314142546410174214.html"},
		}, nil
		//return nil, fmt.Errorf("[News API error] error information:unknown error")
	}
	n := make([]News, 10)
	for i := 0; i < 10; i++ {
		n[i] = News{
			Title: NewsJson.Result.Data[i].Title,
			URL:   NewsJson.Result.Data[i].URL,
		}
	}
	return n, nil
}
