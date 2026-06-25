package finch

import (
	"fmt"
	"time"
)

const (
	UARTServiceUUID = "6e400001-b5a3-f393-e0a9-e50e24dcca9e"
	UARTRXCharUUID  = "6e400002-b5a3-f393-e0a9-e50e24dcca9e"
	UARTTXCharUUID  = "6e400003-b5a3-f393-e0a9-e50e24dcca9e"
)

type FinchBluetooth struct {
	device    Device
	rxChar    Characteristic
	txChar    Characteristic
	dataReady chan []byte
}

type FinchBluetoothBuilder struct{}

type Device interface {
	LocalName() string
	Address() string
	Connect() error
	Disconnect() error
	ServicesResolved() bool
	Services() ([]Service, error)
}

type Service interface {
	UUID() string
	Characteristics() ([]Characteristic, error)
}

type Characteristic interface {
	UUID() string
	Write([]byte) error
	Read() ([]byte, error)
	Notify(bool) error
}

func (b *FinchBluetoothBuilder) ConnectByAddress(macAddress string) (*FinchBluetooth, error) {
	fmt.Printf("[FinchBLE] Starting device scan for %s...\n", macAddress)

	device := &MockDevice{
		name:     "Finch Robot",
		address:  macAddress,
		resolved: true,
	}

	fmt.Printf("[FinchBLE] Found: %s (%s)\n", device.LocalName(), device.Address())

	return b.connect(device)
}

func (b *FinchBluetoothBuilder) connect(device Device) (*FinchBluetooth, error) {
	fmt.Printf("[FinchBLE] Connecting to %s...\n", device.LocalName())

	err := device.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	deadline := time.Now().Add(8 * time.Second)
	for !device.ServicesResolved() && time.Now().Before(deadline) {
		time.Sleep(200 * time.Millisecond)
	}
	if !device.ServicesResolved() {
		return nil, fmt.Errorf("service discovery timed out")
	}

	uartService, err := b.findUARTService(device)
	if err != nil {
		return nil, err
	}

	rxChar, txChar, err := b.findCharacteristics(uartService)
	if err != nil {
		return nil, err
	}

	fmt.Println("[FinchBLE] Enabling notifications on TX characteristic...")
	err = b.enableNotifications(txChar)
	if err != nil {
		return nil, fmt.Errorf("failed to enable notifications: %w", err)
	}

	err = b.sendInitSequence(rxChar)
	if err != nil {
		return nil, fmt.Errorf("failed to send init sequence: %w", err)
	}

	finch := &FinchBluetooth{
		device:    device,
		rxChar:    rxChar,
		txChar:    txChar,
		dataReady: make(chan []byte, 10),
	}

	go finch.monitorNotifications()

	fmt.Println("[FinchBLE] Connected and ready.")
	return finch, nil
}

func (b *FinchBluetoothBuilder) findUARTService(device Device) (Service, error) {
	services, err := device.Services()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	for _, service := range services {
		if service.UUID() == UARTServiceUUID {
			return service, nil
		}
	}

	return nil, fmt.Errorf("UART service not found")
}

func (b *FinchBluetoothBuilder) findCharacteristics(service Service) (Characteristic, Characteristic, error) {
	characteristics, err := service.Characteristics()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get characteristics: %w", err)
	}

	var rxChar, txChar Characteristic
	for _, char := range characteristics {
		uuid := char.UUID()
		if len(uuid) >= 8 && uuid[:8] == "6e400002" {
			rxChar = char
		} else if len(uuid) >= 8 && uuid[:8] == "6e400003" {
			txChar = char
		}
	}

	if rxChar == nil || txChar == nil {
		return nil, nil, fmt.Errorf("could not find RX/TX characteristics in UART service")
	}

	return rxChar, txChar, nil
}

func (b *FinchBluetoothBuilder) enableNotifications(txChar Characteristic) error {
	for i := 0; i < 5; i++ {
		err := txChar.Notify(true)
		if err == nil {
			fmt.Println("[FinchBLE] Notifications enabled.")
			return nil
		}

		fmt.Printf("[FinchBLE] startNotify() attempt %d did not take, retrying...\n", i+1)
		time.Sleep(500 * time.Millisecond)

		if i == 4 {
			return fmt.Errorf("failed to enable notifications on TX characteristic after 5 attempts")
		}
	}
	return nil
}

func (b *FinchBluetoothBuilder) sendInitSequence(rxChar Characteristic) error {
	fmt.Println("[FinchBLE] Sending init sequence (D4 + stopAll)...")

	err := rxChar.Write([]byte{0xD4})
	if err != nil {
		return fmt.Errorf("failed to send D4: %w", err)
	}
	time.Sleep(150 * time.Millisecond)

	err = rxChar.Write([]byte{0x62, 0x70})
	if err != nil {
		return fmt.Errorf("failed to send stop all: %w", err)
	}
	time.Sleep(300 * time.Millisecond)

	return nil
}

func (fb *FinchBluetooth) monitorNotifications() {
}

func (fb *FinchBluetooth) WriteRx(command []byte) error {
	if fb.rxChar == nil {
		return fmt.Errorf("not connected")
	}

	fmt.Printf("[FinchBLE] Writing command: %v\n", command)
	return fb.rxChar.Write(command)
}

func (fb *FinchBluetooth) Read() ([]byte, error) {
	deadline := time.Now().Add(3 * time.Second)

	for time.Now().Before(deadline) {
		select {
		case data := <-fb.dataReady:
			if len(data) >= 20 {
				return data, nil
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil, fmt.Errorf("timed out waiting for sensor data from Finch")
}

func (fb *FinchBluetooth) Close() error {
	fmt.Println("[FinchBLE] Closing connection...")

	var errs []error

	if fb.txChar != nil {
		if err := fb.txChar.Notify(false); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop notifications: %w", err))
		}
	}

	if fb.device != nil {
		if err := fb.device.Disconnect(); err != nil {
			errs = append(errs, fmt.Errorf("failed to disconnect: %w", err))
		}
	}

	close(fb.dataReady)

	if len(errs) > 0 {
		return fmt.Errorf("multiple errors during close: %v", errs)
	}

	return nil
}

func ConnectByAddress(macAddress string) (*FinchBluetooth, error) {
	builder := &FinchBluetoothBuilder{}
	return builder.ConnectByAddress(macAddress)
}
