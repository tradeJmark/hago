package hago

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
)

func Check(e error) {
  if e != nil {
    panic(e)
  }
}

func GetBroker(clientName string, url string, port string) mqtt.Client {
	options := mqtt.NewClientOptions()
	options.AddBroker("tcp://" + url + ":" + port)
	options.SetOrderMatters(false)
	options.SetClientID("HAgo - " + clientName)
	return mqtt.NewClient(options)
}

func BuildController(spec DeviceSpec) Controller {
	switch spec.Type {
	case "mock-light":
		lightSpec := LightSpec{}
		mapstructure.Decode(spec.Properties, &lightSpec)
		return NewMockLight(&lightSpec)
	case "serial-light":
		lightSpec := LightSpec{}
		mapstructure.Decode(spec.Properties, &lightSpec)
		return NewSerialLight(&lightSpec, spec.DeviceSpecific)
	default:
		return nil
	}
}