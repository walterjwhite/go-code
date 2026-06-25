package finch

import (
	"math"
)

type FinchRobot struct {
	bluetooth  *FinchBluetooth
	lightState [15]byte
}

func NewFinchRobot(bluetooth *FinchBluetooth) *FinchRobot {
	return &FinchRobot{
		bluetooth: bluetooth,
	}
}

func (fr *FinchRobot) Close() error {
	return fr.bluetooth.Close()
}


func (fr *FinchRobot) SetBeak(r, g, b int) error {
	fr.lightState[0] = byte(clamp(r))
	fr.lightState[1] = byte(clamp(g))
	fr.lightState[2] = byte(clamp(b))
	return fr.sendLights(0, 0)
}

func (fr *FinchRobot) SetTailLight(port, r, g, b int) error {
	if port == 0 {
		for i := 0; i < 4; i++ {
			fr.lightState[3+i*3] = byte(clamp(r))
			fr.lightState[3+i*3+1] = byte(clamp(g))
			fr.lightState[3+i*3+2] = byte(clamp(b))
		}
	} else {
		offset := 3 + (port-1)*3
		fr.lightState[offset] = byte(clamp(r))
		fr.lightState[offset+1] = byte(clamp(g))
		fr.lightState[offset+2] = byte(clamp(b))
	}
	return fr.sendLights(0, 0)
}

func (fr *FinchRobot) SetLightsOff() error {
	for i := range fr.lightState {
		fr.lightState[i] = 0
	}
	return fr.sendLights(0, 0)
}


func (fr *FinchRobot) PlayNote(midiNote, durationMs int) error {
	freq := 440.0 * math.Pow(2.0, float64(midiNote-69)/12.0)
	period := int(math.Round(1000000.0 / freq))
	return fr.sendLights(period, durationMs)
}


func (fr *FinchRobot) SetMotors(leftSpeed int, leftForward bool, rightSpeed int, rightForward bool) error {
	cmd := make([]byte, 20)
	cmd[0] = 0xD2
	cmd[1] = 0x40
	cmd[2] = encodeSpeed(leftSpeed, leftForward)
	cmd[6] = encodeSpeed(rightSpeed, rightForward)

	return fr.bluetooth.WriteRx(cmd)
}

func (fr *FinchRobot) Forward(speed int) error {
	return fr.SetMotors(speed, true, speed, true)
}

func (fr *FinchRobot) Backward(speed int) error {
	return fr.SetMotors(speed, false, speed, false)
}

func (fr *FinchRobot) TurnLeft(speed int) error {
	return fr.SetMotors(speed, false, speed, true)
}

func (fr *FinchRobot) TurnRight(speed int) error {
	return fr.SetMotors(speed, true, speed, false)
}

func (fr *FinchRobot) StopMotors() error {
	return fr.SetMotors(0, true, 0, true)
}

func (fr *FinchRobot) StopAll() error {
	return fr.bluetooth.WriteRx([]byte{0x62, 0x70})
}

func (fr *FinchRobot) ResetEncoders() error {
	return fr.bluetooth.WriteRx([]byte{0xD4})
}


func (fr *FinchRobot) RequestSensors() error {
	return fr.bluetooth.WriteRx([]byte{0x01})
}

func (fr *FinchRobot) ReadRawSensors() ([]byte, error) {
	return fr.bluetooth.Read()
}

func (fr *FinchRobot) ReadDistanceCm() (int, error) {
	data, err := fr.ReadRawSensors()
	if err != nil {
		return 0, err
	}
	return int(data[1]), nil
}

func (fr *FinchRobot) ReadLightSensors() ([]int, error) {
	data, err := fr.ReadRawSensors()
	if err != nil {
		return nil, err
	}
	return []int{int(data[2]), int(data[3])}, nil
}

func (fr *FinchRobot) ReadAccelerometer() ([]int, error) {
	data, err := fr.ReadRawSensors()
	if err != nil {
		return nil, err
	}
	return []int{
		int(int8(data[4])),
		int(int8(data[5])),
		int(int8(data[6])),
	}, nil
}


func (fr *FinchRobot) sendLights(buzzPeriod, buzzDurationMs int) error {
	cmd := make([]byte, 20)
	cmd[0] = 0xD0
	copy(cmd[1:16], fr.lightState[:])
	cmd[16] = byte((buzzPeriod >> 8) & 0xFF)
	cmd[17] = byte(buzzPeriod & 0xFF)
	cmd[18] = byte((buzzDurationMs >> 8) & 0xFF)
	cmd[19] = byte(buzzDurationMs & 0xFF)

	return fr.bluetooth.WriteRx(cmd)
}

func encodeSpeed(speed int, forward bool) byte {
	s := max(0, min(100, speed))
	if forward {
		return byte(s | 0x80)
	}
	return byte(s)
}

func clamp(value int) int {
	return max(0, min(255, value))
}
