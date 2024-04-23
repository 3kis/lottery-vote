package actors

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"

	"lottery-vote/utils"
)

type TicketActor struct {
	CurrentTicket string
	TicketUsage   int
}

func (state *TicketActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *GenerateTicket:
		state.CurrentTicket = fmt.Sprintf("%s", utils.GenerateRandomString(TicketLength))
		state.TicketUsage = TicketUsage
	case *ValidateTicket:
		//校验传过来的ticket是否和自己的ticket一样, 不一样返回false, 一样返回true
		if state.CurrentTicket != msg.Ticket || state.TicketUsage-1 < 0 {
			fmt.Printf("err!!! TicketActor ValidateTicket: input ticket==>%s  curTicket==>%s\n", msg.Ticket, state.CurrentTicket)
			ctx.Respond(&VoteResponse{Success: false})
			return
		}
		state.TicketUsage--
		ctx.Respond(&VoteResponse{Success: true})
	case *QueryTicket:
		log.Printf("TicketActor QueryTicket start!!!: %s", state.CurrentTicket)
		ctx.Respond(&TicketResponse{Ticket: state.CurrentTicket})

	}
}
