package pathutil

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
)

const (
	PATHTREE_TYPE_PATHS   = 1
	PATHTREE_TYPE_VALUES  = 2
	PATHTREE_TYPE_ENTRIES = 3
)

type PathTree struct {
	Top   *ENTRY
	count int
	lock  sync.Mutex
}

func NewPathTree() *PathTree {
	p := new(PathTree)

	p.Top = NewENTRY()

	return p
}

func (this *PathTree) Insert(path string, value interface{}) interface{} {

	if path == "" {
		return nil
	}

	return this.InsertArray(strings.Split(path, "/"), value)
	//return insert(StringUtil.split(path, '/'), value);
}

func (this *PathTree) InsertArray(paths []string, value interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if paths == nil || len(paths) == 0 || value == nil {
		return nil
	}
	if len(paths) == 1 {
		return nil
	}
	path := NewPATH(paths)

	if this.Top.child == nil {
		cur := this.Top
		cur.child = NewENTRY()
		cur.child.parent = cur
		cur = cur.child
		cur.node = path.Node()
		this.count++

		//fmt.Println("PathTree.InsertArray expand count=" , this.count)
		return this.expand(cur, path, value)
	} else {
		//fmt.Println("PathTree.InsertArray insert")
		return this.insert(this.Top, this.Top.child, path, value)
	}
}

func (this *PathTree) expand(cur *ENTRY, path *PATH, value interface{}) interface{} {
	for path.HasChild() {
		path.Level++
		cur.child = NewENTRY()
		cur.child.parent = cur
		cur = cur.child
		cur.node = path.Node()
		this.count++
	}

	old := cur.value
	cur.value = value
	return old
}

func (this *PathTree) insert(p *ENTRY, cur *ENTRY, path *PATH, value interface{}) interface{} {

	if path.Node() == cur.node {
		if path.HasChild() == false {
			old := cur.value
			cur.value = value
			return old
		}
		path.Level++
		if cur.child != nil {
			return this.insert(cur, cur.child, path, value)
		} else {
			cur.child = NewENTRY()
			cur.child.parent = cur
			cur = cur.child
			cur.node = path.Node()
			return this.expand(cur, path, value)
		}

	} else if cur.right != nil {
		return this.insert(p, cur.right, path, value)
	} else {

		if path.Node() == "*" {
			// *노드는 오른쪽끝에 추가한다.
			cur.right = NewENTRY()
			cur.right.parent = p
			cur = cur.right
			cur.node = path.Node()
			this.count++
			return this.expand(cur, path, value)
		} else {
			// 일반 노드는 부모 노드의 첫번째 자식으로 등록한다.
			cur = NewENTRY()
			cur.parent = p
			cur.right = p.child
			p.child = cur
			cur.node = path.Node()
			this.count++
			return this.expand(cur, path, value)
		}
	}
}

func (this *PathTree) Find(path string) interface{} {
	var ret interface{}
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("WA825", "PathTree.Find Recover", r) //, string(debug.Stack()))
			ret = nil
		}
	}()

	if path == "" {
		return nil
	}
	ret = this.FindArray(strings.Split(path, "/"))

	return ret
}

func (this *PathTree) FindArray(path []string) interface{} {
	var ret interface{}
	defer func() {
		if r := recover(); r != nil {
			// logutil.Println("WA826", "PathTree.Find Recover", r) //, string(debug.Stack()))
			ret = nil
		}
	}()

	if path == nil || len(path) == 0 {
		return nil
	}
	ret = this.find(this.Top.child, NewPATH(path))
	return ret
}

func (this *PathTree) find(cur *ENTRY, m *PATH) interface{} {
	//logutil.Println("find =", m.Node())
	// Node 비어있을 경우 예외처리
	if cur == nil {
		return nil
	}

	if cur.Include(m.Node()) {
		if m.HasChild() == false {
			return cur.value
		}
		m.Level++

		if cur.child != nil {
			return this.find(cur.child, m)
		}
	} else if cur.right != nil {
		return this.find(cur.right, m)
	}
	return nil
}

func (this *PathTree) Size() int {
	return this.count
}

func (this *PathTree) Paths() PathTreeEnumeration {
	return NewPathTreeEnumer(PATHTREE_TYPE_PATHS)
}

func (this *PathTree) Values() PathTreeEnumeration {
	return NewPathTreeEnumer(PATHTREE_TYPE_VALUES)
}

func (this *PathTree) Entries() PathTreeEnumeration {
	return NewPathTreeEnumer(PATHTREE_TYPE_ENTRIES)
}

type PathTreeEnumeration interface {
	HasMoreElements() bool
	NextElement(top *ENTRY) interface{}
}

type PathTreeEnumer struct {
	entry *ENTRY
	Type  int
}

func NewPathTreeEnumer(Type int) *PathTreeEnumer {
	p := new(PathTreeEnumer)
	p.Type = Type
	return p
}

func (this *PathTreeEnumer) HasMoreElements() bool {
	return this.entry != nil
}

func (this *PathTreeEnumer) NextElement(top *ENTRY) interface{} {
	if this.entry == nil {
		//throw new NoSuchElementException("no more next");
		return nil
	}

	e := this.entry
	if this.entry.child != nil {
		this.entry = this.entry.child
	} else {
		for this.entry != nil && this.entry.right == nil {
			this.entry = this.entry.parent
			if this.entry == top {
				this.entry = nil
				switch this.Type {
				case PATHTREE_TYPE_PATHS:
					return e.Path(top)
				case PATHTREE_TYPE_VALUES:
					return e.Value
				default:
					return e
				}
			}
		}
		if this.entry != nil {
			this.entry = this.entry.right
		}
	}
	switch this.Type {
	case PATHTREE_TYPE_PATHS:
		return e.Path(top)
	case PATHTREE_TYPE_VALUES:
		return e.value
	default:
		return e
	}
}

