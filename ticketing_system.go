package main

import (
	"fmt"
	"time"
)

// ticketAgent represents the vehicle through which different database tables will be able to
// receive unique IDs, even in a distributed environment.
type ticketAgent struct {
	idChannel       chan int64
	nextAvailableId *int64
	tableName       string
}

// ticketingService is intended to be run as a Go routine and acts as
// a singular, central authority of generating / returning unique IDs for whatever
// ticketAgent requests an ID.
func ticketingService(agentChannel <-chan ticketAgent) {
	for {
		agent := <-agentChannel
		agent.idChannel <- *agent.nextAvailableId
		*agent.nextAvailableId++
	}
}

// requestID is intended to be run as a Go routine and simulates a process that will require a
// unique ID for the object it is creating.
func requestID(processName string, agent ticketAgent, agentChannel chan<- ticketAgent) {
	for {
		agentChannel <- agent
		id := <-agent.idChannel
		fmt.Printf("%s received ID %d for the %s table\n", processName, id, agent.tableName)
		time.Sleep(2 * time.Second)
	}
}

// initializeTicketAgent initializes a new ticketAgent struct, settings its nextAvailableId to 0.
// This should be used when setting up new database tables that will require unique ID generation.
func initializeTicketAgent(tableName string) ticketAgent {
	return ticketAgent{idChannel: make(chan int64), nextAvailableId: new(int64), tableName: tableName}
}

func main() {
	agentChannel := make(chan ticketAgent)
	tasksAgent := initializeTicketAgent("Tasks")
	notesAgent := initializeTicketAgent("Notes")

	go ticketingService(agentChannel)

	go requestID("Process 1", tasksAgent, agentChannel)
	go requestID("Process 2", tasksAgent, agentChannel)
	go requestID("Process 3", notesAgent, agentChannel)
	go requestID("Process 4", notesAgent, agentChannel)

	for {

	}
}
