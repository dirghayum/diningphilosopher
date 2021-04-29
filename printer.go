//DirghayuMainali(L20445249)
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var header_printed bool = false
var isRunning_printer bool = false //is the printer running or exited
var printerTimer = DP_TIMER

func main() {
	go printer_timer()
	display()
}

func display() {
	//fmt.Println("Accepting connection....")
	//Start a UDP conection
	//--------------------------------------------------
	con, err := net.ResolveUDPAddr(UDP_NETWORK, ":"+strconv.Itoa(PRINTER_PORT))
	dd(err)
	connection, err := net.ListenUDP(UDP_NETWORK, con)
	defer connection.Close()
	//-----------------------------------------------------

	//create a buffer to hold incoming message
	buffer := make([]byte, 1024)

	//Accept incoming message in an infinite loop and take
	fmt.Println("Accepting message from philosophers")

	for {
		var philo_status []string
		len, _, _ := connection.ReadFromUDP(buffer)
		message := buffer[0:len]
		//fmt.Println(message)
		json.Unmarshal(message, &philo_status)
		if !header_printed {
			print_header(philo_status)
		}
		print_philo_status(philo_status)

	}
}

func print_header(philo_status []string) {
	if !isRunning_printer {
		isRunning_printer = true
	}
	var header string = "Current time      " + get_gap(" ")
	var philo_status_temp []string
	for i := 0; i < len(philo_status); i++ {
		header += "Philo#" + strconv.Itoa(i+1) + " " + get_gap(" ")
		philo_status_temp = append(philo_status_temp, STATE_THINKING)
	}
	header += "\n" + "===========" + get_gap("=")
	for i := 0; i < len(philo_status)*2; i++ {
		header += "========"
	}
	fmt.Println(header)
	print_philo_status(philo_status_temp)
	header_printed = true
}

func print_philo_status(philo_status []string) {
	var row string = ""
	dt := time.Now()
	day := strconv.Itoa(dt.Day())
	row += dt.Month().String() + " " + day + " " + dt.Format("15:04:05")
	for i := 0; i < len(philo_status); i++ {
		row += get_gap(" ") + philo_status[i]
	}
	fmt.Println(row)
}

func get_gap(c string) string {
	var gap string = ""
	for i := 0; i < COLUMN_GAP; i++ {
		gap += c
	}
	return gap
}

func printer_timer() {
	for {
		if isRunning_printer {
			printerTimer--
			//fmt.Println("Timer : ",forkTimer)
		}
		time.Sleep(1000 * time.Millisecond)
		if printerTimer <= 0 {
			fmt.Println("Shutting down...")
			os.Exit(0)
		}
	}
}
