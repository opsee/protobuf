package types

import (
	"encoding/json"
	"sync"
)

var PermissionsBitmap = &permissionsBitmap{bitmap: make(map[int]string)}

type permissionsBitmap struct {
	bitmap map[int]string
	sync.RWMutex
}

func (p *permissionsBitmap) Get(i int) (string, bool) {
	p.RLock()
	defer p.RUnlock()
	t, ok := p.bitmap[i]
	return t, ok
}

func (p *permissionsBitmap) Length() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.bitmap)
}

func (p *permissionsBitmap) Register(i int, name string) {
	p.Lock()
	defer p.Unlock()
	p.bitmap[i] = name
}

type BitFlags interface {
	// Set bit to 1 at index i
	Set(int)

	// Set bit to 0 at index i
	Clear(int)

	// should return false if outside of bit range
	Test(int) bool

	// return dank bits
	HighBits() []int
}

// Set flag i in permission
func (p *Permission) Set(i int) {
	p.Perm |= (uint64(1) << uint(i))
}

// UnSet flag i in permission
func (p *Permission) Clear(i int) {
	p.Perm &= ^(uint64(1) << uint(i))
}

// Flag i in permission contains 1
func (p *Permission) Test(i int) bool {
	return (p.Perm&(uint64(1)<<uint(i)) > 0)
}

// Returns dank bits
func (p *Permission) HighBits() []int {
	var hb []int
	for i := 0; i < 64; i++ {
		if p.Test(i) {
			hb = append(hb, i)
		}
	}
	return hb
}

func (p *Permission) Permissions() []string {
	var perms []string
	for _, bit := range p.HighBits() {
		if perm, ok := PermissionsBitmap.Get(bit); ok {
			perms = append(perms, perm)
		}
	}
	return perms
}

// Override MarshalJson to return a list of permissions
func (p *Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Permissions())
}
