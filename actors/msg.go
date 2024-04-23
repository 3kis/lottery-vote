package actors

const (
	TicketActorCount = 1000
	VoteActorCount   = 1000
	TicketLength     = 10
	TicketUsage      = 100
)

type ConsistentHash struct {
	circle           map[uint32]string
	actors           []string
	virtualNodeCount int
}

/*func NewConsistentHash(actors []string, virtualNodeCount int) *ConsistentHash {
	ch := &ConsistentHash{
		circle:           make(map[uint32]string),
		actors:           actors,
		virtualNodeCount: virtualNodeCount,
	}
	for _, actor := range actors {
		for i := 0; i < virtualNodeCount; i++ {
			hash := sha1.Sum([]byte(fmt.Sprintf("%s%d", actor, i)))
			ch.circle[binary.BigEndian.Uint32(hash[:])] = actor
		}
	}
	return ch
}

func (ch *ConsistentHash) Get(key string) string {
	if len(ch.circle) == 0 {
		return ""
	}
	hash := sha1.Sum([]byte(key))
	h := binary.BigEndian.Uint32(hash[:])
	for i := h; ; i++ {
		if i == math.MaxUint32 {
			i = 0
		}
		if actor, ok := ch.circle[i]; ok {
			return actor
		}
	}
}*/

type GenerateTicket struct{}
type TicketResponse struct{ Ticket string }
type ValidateTicket struct {
	Ticket string
}
type VoteRequest struct {
	Usernames []string
}
type VoteResponse struct {
	Success bool
}
type QueryTicket struct{}

// 查询user有多少票的请求体
type VoteQuery struct {
	Username string
}

// 查询user有多少票的响应体
type VoteQueryResponse struct {
	Votes int
}
