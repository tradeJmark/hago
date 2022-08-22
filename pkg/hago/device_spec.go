package hago

type DeviceSpec struct {
	Type string
	Properties map[string]string
	DeviceSpecific map[string]string `yaml:"device_specific"`
}

type LightSpec struct {
	StateTopic string `mapstructure:"state_topic"`
	CommandTopic string `mapstructure:"command_topic"`
	PayloadOn string `mapstructure:"payload_on"`
	PayloadOff string `mapstructure:"payload_off"`
}