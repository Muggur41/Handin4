package node

import (
	"sync"
)

type State int

const (
	RELEASED State = iota //constants for the state of the node
	WANTED
	HELD
)

type Node struct {
	ID    int32
	Timestamp int64
	State State
	DeferredQueue []int32
	ReplyCount int
	PeerCount int
	mu sync.Mutex
}

func NewNode(id int32, peerCount int) *Node {
	return &Node{
		ID:            id,
		Timestamp:     0,
		State:         RELEASED,
		DeferredQueue: make([]int32, 0),
		ReplyCount:    0,
		PeerCount:     peerCount,
	}
}

func (n *Node) IncrementTimestamp() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.Timestamp++
}

func (n *Node) UpdateTimestamp(otherTimestamp int64) { //lamport clock rule, take the max of the two timestamps and increment by 1
	n.mu.Lock()
	defer n.mu.Unlock()

	if otherTimestamp > n.Timestamp {
		n.Timestamp = otherTimestamp
	}
	n.Timestamp++
}


func (n *Node) ShouldDefer(requestingID int32, requestingTimestamp int64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	// 1. If we're in HELD -> defer (make them wait)
	// 2. If we're in WANTED (also requesting):
	//   - + our timestamp is lower -> we have priority -> defer
	//   - + same timestamp but our ID is lower -> we have priority -> defer
	// 3. Otherwise -> reply immediately (let them in)

	if n.State == HELD {
		return true
	}
	if n.State == WANTED {
		if requestingTimestamp > n.Timestamp {
			return true
		}
		if n.Timestamp == requestingTimestamp && n.ID < requestingID {
			return true
		}
	}
	return false
}

func (n *Node) DeferReply(nodeID int32) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.DeferredQueue = append(n.DeferredQueue, nodeID)
}
		
func (n *Node) RequestCS() int64{
	n.mu.Lock()
	defer n.mu.Unlock()
	
	n.State = WANTED //Set state to wanted
	n.Timestamp++ //incement timestamp for the request
	n.ReplyCount = 0 //reset reply count 
	return n.Timestamp //return the timestamp of the request so we can send it in our requests
}

func (n *Node) EnterCS() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = HELD //set state to held
}

func (n *Node) ReleaseCS() []int32 {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.State = RELEASED //set state back to released
	deferred := n.DeferredQueue //grab the list of deferred nodes waiting for a reply
	n.DeferredQueue = make([]int32, 0) //clear the deferred queue
	return deferred //return the list so we can send a reply to them
}

func (n *Node) ReceiveReply() bool{
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ReplyCount++ //increment the reply count when we receive a reply

	return n.ReplyCount >= n.PeerCount //return true if we've received replies from all peers, false otherwise
}