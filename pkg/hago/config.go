package hago

type Config struct {
	Broker struct {
		Url string
		Port string
	}
	Devices []DeviceSpec `yaml:",flow"`
	Client struct {
		Name string
	}
}