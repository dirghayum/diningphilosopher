# Programming Project – 2
### Dirghayu Mainali (L20445249)

## PROJECT DESCRIPTION

This project is based on the implementation of the Dining Philosophers problem without using shared memory or a shared clock. The Dining Philosopher Problem states that K philosophers seated around a circular table with one fork between each pair of philosophers. There is one fork between each philosopher. A philosopher may eat if he can pick up the two forks adjacent to him. One fork may be picked up by any one of its adjacent followers but not both.  [source: internet] 

When we are solving the Dining Philosopher Problem in a distributed system, things get little complicated because we must take into consideration the network and the issues with network communication. Since there is no shared memory, message passing is the only way to communicate between forks, philosophers and server and printer.

## REQUIREMENTS
- Go Programming
- Any Unix, CentOS Windows system

## HOW TO RUN THE PROGRAM
 
Before executing the program, decide in what machine you are going to run server and printer and then update the IP address of those two machine in helper.go file
   	var MONITOR_IP string = "140.158.130.67”
              var PRINTER_IP string = "140.158.130.53"

To run this program in LAB, I have used the following IPs:

OS:Ubuntu 18.0.4
delta19.bmt.lamar.edu, 140.158.130.67 - server
delta21.bmt.lamar.edu, 140.158.130.53 - printer
delta22.bmt.lamar.edu, 140.158.130.41 - philosopher

OS:CentOS 6.0
sigma24.bmt.lamar.edu, 140.158.128.21 - fork1 
sigma25.bmt.lamar.edu, 140.158.128.15 - fork2 
sigma26.bmt.lamar.edu, 140.158.130.213 -fork3 
sigma23.bmt.lamar.edu, 140.158.128.12 - fork4 
sigma27.bmt.lamar.edu, 140.158.131.163 -fork5 
STEPS:

## STEPS:

1. Copy the Dining philosopher source code and paste it in desktop

You need to copy paste the source code in all different computers OR paste the source code in one computer and mount that folder in all other computers. 
Since I used git bash to copy my source code from my laptop to the computer with the given IPs, I used the following command.

----
scp *.go dmainali@140.158.130.67:/home/bmt.lamar.edu/dmainali/dp
----


This way all the files will be copied in the given IP computer.
(Note: Before copying the file, please ping all the computer with the one you have marked as server computer to make sure those computers can communicate)

.Server

CD into the folder where you have copied the source code in remote computer
----    
$ cd dp
----

Build server.go
----
$ go build server.go helper.go
----

Run server
----
$ ./server
----

 
-------------------------------------------------------------

.We need to run 5 forks. 

Copy paste the source code in 5 different machines and repeat this process. 
Cd inside the source code folder (dp)
$ cd dp

Build fork.go
$ go build fork.go helper.go

Run fork
$ ./fork

NOTE# Since the fork program is using a hardcoded port, you cannot run two forks in a single machine. If you want to run more than two forks in a machine, you need to change the port number of fork and build and run so it might be easier to run it in 5 different machines)

 

-------------------------------------------------------------
.Printer
Copy paste the source code 

Cd inside the source code folder (dp)
$ cd dp

Build printer.go
$ go build printer.go helper.go

Run printer
$ ./printer


 
-------------------------------------------------------------

.
philosopher
Copy paste the source code. 

Cd inside the source code folder (dp)
$ cd dp

Build philosopher.go
$ go build philosopher.go helper.go

Run philosopher
$ ./philosopher

When you run a philosopher, it will create multiple philosophers in separate independent processes. The program will create the same number of philosophers as the number of forks. After the philosopher is started, you can see the status of all philosophers in the printer window.

 

## CONNECTION IN THE SYSTEM
FORK →(UDP Connection) →  SERVER
PHILOSOPHER →(UDP Connection) →  SERVER
PHILOSOPHER →(TCP Connection) →  FORK
PHILOSOPHER →(UDP Connection) →  PRINTER

### MODULES
There are 5 go files  in our program:
1.	Helper
2.	Forks
3.	Philosopher
4.	Server
5.	Printer	
Helper:

This file contains all the main variables and constants that are used in the program, values like display format, server ip and port, timer etc.

##### Server:

Server is the first program that we need to run. It is responsible to 
●	Save the IP addresses of the Forks in a list.
●	Provide philosopher with the list of forks

Once all forks are registered and philosopher gets the list of forks, this server just starts the countdown timer for 90 seconds. Fork and philosopher will start to talk with each other and do not communicate with the server afterwards.

##### Fork:

