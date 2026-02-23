package tracetree

import (
	"fmt"
)

type Node struct {
	PID            string
	ID             string
	PStepID        int
	StepID         int
	Parent         *Node
	Child          *Node
	Next           *Node
	Tail           *Node
	Data           interface{}
	Virtual        bool
	Count          int
	StartTimestamp uint64
	EndTimestamp   uint64
}

func NewNode(pid, id string, data interface{}) *Node {
	p := new(Node)
	p.PID = pid
	p.ID = id
	p.PStepID = -1
	p.Data = data
	p.Child = nil
	p.Next = nil
	p.Tail = nil
	p.Count = 0

	return p
}

func (this *Node) String() string {
	return fmt.Sprintf("Node pid: %s id: %s child: %v next: %v", this.PID, this.ID, (this.Child != nil), (this.Next != nil))
}
