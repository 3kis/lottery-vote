package actors

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type VoteActor struct {
	userVotes map[string]int
}

func (state *VoteActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *VoteRequest:
		for _, username := range msg.Usernames {
			state.userVotes[username] += 1
			log.Printf("Received vote request for user %s", username)
		}

		ctx.Send(VoteWriterPID, msg)
		ctx.Respond(&VoteResponse{Success: true})
	//查找自己的投过的user现在有了多少票了
	case *VoteQuery:
		votes, ok := state.userVotes[msg.Username]
		if !ok {
			votes = 0
		}
		ctx.Respond(&VoteQueryResponse{Votes: votes})

	}
}