Each fork has to register itself with the server. To talk with the server, it will need the IP and port of the server. We can run multiple forks on different machine and all the forks will register itself in the server. After registering, it will wait for the message from philosopher. Each fork has two states - clean or dirty. 
●	If any philosopher picks the fork, it will set the status as dirty and 
●	if the philosopher drops the fork, it will set the status as clean. 

Philosopher and fork talk thru TCP socket. Whenever a fork gets the first message from a philosopher, it will start to run the timer. After 90 seconds, the fork will exit.

##### Philosopher:
	
When we run the philosopher, it will first communicate with the server to get a list of forks.
Philosopher will then establish TCP connection with the nearest two forks. It will also establish a UDP connection with the printer.

At first each philosopher is in THINKING state. After waiting a random amount of time (1-3 seconds) it will enter WAITING state. It will then send a request to PICK the left fork. If the left fork is in a clean state, the philosopher will pick that fork and the fork will mark itself as dirty. Then it will send a request to pick the right fork.
●	If the RIGHT fork is clean, it will pick that fork and the fork will mark itself as dirty. This means philosopher now has both forks and it enters EATING state
●	If RIGHT fork is already dirty, it will drop the LEFT fork, and the LEFT fork will mark itself as clean again. The philosopher will wait for some time and same process of picking forks repeats
If philosopher is in EATING state, it will eat for a random time (2-6 seconds) and then after eating it will drop both forks. It then enters THINKING state again. And the cycle repeats.
After each status change, it will send a message to the printer to log the time and its status.

The philosopher will run for 90 seconds and then exit.

##### Printer:	

As the name suggests, the printer is responsible to print the status of philosophers. After each status change the philosophers send their status to the printer. This communication takes place thru UDP socket. Printer will also run for 90 seconds and when all philosophers exit, it will terminate its own process.

## PROGRAM FLOW
At first, we run Server on one machine. It will serve as a UDP server and start accepting connections. It will serve to forks and philosophers. The server will keep running in background throughout the program life cycle

1.	After running Server, we need to run forks in N different machines. When fork.go is run, it will create a socket to communicate with Server in same port and IP.  The fork will then send its ip and port to the server, which on receiving by server is stored in a list. You can run as many client forks as you like in different machines, the server will keep track of all. On running the fork module, it will create child processes for each fork. As the forks are added to server list, it will run a timer in a separate thread. The function of the timer is to send a kill signal to the forks after "timer" duration of time. The forks will terminate all their background processes on receiving the kill signal sent by Server. It's like an hourglass. Server also serves philosophers on providing the forks details. When philosopher process first runs, it will send a request to server to get the list of registered forks. After getting the fork detail, it will establish TCP connection with fork and then start   sending and receiving messages to and from the nearest two forks.

2.	All forks are initially clean. So, at first each fork can be grabbed by any philosopher. After the forks are created, it just waits for the philosopher to send message. 

3.	After philosopher are created, each philosopher will find its nearest two forks. This is done with the help of fork id. The fork id is the index of the fork that was maintained by the Server. Each philosopher will have access to its nearest two forks. Philosopher 'n' will have access to fork 'n-1' and n. For example, philosopher 5 will have access to fork 4 and fork 5. 
At first all philosophers will pick the n-1 th fork. After that it will try to access nth fork. 
If the second fork is not available, it will put down the first fork, wait for a random time and try again. This works because the wait time of different philosopher are not same. So, if a philosopher puts down fork and waits 3 second, other philosopher might just wait two second and grab the clean fork before the first philosopher grabs it. The first and second forks are swapped in when the philosopher cannot access any of the fork, so that in next iteration it will try to pick up the second fork first.
Throughout the program, the philosopher can be in any of the following three states:
- When the philosopher has no fork in its any hand, it is in *THINKING* state. 
- When the philosopher has 1 fork in its hand and waiting for another fork, it is in *WAITING* state.
- When the philosopher has both fork in its hands, it is in *EATING* state.
When philosophe gets both fork1 and fork2, it will start eating for a random duration of time and then clean the fork and put it down. Philosopher will know if the fork is dirty or clean by sending it a message and checking the response.

4.	When philosopher requests pick up or drop a fork, it will send a message to fork. The fork will check the message to see it is for the pick request or drop request.

-If it is for pick request, it reads the status if it is currently dirty or clean. If dirty , it replies with a message which contains the flag 'pickup_fail' and the id of philosopher which is currently using it.

-If it is a drop request, it change the status from dirty to clean and sets fork_user to None. and sends the response to philosopher which contains the message 'drop_success'

5.	After program runs for 90 seconds, all background processes are terminated. Server will send kill signal to all forks simultaneously to let them know that time is over. The philosophers and Printer have their own timers so they all will terminate.

 
This application has been tested on all major OS, windows, mac and linux.
