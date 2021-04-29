//DirghayuMainali(L20445249)

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
        "net/http"
	"os"
	"strconv"
	"time"
)

type philosopher struct {
	leftFork   string //ip of left fork
	rightFork  string //ip of right fork
	status     string //thinking , waiting or eating
	index      int    //index of the philosopher
}

var fork_list_local []string
var philosopherList []philosopher

func main() {
        http.DefaultClient.Timeout = time.Minute * 10
	rand.Seed(time.Now().UTC().UnixNano()) //for printing time
	get_fork_list()
	create_philosophers()
}

func create_philosophers() {
	for i := 0; i < len(fork_list_local); i++ {
		p := philosopher{}
		p.leftFork = get_left_fork(i)
		p.rightFork = get_right_fork(i)
		p.status = STATE_THINKING //Initially all philosopher are in thinking state
		p.index = i

		go run_philosopher(p)

	}
	//Run for 90 seconds and exit
	time.Sleep(time.Duration(DP_TIMER*1000) * time.Millisecond)
	fmt.Println("Shutting down ..")
	os.Exit(0)
}

func run_philosopher(phil philosopher) {
	//fmt.Println("This is philosopher : "+string(phil.index))
	time.Sleep(time.Duration(randInt(2, 6))*1000 * time.Millisecond)
	phil.status = STATE_WAITING
	print_status(phil)

	for {
		isPicked_left := pick_fork(phil.leftFork)
		if isPicked_left {
			fmt.Println("Philosopher ",phil.index," picked the left fork")
			isPicked_right := pick_fork(phil.rightFork)
			if isPicked_right {
				fmt.Println("Philosopher ",phil.index, " now has both forks and is eating")
				phil.status = STATE_EATING
				print_status(phil)
				time.Sleep(time.Duration(randInt(2, 6)) * 1000 * time.Millisecond)
				drop_fork(phil.leftFork)
				drop_fork(phil.rightFork)
				fmt.Println("Philosopher ",phil.index," dropped both forks")
				phil.status = STATE_THINKING
				print_status(phil)
				time.Sleep(2000 * time.Millisecond)
				phil.status = STATE_WAITING
				print_status(phil)
			} else {
				drop_fork(phil.leftFork)
				fmt.Println("Philosopher ",phil.index," dropped the left fork because right fork was in use")
				//wait for half second and try again
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

func get_conection(forkIp string) net.Conn{
	//send message PICK_FORK,
	conn, err := net.Dial("tcp", forkIp+":"+strconv.Itoa(FORK_PORT))
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	return conn
}

func pick_fork(forkIp string) bool {
	fork_picked := true

	//send message PICK_FORK,
	conn, err := net.Dial("tcp", forkIp+":"+strconv.Itoa(FORK_PORT))
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(PICK_FORK))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	// Read message from fork
	var readBuffer = make([]byte, 4)
	_, err = conn.Read(readBuffer)
	var forkReply = string(readBuffer)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	//time.Sleep(1000 * time.Millisecond)
	if forkReply == PICK_FORK_ACTION_PASS {
		fork_picked = true
	}else{
		fork_picked = false
	}
        time.Sleep(200 * time.Millisecond)        
        conn.Close()
	return fork_picked
}

func drop_fork(forkIp string) bool {
	fork_dropped := true
	//send message  DROP_FORK,
	conn, err := net.Dial("tcp", forkIp+":"+strconv.Itoa(FORK_PORT))
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(DROP_FORK))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	// Read message from fork
	var readBuffer = make([]byte, 9)
	_, err = conn.Read(readBuffer)
	var forkReply = string(readBuffer)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
	if forkReply == DROP_FORK_ACTION_PASS {
		fork_dropped = true
	}else{
		fork_dropped=false
	}
        time.Sleep(100 * time.Millisecond)        
        conn.Close()
	return fork_dropped

}

func get_left_fork(philo_index int) string {
	forkIndex := philo_index % len(fork_list_local)
	if forkIndex < len(fork_list_local) {
		return fork_list_local[forkIndex]
	} else {
		return fork_list_local[len(fork_list_local)-1]
	}
}

func get_right_fork(philo_index int) string {
	forkIndex := (philo_index + 1) % len(fork_list_local)
	if forkIndex < len(fork_list_local) {
		return fork_list_local[forkIndex]
	} else {
		return fork_list_local[len(fork_list_local)-1]
	}
}

func print_status(phil philosopher) {
	//fmt.Println("Sending status to printer")
		udpaddr, err := net.ResolveUDPAddr(UDP_NETWORK, PRINTER_IP+":"+strconv.Itoa(PRINTER_PORT))
		udpconn, err := net.DialUDP(UDP_NETWORK, nil, udpaddr)
		dd(err)

		defer udpconn.Close()
	//fmt.Println(philosopherList[0].status)
	//for{
		var philostatus []string
		for i :=0;i<len(fork_list_local);i++{
			if i == phil.index {
				philostatus = append(philostatus, phil.status)
			}else{
				philostatus = append(philostatus, "--------")
			}
		}

		philosopher_serialized, _ := json.Marshal(philostatus)
		//fmt.Println(philosopher_serialized)
		_, err = udpconn.Write(philosopher_serialized)
		//time.Sleep(1000*time.Millisecond)
//	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func get_fork_list() {
	udpaddr, err := net.ResolveUDPAddr(UDP_NETWORK, MONITOR_IP+":"+strconv.Itoa(MONITOR_PORT))
	udpconn, err := net.DialUDP(UDP_NETWORK, nil, udpaddr)
	dd(err)

	//fmt.Printf("The UDP server is %s\n", udpconn.RemoteAddr().String())
	defer udpconn.Close()

	data := []byte("phil")
	_, err = udpconn.Write(data)
	dd(err)

	buffer := make([]byte, 1024)
	n, _, err := udpconn.ReadFromUDP(buffer)
	dd(err)

	//get json array of fork list and decode it into an array
	monitor_reply := buffer[0:n]
	json.Unmarshal(monitor_reply, &fork_list_local)
	//fmt.Println(monitor_reply)
	//fmt.Println(fork_list_local[0])
	//fmt.Println(fork_list_local[1])
}
