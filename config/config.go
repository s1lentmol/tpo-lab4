package config

const (
	Token   = "518499974"
	User    = "-1355033797"
	BaseURL = "http://127.0.0.1:8888"
)

type HardwareConfig struct {
	ID    int
	Price int
}

var HardwareConfigs = []HardwareConfig{
	{ID: 1, Price: 2700},
	{ID: 2, Price: 2900},
	{ID: 3, Price: 5400},
}