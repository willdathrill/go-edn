package edn

import (
	//"regexp"
	//"strings"
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
	}
	return k
}

func (k *Keyword) String() string { return k._str }
func (k *Keyword) Name() string { return k.Name() }


// todo: implement immutable mappings
// maybe with https://github.com/steveyen/gtreap
type PMap struct{}


