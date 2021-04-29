//DirghayuMainali(L20445249)

package main
import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)


var fork_list []string    //This will hold the IP of all forks
var timer int = DP_TIMER  //90 seconds timer
var isRunning bool = false

//entry point
func main() {
	done := make(chan int)
	 go run_timer(done)  //run the timer in the background
	 go start_accepting_message(done)   //open conection to receive messages from fork and philosopher
	<-done
}

func start_accepting_message(done chan int) {
	//Start a UDP conection
	//--------------------------------------------------
	con, err := net.ResolveUDPAddr(UDP_NETWORK, ":"+strconv.Itoa(MONITOR_PORT))
	dd(err)
	connection, err := net.ListenUDP(UDP_NETWORK, con)
	dd(err)
	defer connection.Close()
	//-----------------------------------------------------

	//create a buffer to hold incoming message
	buffer := make([]byte, 1024)

	//Accept incoming message in an infinite loop and take
	//action based on who is sending the message
	for {
		len, addr, err := connection.ReadFromUDP(buffer)
		message := strings.TrimSpace(string(buffer[0:len]))
		dd(err)
		if message == "phil" {
			fmt.Println ("Received request from philosopher to get the list of forks ..")
			fork_list_serialized, _ := json.Marshal(fork_list)
			//send the list of forks to philosopher
			_, err = connection.WriteToUDP(fork_list_serialized, addr)
			dd(err)
			//once philosopher had contacted monitor, start the timer
			isRunning=true;
			fmt.Println("Timer of ",DP_TIMER," seconds started")

		} else if message == "fork" {
			_, err = connection.WriteToUDP([]byte("Fork registered"), addr)
			if !Contains(fork_list, addr.IP.String()) {
				fork_list = append(fork_list, addr.IP.String())
			}
			fmt.Println ("Fork Registered :"+addr.IP.String())
			fmt.Println(fork_list)

			if err != nil {
				fmt.Println(err)
				return
			}
		} else if message == "printer" {
			fork_list_serialized, _ := json.Marshal(fork_list)
			_, err = connection.WriteToUDP(fork_list_serialized, addr)

			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}

func run_timer(done chan int) {
	for {
		if isRunning {
			timer--
			fmt.Println("Timer : ",timer)
		}
		time.Sleep(1000 * time.Millisecond)
		if timer<=0{
			fmt.Println("Shutting down...")
			os.Exit(0)
		}
	}
}

func send_kill_signal(){
	//send kill to fork and philo and printer
}

func Contains(forks []string, fork string) bool {
	var found bool = false
	for _, f := range forks {
		if fork == f {
			found = true
			break
		}
	}
	return found
}


