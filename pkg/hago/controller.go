package hago

import (
	"github.com/eclipse/paho.mqtt.golang"
	"fmt"
	"go.bug.st/serial"
	"io"
	"bufio"
)

type Controller interface {
	Loop(mqtt.Client)
	MakeSubscriptions(mqtt.Client)
}

type LightController interface {
	SetState(bool)
	UpdateState()
	Spec() *LightSpec
}

type MockLight struct {
	State bool
	LightSpec *LightSpec
	Updates chan string
}

func NewMockLight(spec *LightSpec) *MockLight {
	ml := MockLight{State: false, LightSpec: spec, Updates: make(chan string)}
	return &ml
}

func (ml *MockLight) Spec() *LightSpec {
	return ml.LightSpec
}

func (ml *MockLight) SetState(state bool) {
	ml.State = state
	var msg string
	if state {
		msg = ml.LightSpec.PayloadOn
	} else {
		msg = ml.LightSpec.PayloadOff
	}
	ml.Updates <- string(msg)
}

func (ml *MockLight) UpdateState() {
}

func makeLCSubscriptions(lc LightController, broker mqtt.Client) {
	token := broker.Subscribe(lc.Spec().CommandTopic, 0, getLCCommandHandler(lc))
	token.Wait()
	fmt.Println("subscribed to " + lc.Spec().CommandTopic)
}

func getLCCommandHandler(lc LightController) mqtt.MessageHandler {
	return func(_ mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) == lc.Spec().PayloadOn {
			lc.SetState(true)
		} else if string(msg.Payload()) == lc.Spec().PayloadOff {
			lc.SetState(false)
		}
	}
}

func (ml *MockLight) MakeSubscriptions(broker mqtt.Client) {
	makeLCSubscriptions(ml, broker)
}

func (ml *MockLight) Loop(broker mqtt.Client) {
	var msg string
	if ml.State {
		msg = "ON"
	} else {
		msg = "OFF"
	}
	publishLCStateMsg(ml, broker, msg)
	for state := range ml.Updates {
		fmt.Println(state)
		publishLCStateMsg(ml, broker, state)
	}
}

func publishLCStateMsg(lc LightController, broker mqtt.Client, msg string) {
	token := broker.Publish(lc.Spec().StateTopic, 0, true, msg)
	token.Wait()
}











type SerialLight struct {
	LightSpec *LightSpec
	Device string
	StatusCommand string
	SerialPort serial.Port
}

func NewSerialLight(spec *LightSpec, deviceSpecific map[string]string) *SerialLight {
	sl := SerialLight{
		LightSpec: spec,
		Device: deviceSpecific["device"],
		StatusCommand: deviceSpecific["status_command"],
	}
	mode := &serial.Mode{BaudRate: 115200}
	port, err := serial.Open(sl.Device, mode)
	Check(err)
	sl.SerialPort = port
	return &sl
}

func (sl *SerialLight) Spec() *LightSpec {
	return sl.LightSpec
}

func (sl *SerialLight) SetState(state bool) {
	var msg string
	if state {
		msg = sl.LightSpec.PayloadOn
	} else {
		msg = sl.LightSpec.PayloadOff
	}
	io.WriteString(sl.SerialPort, msg + "\r\n")
}

func (sl *SerialLight) UpdateState() {
	io.WriteString(sl.SerialPort, sl.StatusCommand)
}

func (sl *SerialLight) MakeSubscriptions(broker mqtt.Client) {
	makeLCSubscriptions(sl, broker)
}

func (sl *SerialLight) Loop(broker mqtt.Client) {
	scanner := bufio.NewScanner(sl.SerialPort)
	for scanner.Scan() {
		publishLCStateMsg(sl, broker, scanner.Text())
	}
}