package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type PermissionTestable interface {
	Run() error
}

type permissionTest struct {
	input *Permission
	i     int
}

// Describes test case for Permission.Set
type SetPermissionTest struct {
	permissionTest
	expected uint64
}

// Testable interface for Permission.Set
func (spt *SetPermissionTest) Run() error {
	spt.input.Set(spt.i)
	if spt.input.Perm != spt.expected {
		return fmt.Errorf("Expected %x, Got %x", spt.expected, spt.input.Perm)
	}
	return nil
}

// Describes test case for Permission.Clear
type ClearPermissionTest struct {
	permissionTest
	expected uint64
}

// Testable interface for Permission.Clear
func (cpt *ClearPermissionTest) Run() error {
	cpt.input.Clear(cpt.i)
	if cpt.input.Perm != cpt.expected {
		return fmt.Errorf("Expected %x, Got %x", cpt.expected, cpt.input.Perm)
	}
	return nil
}

// Describes test case for Permission.Test
type TestPermissionTest struct {
	permissionTest
	expected bool
}

// Testable interface for Permission.Test
func (tpt *TestPermissionTest) Run() error {
	res := tpt.input.Test(tpt.i)
	if res != tpt.expected {
		return fmt.Errorf("Expected %t, Got %t", tpt.expected, res)
	}
	return nil
}

var permissionTests = []PermissionTestable{
	// lsb
	&SetPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0x0)},
			i:     0,
		},
		uint64(0x1),
	},
	// msb
	&SetPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0x0)},
			i:     63,
		},
		uint64(0x8000000000000000),
	},
	&ClearPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0x8000000000000000)},
			i:     63,
		},
		uint64(0x0),
	},
	&ClearPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0x8)},
			i:     3,
		},
		uint64(0x0),
	},
	&TestPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0xf)},
			i:     3,
		},
		true,
	},
	&TestPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0x0)},
			i:     4,
		},
		false,
	},
	&TestPermissionTest{
		permissionTest{
			input: &Permission{Perm: uint64(0xfe)},
			i:     0,
		},
		false,
	},
}

// Run all tests in above list
func TestRunPermissionTests(t *testing.T) {
	for i, z := range permissionTests {
		if err := z.Run(); err != nil {
			t.Error(fmt.Sprintf("Test %d", i), err)
		}
	}
}

func TestPermissionsHighBits(t *testing.T) {
	// register permissions types
	PermissionsBitmap.Register(0, "admin")
	PermissionsBitmap.Register(1, "edit")
	PermissionsBitmap.Register(2, "billing")

	// test marshalling of json 011
	p := &Permission{Perm: uint64(0x3)}

	expected := []int{0, 1}
	res := p.HighBits()
	if !reflect.DeepEqual(expected, res) {
		t.Error(fmt.Sprintf("TestPermissionsHighBits expected %v, got %v", expected, res))
	}
}

func TestPermissionsMarshalJSON(t *testing.T) {
	// register permissions types
	PermissionsBitmap.Register(0, "admin")
	PermissionsBitmap.Register(1, "edit")
	PermissionsBitmap.Register(2, "billing")

	// test marshalling of json 011
	expected := []string{"admin", "edit"}
	jb, _ := json.Marshal(&Permission{Perm: uint64(0x3)})

	var res []string
	_ = json.Unmarshal(jb, &res)
	if !reflect.DeepEqual(expected, res) {
		t.Error(fmt.Sprintf("TestPermissionsMarshalJSON expected %v, got %v", expected, res))
	}
}
