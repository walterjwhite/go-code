package main

import (
	"fmt"
	"time"

	"github.com/walterjwhite/robotics/finch/v2_0/finch"
)

func main() {
	bt, err := finch.ConnectByAddress("CA:C7:3C:22:BD:64")
	if err != nil {
		fmt.Printf("Failed to connect to Finch: %v\n", err)
		return
	}
	defer func() {
		if err := bt.Close(); err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
	}()

	robot := finch.NewFinchRobot(bt)

	lightShow(robot)
	playNotes(robot)
	drive(robot)
	readSensors(robot)
}


func lightShow(finch *finch.FinchRobot) {
	fmt.Println("=== Light show ===")

	if err := finch.SetBeak(255, 0, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	time.Sleep(500 * time.Millisecond)

	if err := finch.SetBeak(0, 255, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	time.Sleep(500 * time.Millisecond)

	if err := finch.SetBeak(0, 0, 255); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	time.Sleep(500 * time.Millisecond)

	if err := finch.SetBeak(255, 128, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	time.Sleep(500 * time.Millisecond)

	if err := finch.SetBeak(0, 0, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}

	for port := 1; port <= 4; port++ {
		if err := finch.SetTailLight(port, 0, 200, 200); err != nil {
			fmt.Printf("Error setting tail light: %v\n", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if err := finch.SetTailLight(0, 0, 0, 0); err != nil {
		fmt.Printf("Error setting tail light: %v\n", err)
	}
}

func playNotes(finch *finch.FinchRobot) {
	fmt.Println("=== Playing notes ===")

	notes := []int{60, 64, 67, 72}
	durations := []int{300, 300, 300, 500}

	for i := 0; i < len(notes); i++ {
		err := finch.PlayNote(notes[i], durations[i])
		if err != nil {
			fmt.Printf("Error playing note: %v\n", err)
		}
		time.Sleep(time.Duration(durations[i]+50) * time.Millisecond)
	}
}

func drive(finch *finch.FinchRobot) {
	fmt.Println("=== Driving ===")

	if err := finch.SetBeak(0, 255, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	if err := finch.Forward(50); err != nil {
		fmt.Printf("Error moving forward: %v\n", err)
	}
	time.Sleep(1000 * time.Millisecond)

	if err := finch.SetBeak(255, 255, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	if err := finch.TurnRight(40); err != nil {
		fmt.Printf("Error turning right: %v\n", err)
	}
	time.Sleep(700 * time.Millisecond)

	if err := finch.SetBeak(0, 255, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	if err := finch.Forward(50); err != nil {
		fmt.Printf("Error moving forward: %v\n", err)
	}
	time.Sleep(1000 * time.Millisecond)

	if err := finch.SetBeak(255, 0, 0); err != nil {
		fmt.Printf("Error setting beak color: %v\n", err)
	}
	if err := finch.StopMotors(); err != nil {
		fmt.Printf("Error stopping motors: %v\n", err)
	}
}

func readSensors(finch *finch.FinchRobot) {
	fmt.Println("=== Reading sensors ===")

	distance, err := finch.ReadDistanceCm()
	if err != nil {
		fmt.Printf("Error reading distance: %v\n", err)
	} else {
		fmt.Printf("  Distance: %d cm\n", distance)
	}

	lights, err := finch.ReadLightSensors()
	if err != nil {
		fmt.Printf("Error reading light sensors: %v\n", err)
	} else {
		fmt.Printf("  Light sensors — left: %d  right: %d\n", lights[0], lights[1])
	}
}
