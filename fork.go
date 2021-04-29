//Dirghayu Mainali(L20445249)

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var is_clean bool = true            //status of the fork, if it is dirty or clean
var isRunning_fork bool = false      //is the fork running or exited
var forkTimer = DP_TIMER
var current_philosopher string = "" //which philosopher has this fork at this time
var previos_philosopher string = "" // which philosopher had hold this fork previously
var killfork chan bool

func main() {
	//registers its IP to the fork_list
	go register_fork()
	//start accepting message from the philosopher
	go accept_message_from_philosopher()
	go fork_timer()
	<-killfork
	//todo: start another thread to kill the fork after a timer
}

func register_fork() {
	udpaddr, err := net.ResolveUDPAddr("udp4", MONITOR_IP+":"+strconv.Itoa(MONITOR_PORT))
	udpconn, err := net.DialUDP("udp4", nil, udpaddr)
	dd(err)

	fmt.Printf("The UDP server is %s\n", udpconn.RemoteAddr().String())
	//defer udpconn.Close()

	data := []byte("fork" + "\n")
	_, err = udpconn.Write(data)
	dd(err)

	buffer := make([]byte, 1024)
	n, _, err := udpconn.ReadFromUDP(buffer)
	dd(err)

	//if monitor sends kill signal, just die
	monitor_reply := string(buffer[0:n])

	fmt.Println("monitor reply = ", monitor_reply)
	if monitor_reply == "KILL" {
		isRunning_fork = false
		os.Exit(0)
	}

	fmt.Println(isRunning_fork)
}

func accept_message_from_philosopher() {
	fmt.Println("Waiting for philosopher to conect ..")
	l, err := net.Listen("tcp", ":"+strconv.Itoa(FORK_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			os.Exit(1)
		}

		buf := make([]byte, 4) //pick or drop
		_, err1 := conn.Read(buf)
		dd(err1)
		data := string(buf)

		//if philosopher contacts the fork, start the timer
		if !isRunning_fork {
			isRunning_fork = true
		}

		//pick or drop fork based on philosophers request
		if data == PICK_FORK {
			if is_clean {
				current_philosopher = conn.RemoteAddr().String()
				conn.Write([]byte(PICK_FORK_ACTION_PASS))
				is_clean = false
				fmt.Println("Fork picked by ",conn.RemoteAddr().String())
			} else {
				conn.Write([]byte(PICK_FORK_ACTION_FAIL))
			}
		} else if data == DROP_FORK {
				is_clean = true
				current_philosopher = ""
				conn.Write([]byte(DROP_FORK_ACTION_PASS))
			    fmt.Println("Fork dropped by ",conn.RemoteAddr().String())
		}
	}
}

func fork_timer() {
	for {
		if isRunning_fork {
			forkTimer--
			//fmt.Println("Timer : ",forkTimer)
		}
		time.Sleep(1000 * time.Millisecond)
		if forkTimer<=0{
			fmt.Println("Shutting down...")
			os.Exit(0)
		}
	}
}
