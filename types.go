package edn

import (
	"reflect"
	"sync"
)

type Eq interface {
	Equal(other interface{}) bool
}

type Named interface {
	Ns() string
	Name() string
}

type Hasher interface {
	Hash() int32
	Hasheq(other Hasher) bool
}

type Counted interface {
	Count() int
}

type Seq interface {
	First() interface{}
	Next() Seq
	More() Seq
	Cons(interface{}) Seq
}

type Stack interface {
	Peek() interface{}
	Pop() Stack
}

type List interface {
	Seq
	// todo: formalize equality
	//Eq
	Stack
	Contains(o interface{}) bool
	// todo: work out func signature, maybe Coll inter?
	//ContainsAll
	Index(interface{}) (int, bool)
	LastIndex(interface{}) (int, bool)
	EmptyP() bool
	SubList(int, int) List
}

type Symbol struct {
	ns, name string
	hash int
	_str string // for caching the value of 
}

func NewSym(ns, name string) Symbol {
	var str string
	if ns != "" {
		str = ns+"/"+name
	}else{
		str = name
	}
	return Symbol{ //todo: implement hashing
		ns:		ns, 
		name:	name,
		_str:		str	
	}
}

func (sym Symobl) String() string { return sym._str }
func (sym Symbol) Name() string { return sym.name }
func (sym Symbol) Ns() string { return sym.ns }

func (sym Symbol) Equal(other interface{}) bool {
	so, ok := other.(Symbol)
	return ok && so.Name() == sym.name && so.Ns() == sym.ns
}

// Keyword implementation
var (
	keyTable = make(map[Symbol]*Keyword) //cache of Keywords
	keyLock sync.Mutex // lock ^
)

type Keyword struct {
	Symbol
	hash int
	_str string
}

func NewKwd(ns, name string) *Keyword {
	keyLock.Lock()
	defer keyLock.Unlock()
	sym := NewSym(ns,name)
	k, ok := keyTable[sym]
	if !ok { // todo: implement hashing
		k = &Keyword{Symbol:sym,_str:":"+sym.String()}
		keyTable[sym] = k
	}
	return k
}

func (k *Keyword) String() string { return k._str }
func (k *Keyword) Equal(other interface{}) bool {
	ko, ok := other.(*Keyword)
	return ok && ko == k
	// reference equality deliberate: cached in keyTable
}

// Persistent List

type PList struct {
	first interface{}
	next *PList
	count int
}

func NewList1(fst interface{}) *PList {
	return &PList{first: fst, next: nil, count: 1}
}

func NewList(fst interface{}, rst *PList, count int) *PList {
	return &PList{first: fst, rest: rst, count: count}
}

var emptyList = 

func (p *PList) Count() int {
	return p.count
}

func (p *PList) FIrst() interface{} {
	return p.first
}

func (p *PList) Next() Seq {
	if p.count == 1 { return nil }
	return p.next
}

func (p *PList) More() Seq {
	s := p.Next()
	if s == nil { s = NewList(nil,nil,0) }
	return s
}

func (p *PList) Cons(o interface{}) Seq {
	return NewList(o, p, p.count+1)
}

func (p *PList) Peek() interface{} { return p.first }

func (p *PList) Pop() Stack {
	if p.rest == nil {
		return NewList(nil,nil,0)
	}
	return p.rest
}

func (p *PList) Empty() *PList { return NewList(nil,nil,0) }

func (p *PList) EmptyP() bool {
	return p.first == nil && p.next == nil
}

func (p *PList) reify() []interface{} {
	np := p
	ret := make([]interface{}, np.count)
	for i := range ret {
		ret[i] = np.First()
		np = np.Next()
		if np == nil {
			break
		}
	}
	return ret
}

func sliceToList(all ...interface{}) *PList {
	p := emptyList
	for _,v := range all {
		p = p.Cons(v)
	}
	return p
}

func (p *PList) SubList(from, to int) List {
	return sliceToList(p.reify()[from:to]...)
}

func (p *PList) Index(o interface{}) (int, bool) {
	np := p
	for i := 0; np != nil; i++ {
		// todo: *might* need to fix this to enable nesting
		if reflect.DeepEqual(np.first, o) {
			return i, true
		}
		np := np.Next()
	}
	return -1, false // <-- indicating not found
}

func (p *PList) LastIndex(o interface{}) (int, bool) {
	r := p.reify()
	for i := len(r); i > 0; i-- {
		// todo: ref Index ^
		if reflect.DeepEqual(r[i], o) {
			return i, true
		}
	}
	return -1, false
}

func (p *PList) Contains(o interface{}) bool {
	for np := p; np != nil; np = np.Next() {
		if reflect.DeepEqual(np.First(), o) {
			return true
		}
	}
	return false
}

// Vector implementation
/* todo: how do we ensure transients aren't edited after
	call to persistent?
type vecNode struct{
	sema chan struct{}
	arr []interface{}
}

type tVec struct {
	cnt uint
	shift int
	root vecNode
	tail []interface{}
}

func (t *tVec) Count() int {
	return int(t.cnt)
}
*/