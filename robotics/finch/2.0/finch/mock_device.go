package finch

import (
	"fmt"
)

type MockDevice struct {
	name     string
	address  string
	resolved bool
}

func (m *MockDevice) LocalName() string {
	return m.name
}

func (m *MockDevice) Address() string {
	return m.address
}

func (m *MockDevice) Connect() error {
	fmt.Printf("[MockDevice] Connecting to %s\n", m.name)
	return nil
}

func (m *MockDevice) Disconnect() error {
	fmt.Printf("[MockDevice] Disconnecting from %s\n", m.name)
	return nil
}

func (m *MockDevice) ServicesResolved() bool {
	return m.resolved
}

func (m *MockDevice) Services() ([]Service, error) {
	uartService := &MockService{
		uuid: UARTServiceUUID,
		characteristics: []Characteristic{
			&MockCharacteristic{uuid: UARTRXCharUUID},
			&MockCharacteristic{uuid: UARTTXCharUUID},
		},
	}
	return []Service{uartService}, nil
}

type MockService struct {
	uuid            string
	characteristics []Characteristic
}

func (m *MockService) UUID() string {
	return m.uuid
}

func (m *MockService) Characteristics() ([]Characteristic, error) {
	return m.characteristics, nil
}

type MockCharacteristic struct {
	uuid string
}

func (m *MockCharacteristic) UUID() string {
	return m.uuid
}

func (m *MockCharacteristic) Write(data []byte) error {
	fmt.Printf("[MockCharacteristic] Writing to %s: %v\n", m.uuid, data)
	return nil
}

func (m *MockCharacteristic) Read() ([]byte, error) {
	mockData := make([]byte, 20)
	mockData[0] = 1   // sequence counter
	mockData[1] = 50  // distance sensor
	mockData[2] = 100 // left light sensor
	mockData[3] = 150 // right light sensor
	mockData[4] = 13  // accelerometer X
	mockData[5] = 11  // accelerometer Y
	mockData[6] = 67  // accelerometer Z
	mockData[7] = 0
	mockData[8] = 1   // magnetometer flag
	mockData[9] = 200 // magnetometer value
	mockData[10] = 0
	mockData[11] = 1   // magnetometer flag
	mockData[12] = 180 // magnetometer value
	mockData[13] = 0   // line sensor
	mockData[14] = 216 // encoder left high
	mockData[15] = 200 // encoder right high
	mockData[16] = 50  // always 50
	mockData[17] = 8   // encoder left low
	mockData[18] = 9   // encoder right low
	mockData[19] = 129 // always 129

	return mockData, nil
}

func (m *MockCharacteristic) Notify(enable bool) error {
	fmt.Printf("[MockCharacteristic] Notifications %s for %s\n",
		map[bool]string{true: "enabled", false: "disabled"}[enable], m.uuid)
	return nil
}
