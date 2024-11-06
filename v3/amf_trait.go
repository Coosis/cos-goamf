package v3

import (
	// "fmt"
)

type AmfTraitSet struct {
	ClassName string
	Traits []string
	IsDynamic bool
	IsExternalizable bool
}

// TODO: Increase efficiency
func(ts *AmfTraitSet) Hash() string {
	res := ts.ClassName
	for _, trait := range ts.Traits {
		res += trait
	}
	dyn := "0"
	if ts.IsDynamic {
		dyn = "1"
	}
	ext := "0"
	if ts.IsExternalizable {
		ext = "1"
	}
	res += dyn + ext
	return res
}

func NewAmfTraitSet() *AmfTraitSet {
	return &AmfTraitSet{
		ClassName: "",
		Traits: make([]string, 0),
		IsDynamic: false,
		IsExternalizable: false,
	}
}
