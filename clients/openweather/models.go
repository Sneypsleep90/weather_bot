package openweather

type CoordinatesResponse struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}
type Coordinates struct {
	Lat float64
	Lon float64
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type Weather struct {
	Temp float64
}
