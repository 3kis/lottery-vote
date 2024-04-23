package actors

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

var (
	System          *actor.ActorSystem
	TicketActorPIDs []*actor.PID
	VoteActorPIDs   []*actor.PID
	VoteWriterPID   *actor.PID
	//consistentHash  *ConsistentHash
)

func InitActorSystem() {
	rand.Seed(time.Now().UnixNano())

	System = actor.NewActorSystem()

	ticketProps := actor.PropsFromProducer(func() actor.Actor {
		return &TicketActor{}
	})
	for i := 0; i < TicketActorCount; i++ {
		pid := System.Root.Spawn(ticketProps)
		TicketActorPIDs = append(TicketActorPIDs, pid)
	}

	voteProps := actor.PropsFromProducer(func() actor.Actor {
		return &VoteActor{userVotes: make(map[string]int)}
	})
	for i := 0; i < VoteActorCount; i++ {
		pid := System.Root.Spawn(voteProps)
		VoteActorPIDs = append(VoteActorPIDs, pid)
	}

	writerProps := actor.PropsFromProducer(func() actor.Actor {
		return &VoteWriter{}
	})
	VoteWriterPID = System.Root.Spawn(writerProps)

	actorNames := make([]string, TicketActorCount)
	for i := range actorNames {
		actorNames[i] = fmt.Sprintf("TicketActor-%d", i)
	}
	//consistentHash = NewConsistentHash(actorNames, 10)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			for _, pid := range TicketActorPIDs {
				System.Root.Send(pid, &GenerateTicket{})
			}
		}
	}()
}
