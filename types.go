package edn

import (
	//"regexp"
	//"strings"
)

type Eq interface {
	Equal(other interface{}) bool
}

type Metaer interface {
	WithMeta(*PMap) Meta
	Meta() *PMap
}

type Named interface {
	Ns() string
	Name() string
}


type Symbol struct {
	ns, name string
	meta *PMap
	//hash int
	_str string // for caching the value of 
}

func NewSym(ns, name string) Symbol {
	var str string
	if ns != "" {
		str = ns+"/"+name
	}else{
		str = name
	}
	return Symbol{
		ns:		ns, 
		name:	name,
		_str:		str	
	}
}

func (sym Symobl) String() string {
	return sym._str
}

func (sym Symbol) Equal(other interface{}) bool {
	// if other isn't a symbol, it (obviously) can't be equal
	so, ok := other.(Symbol)
	return ok && so.Name() == sym.name && so.Ns() == sym.ns
}

func (sym Symbol) WithMeta(m *PMap) Meta {
	return Symbol{ns:ns, name:name, meta:m}
}

// some accessors, to fill out interfaces above
func (sym Symbol) Meta() *PMap { return sym.meta }
func (sym Symbol) Name() string { return sym.name }
func (sym Symbol) Ns() string { return sym.ns }

// todo: implement immutable mappings
// maybe with https://github.com/steveyen/gtreap
type PMap struct{}