type PATH struct {
	Nodes []string
	Level int
}

func NewPATH(nodes []string) *PATH {
	p := new(PATH)
	p.Nodes = nodes

	return p
}

func (this *PATH) HasChild() bool {
	return (this.Level+1 < len(this.Nodes))
}

func (this *PATH) Node() string {
	return this.Nodes[this.Level]
}

type ENTRY struct {
	node   string
	value  interface{}
	right  *ENTRY
	child  *ENTRY
	parent *ENTRY
}

func NewENTRY() *ENTRY {
	p := new(ENTRY)
	return p
}

func (this *ENTRY) Include(v string) bool {
	return (this.node == "*" && v != "") || this.node == v
}

func (this *ENTRY) Value() interface{} {
	return this.value
}

func (this *ENTRY) Node() string {
	return this.node
}

func (this *ENTRY) Path(top *ENTRY) []string {
	cur := this

	// TODO
	sk := list.New()
	for cur != top {
		sk.PushFront(cur.node)
		//sk.add(cur.node);
		cur = cur.parent
	}

	arr := make([]string, sk.Len())
	i := 0
	for e := sk.Back(); e != nil; e = e.Next() {
		arr[i] = e.Value.(string)
		i++
	}

	return arr
}

//	public static void main(String[] args) {
//		PathTree<String> t = new PathTree<String>();
//		t.insert("/cube/pcode/*/history/series", "/cube/pcode/{pcode}/history/series");
//		t.insert("/tx_country/pcode/*/cube/point", "/tx_country/pcode/{pcode}/cube/point");
//		t.insert("/rt_user/pcode/*/cube/series", "/rt_user/pcode/{pcode}/cube/series");
//		t.insert("/summary/pcode/*/cube", "/summary/pcode/{pcode}/cube");
//		t.insert("/hitmap/pcode/*/cube/series", "/hitmap/pcode/{pcode}/cube/series");
//		t.insert("/tx/pcode/*/top/cube", "/tx/pcode/{pcode}/top/cube");
//		t.insert("/event/pcode/*/top/cube", "/event/pcode/{pcode}/top/cube");
//		t.insert("/tps_res_time/pcode/*/cube/series", "/tps_res_time/pcode/{pcode}/cube/series");
//		t.insert("/cpu_heap/pcode/*/cube/series", "/cpu_heap/pcode/{pcode}/cube/series");
//		t.insert("/summary/pcode/*", "/summary/pcode/{pcode}");
//		t.insert("/counter/pcode/*/oid/*/{name}", "/counter/pcode/{pcode}/oid/{oid:.*}/{name}");
//		t.insert("/counter/pcode/*/name", "/counter/pcode/{pcode}/name");
//		t.insert("/stat/pcode/*/oid/*", "/stat/pcode/{pcode}/oid/{oid}");
//		t.insert("/report/pcode/*/daily/summary", "/report/pcode/{pcode}/daily/summary");
//		t.insert("/config/*/*/get", "/config/{pcode}/{oid:.*}/get");
//		t.insert("/config/*/*/set", "/config/{pcode}/{oid:.*}/set");
//		t.insert("/agent/pcode/*/oid/*/show_config", "/agent/pcode/{pcode}/oid/{oid:.*}/show_config");
//		t.insert("/agent/pcode/*/oid/*/add_config", "/agent/pcode/{pcode}/oid/{oid:.*}/add_config");
//		t.insert("/agent/pcode/*/oids", "/agent/pcode/{pcode}/oids");
//		t.insert("/agent/pcode/*/oid/*/remove", "/agent/pcode/{pcode}/oid/{oid}/remove");
//		t.insert("/agent/pcode/*/oid/*/env", "/agent/pcode/{pcode}/oid/{oid}/env");
//		t.insert("/agent/pcode/*/oid/*/threadlist", "/agent/pcode/{pcode}/oid/{oid:.*}/threadlist");
//		t.insert("/agent/pcode/*/oid/*/thread/{threadId}", "/agent/pcode/{pcode}/oid/{oid:.*}/thread/{threadId}");
//		t.insert("*", "${springfox.documentation.swagger.v2.path:/v2/api-docs}");
//		t.insert("/yard/pcode/*/oid", "/yard/pcode/{pcode}/oid");
//		t.insert("/yard/pcode/*/oid/*/oname", "/yard/pcode/{pcode}/oid/{oid}/oname");
//		t.insert("/yard/pcode/*/disk", "/yard/pcode/{pcode}/disk");
//		t.insert("/yard/pcode/*/disk/clear", "/yard/pcode/{pcode}/disk/clear");
//
//
////		System.out.println(t.find("/pcode/123/disk/clear"));
////		System.out.println(t.find("/pcode/1234/disk"));
//		System.out.println(t.find("/config/123/4/get"));
////		Enumeration<PathTree<String>.ENTRY> en = t.entries();
////		while (en.hasMoreElements()) {
////			PathTree<String>.ENTRY e = en.nextElement();
////			System.out.println(e.path() + "=>" + e.value());
////		}
//	}

func main() {
	t := NewPathTree()
	t.Insert("/api/internal/v1/panels/*", "/api/internal/v1/panels/{panelNo}")
	t.Insert("/api/internal/v1/panels/push", "/api/internal/v1/panels/push")
	fmt.Println(t.Find("/api/internal/v1/panels/123"))
	fmt.Println(t.Find("/api/internal/v1/panels/push"))
}
