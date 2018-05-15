package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println("Some error occured", err)
		os.Exit(-1)
	}

	defer rpio.Close()

	var (
		j         uint8
		f         float32
		dht11_dat [5]int
	)

	pin := rpio.Pin(7)
	pin.Output()
	pin.Low()
	time.Sleep(18)
	pin.PullUp()
	time.Sleep(40)

	pin.Input()

	laststate := rpio.High
	for i := 0; i < 85; i++ {
		counter := 0
		for pin.Read() == laststate {
			counter++
			time.Sleep(1)
			if counter == 255 {
				break
			}
		}

		laststate = pin.Read()

		if counter == 255 {
			break
		}

		if (i >= 4) && (i%2 == 0) {
			dht11_dat[j/8] <<= 1
			if counter > 16 {
				dht11_dat[j/8] |= 1
			}
			j++
		}
	}

	if (j >= 40) &&
		(dht11_dat[4] == ((dht11_dat[0] + dht11_dat[1] + dht11_dat[2] + dht11_dat[3]) & 0xFF)) {
		f = float32(dht11_dat[2]*9./5. + 32)
		fmt.Printf("Humidity = %d.%d %% Temperature = %d.%d C (%.1f F)\n",
			dht11_dat[0], dht11_dat[1], dht11_dat[2], dht11_dat[3], f)
	} else {
		fmt.Println("Data not good, skip\n")
	}

}
