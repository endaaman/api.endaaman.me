package models


type Base struct {
	identified bool
}

func (a *Base) Identify() {
	a.identified = true
}

func (a *Base) Identified() bool {
	return a.identified
}

