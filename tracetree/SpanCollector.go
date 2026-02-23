package tracetree

import (
	"strings"

	"github.com/whatap/golib/util/hmap"
)

const (
	DEFAULT_TRACE_SIZE = 2048
	DEFAULT_SIZE       = 100000
)

type SpanCollector struct {
	entry *hmap.StringKeyLinkedMap
}

func NewSpanCollector() *SpanCollector {
	p := new(SpanCollector)
	p.entry = hmap.NewStringKeyLinkedMap().SetMax(DEFAULT_TRACE_SIZE)
	return p
}

func (this *SpanCollector) Size() int {
	return this.entry.Size()
}

func (this *SpanCollector) PutRoot(traceID string, newNode *Node) {
	traceTree := this.getTraceTree(traceID)
	if traceTree == nil {
		traceTree = NewTraceTree(traceID)
		this.entry.Put(traceID, traceTree)
	}
	traceTree.PutRoot(newNode)

	if traceTree.Size() == 0 {
		traceTree.Clear()
		this.entry.Remove(traceID)
	}
}

func (this *SpanCollector) Put(traceID string, newNode *Node) {
	traceTree := this.getTraceTree(traceID)
	if traceTree == nil {
		traceTree = NewTraceTree(traceID)
		this.entry.Put(traceID, traceTree)
	}
	traceTree.Put(newNode)
}

func (this *SpanCollector) getTraceTree(traceID string) *TraceTree {
	var traceTree *TraceTree = nil
	if tmp := this.entry.Get(traceID); tmp != nil {
		if val, ok := tmp.(*TraceTree); ok {
			traceTree = val
		}
	}
	return traceTree
}

func (this *SpanCollector) last(me *Node) *Node {
	for me.Next != nil {
		me = me.Next
	}
	return me
}

func (this *SpanCollector) GetAndClear(traceID string, pID string) []*Node {
	if strings.TrimSpace(traceID) == "" || strings.TrimSpace(pID) == "" {
		return nil
	}
	var traceTree *TraceTree = nil
	if tmp := this.entry.Get(traceID); tmp != nil {
		if val, ok := tmp.(*TraceTree); ok {
			traceTree = val
		}
	}
	if traceTree == nil {
		return nil
	}
	result := traceTree.GetAndClear(pID)
	if traceTree.Size() == 0 {
		traceTree.Clear()
		this.entry.Remove(traceID)
	}
	return result
}

func (this *SpanCollector) Get(traceID string, pID string) []*Node {
	if strings.TrimSpace(traceID) == "" || strings.TrimSpace(pID) == "" {
		return nil
	}
	var traceTree *TraceTree = nil
	if tmp := this.entry.Get(traceID); tmp != nil {
		if val, ok := tmp.(*TraceTree); ok {
			traceTree = val
		}
	}
	if traceTree == nil {
		return nil
	}
	result := traceTree.Get(pID)

	if traceTree.Size() == 0 {
		traceTree.Clear()
		this.entry.Remove(traceID)
	}
	return result
}

func (this *SpanCollector) SetMax(sz int) {
	this.entry.SetMax(sz)
}
