package component

import "testing"

func TestNewRegistrar(t *testing.T) {
	r := NewRegistrar()
	if r == nil {
		t.Error("Couldn't create registrar!")
	}
}

func TestRegistrarImpl_Register(t *testing.T) {

}
