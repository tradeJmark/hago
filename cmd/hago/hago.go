package main

import (
  "os"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "github.com/tradejmark/hago/pkg/hago"
  "fmt"
)

func main() {
  config := hago.Config{}

  configFile := os.Args[1]
  configData, err := ioutil.ReadFile(configFile)
  hago.Check(err)
  
  err = yaml.Unmarshal(configData, &config)
  hago.Check(err)

  fmt.Println(config)

  broker := hago.GetBroker(config.Client.Name, config.Broker.Url, config.Broker.Port)
  token := broker.Connect()
  token.Wait()
  if token.Error() != nil {
    panic(token.Error())
  }
  fmt.Println("connected")
  for _, dev := range config.Devices {
    controller := hago.BuildController(dev)
    controller.MakeSubscriptions(broker)
    go controller.Loop(broker)
  }
  for {}
}
