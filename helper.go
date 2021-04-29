//DirghayuMainali(L20445249)

package main

import "fmt"

var MONITOR_IP string = "140.158.130.67"
var MONITOR_PORT int = 9000

var PRINTER_IP string = "140.158.130.53"
var PRINTER_PORT int = 9001

var FORK_PORT int = 9002

var UDP_NETWORK string = "udp4"

var PICK_FORK string = "PICK"
var DROP_FORK string = "DROP"

var PICK_FORK_ACTION_FAIL string ="FAIL"
var PICK_FORK_ACTION_PASS string ="PASS"

var DROP_FORK_ACTION_FAIL string ="DROP_FAIL"
var DROP_FORK_ACTION_PASS string ="DROP_PASS"

//Philosopher states
var STATE_THINKING string= "Thinking"
var STATE_WAITING string = "Waiting "
var STATE_EATING string  = "Eating  "

//Timer or the lifespan of the DP program
var DP_TIMER int = 90

//Initial thinking time of the philosopher
var THINKING_TIME int = 1

var COLUMN_GAP int= 10

//----------------helper methods --------------------
func dd(err error){
	if err != nil {
		fmt.Println(err)
		return
	}
}
