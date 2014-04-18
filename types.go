package edn

import (
	//"regexp"
	//"strings"
)

type Eq interface {
	func Equal(other interface{}) bool
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
	hash int
	_str string
}

func NewSym(ns, name string) Symbol {
	return Symbol{ns:ns, name:name}
}

func (sym Symobl) String() string {
	if sym._str == "" {
		if sym.ns != "" {
			sym._str = sym.ns + sym.name
		}else{
			sym._str = sym.name
		}
	}
	return sym._str
}

func (sym Symbol) Equal(other interface{}) bool {
	so, ok := other.(Symbol)
	return ok && so.name == sym.name && so.ns == sym.ns
}

func (sym Symbol) WithMeta(m *PMap) Meta {
	return Symbol{ns:ns, name:name, meta:m}
}

func (sym Symbol) Meta() *PMap { return sym.meta }
func (sym Symbol) Name() string { return sym.name }
func (sym Symbol) Ns() string { return sym.ns }



