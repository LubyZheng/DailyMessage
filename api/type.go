package api

type NewsResponse struct {
	Success string            `json:"reason"`
	Result  NewsResponseChild `json:"result"`
}

type NewsResponseChild struct {
	Data []News `json:"data"`
}

type News struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Weather struct {
	City     string
	ImgURL   string
	Temp     string
	Weather  string
	Air      string
	Humidity string
	Wind     string
	Note     string
}

type CovidResponse struct {
	Success bool                 `json:"success"`
	Results []CovidResponseChild `json:"results"`
}

type CovidResponseChild struct {
	Cities     []CityStatistics `json:"cities"`
	UpdateTime int64            `json:"updateTime"`
}

type CityStatistics struct {
	ExistingInfectedPopulation  int    `json:"currentConfirmedCount"`
	TotalInfectedPopulation     int    `json:"confirmedCount"`
	SuspectedInfectedPopulation int    `json:"suspectedCount"`
	CuredPopulation             int    `json:"curedCount"`
	DeadPopulation              int    `json:"deadCount"`
	CityEnglishName             string `json:"cityEnglishName"`
}

type Covid struct {
	ExistingInfectedPopulation  string
	TotalInfectedPopulation     string
	SuspectedInfectedPopulation string
	CuredPopulation             string
	DeadPopulation              string
	UpdateTime                  string
}
