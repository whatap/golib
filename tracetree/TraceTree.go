package tracetree

import (
	"strings"

	"github.com/whatap/golib/util/hmap"
)

const (
	DEFAULT_SPAN_SIZE      = 2048
	DEFAULT_CHILD_NODE_MAX = 1024
)

type Orphans struct {
	List []*Node
}

func NewOrphans() *Orphans {
	p := new(Orphans)
	p.List = make([]*Node, 0)
	return p
}

type TraceTree struct {
	traceID      string
	entry        *hmap.StringKeyLinkedMap
	orphansEntry *hmap.StringKeyLinkedMap
}

func NewTraceTree(traceID string) *TraceTree {
	p := new(TraceTree)
	p.traceID = traceID
	p.entry = hmap.NewStringKeyLinkedMap().SetMax(DEFAULT_SPAN_SIZE)
	p.orphansEntry = hmap.NewStringKeyLinkedMap().SetMax(DEFAULT_SPAN_SIZE)
	return p
}

func (this *TraceTree) Size() int {
	return this.entry.Size()
}

func (this *TraceTree) PutRoot(newNode *Node) {
	newNode.ID = strings.TrimSpace(newNode.ID)
	newNode.PID = strings.TrimSpace(newNode.PID)
	this.entry.Put(newNode.ID, newNode)
	this.linkChildren(newNode)
	this.linkOperhans(newNode)
}

func (this *TraceTree) Put(newNode *Node) {
	newNode.ID = strings.TrimSpace(newNode.ID)
	newNode.PID = strings.TrimSpace(newNode.PID)

	this.entry.Put(newNode.ID, newNode)
	this.linkParent(newNode)
	this.linkChildren(newNode)
}

func (this *TraceTree) linkParent(newNode *Node) {
	if newNode.PID != "" {
		if tmp := this.entry.Get(newNode.PID); tmp != nil {
			if node, ok := tmp.(*Node); ok {
				newNode.Parent = node
				if node.Child == nil {
					node.Child = newNode
					node.Tail = newNode
				} else {
					node.Tail.Next = newNode
					node.Tail = newNode
				}
				node.Count++
				return
			}
		}
	}
	this.orphansEntry.Put(newNode.ID, newNode)
}

func (this *TraceTree) linkChildren(newNode *Node) {
	removeID := make([]string, 0)
	en := this.orphansEntry.Values()
	for en.HasMoreElements() {
		tmp := en.NextElement()
		if node, ok := tmp.(*Node); ok {
			if node.PID == newNode.ID {
				node.Parent = newNode
				if newNode.Child == nil {
					newNode.Child = node
					newNode.Tail = node
				} else {
					newNode.Tail.Next = node
					newNode.Tail = node
				}
				newNode.Count++
				removeID = append(removeID, node.ID)
			}
		}
	}

	for _, k := range removeID {
		this.orphansEntry.Remove(k)
	}
}

func (this *TraceTree) linkOperhans(newNode *Node) {
	removeID := make([]string, 0)
	en := this.orphansEntry.Values()
	for en.HasMoreElements() {
		tmp := en.NextElement()
		if node, ok := tmp.(*Node); ok {
			if newNode.StartTimestamp <= node.StartTimestamp && node.StartTimestamp <= newNode.EndTimestamp {
				node.Parent = newNode
				if newNode.Child == nil {
					newNode.Child = node
					newNode.Tail = node
				} else {
					newNode.Tail.Next = node
					newNode.Tail = node
				}
				newNode.Count++
				removeID = append(removeID, node.ID)
			}
		}
	}

	for _, k := range removeID {
		this.orphansEntry.Remove(k)
	}
}

func (this *TraceTree) GetAndClear(pid string) []*Node {
	if strings.TrimSpace(pid) == "" {
		return nil
	}
	var top *Node = nil
	if tmp := this.entry.Remove(pid); tmp != nil {
		if val, ok := tmp.(*Node); ok {
			top = val
		}
	}
	if top == nil {
		return nil
	}

	lst := make([]*Node, 0)

	idx := 0
	count := 0
	top.StepID = -1

	if top.Child != nil {
		count, idx, lst = this.getAndClear(lst, top.Child, top.StepID, idx, count)
	}
	return lst
}

func (this *TraceTree) getAndClear(lst []*Node, node *Node, pStepID, idx, count int) (int, int, []*Node) {
	for node != nil {
		node.PStepID = pStepID

		if count <= DEFAULT_CHILD_NODE_MAX {
			lst = append(lst, node)
			node.StepID = idx
			idx += 10
		}
		count++
		this.entry.Remove(node.ID)

		if node.Child != nil {
			count, idx, lst = this.getAndClear(lst, node.Child, node.StepID, idx, count)
		}

		node = node.Next
	}
	return count, idx, lst
}

func (this *TraceTree) Get(pid string) []*Node {
	if strings.TrimSpace(pid) == "" {
		return nil
	}
	var top *Node = nil
	if tmp := this.entry.Get(pid); tmp != nil {
		if val, ok := tmp.(*Node); ok {
			top = val
		}
	}
	if top == nil {
		return nil
	}
	lst := make([]*Node, 0)
	lst = append(lst, top)
	count := 0

	if top.Child != nil {
		count, lst = this.get(lst, top.Child, count)
	}
	return lst
}

func (this *TraceTree) get(lst []*Node, node *Node, count int) (int, []*Node) {
	for node != nil {
		if count > DEFAULT_CHILD_NODE_MAX {
			return count, lst
		}
		lst = append(lst, node)
		count++

		if node.Child != nil {
			count, lst = this.get(lst, node.Child, count)
		}
		node = node.Next
	}
	return count, lst
}

func (this *TraceTree) SetMax(sz int) {
	this.entry.SetMax(sz)
}

func (this *TraceTree) Clear() {
	this.entry.Clear()
	this.orphansEntry.Clear()
}
