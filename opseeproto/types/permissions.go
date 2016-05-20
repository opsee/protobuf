package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
)

// the permission which corresponds to opsee administrator
const OpseeAdmin = "opsee_admin"

var PermissionsRegistry = NewPermsRegistry()

func NewPermissionsBitmap() *PermissionsBitmap {
	return &PermissionsBitmap{Bitmap: make(map[int]string)}
}

type PermissionsBitmap struct {
	Bitmap map[int]string
	sync.RWMutex
}

func (p *PermissionsBitmap) Get(i int) (string, bool) {
	p.RLock()
	defer p.RUnlock()
	t, ok := p.Bitmap[i]
	return t, ok
}

func (p *PermissionsBitmap) Length() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.Bitmap)
}

func (p *PermissionsBitmap) Register(i int, name string) {
	p.Lock()
	defer p.Unlock()
	p.Bitmap[i] = name
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

// Returns a list of permissions we have
func (p *Permission) Permissions() []string {
	var perms []string
	r, ok := PermissionsRegistry.Get(p.Name)
	if !ok {
		return []string{}
	}
	for _, bit := range p.HighBits() {
		if perm, ok := r.Get(bit); ok {
			perms = append(perms, perm)
		}
	}
	return perms
}

// Checks permissions map for permission names, returns errors for those that do not exist
func (p *Permission) HasPermissions(pnames ...string) map[string]bool {
	hasPerms := make(map[string]bool)
	retPerms := make(map[string]bool)
	for _, p := range p.Permissions() {
		hasPerms[p] = true
	}

	for _, name := range pnames {
		if _, ok := hasPerms[name]; ok {
			retPerms[name] = true
		} else {
			retPerms[name] = false
		}
	}
	return retPerms
}

// Checks permissions map for permission names, returns errors for those that do not exist
func (p *Permission) CheckPermissions(pnames ...string) map[string]error {
	permErrs := make(map[string]error)
	for name, has := range p.HasPermissions(pnames...) {
		if !has {
			permErrs[name] = NewPermissionsError(name)
		} else {
			permErrs[name] = nil
		}
	}
	return permErrs
}

// Override MarshalJson to return a list of permissions
func (p *Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Permissions())
}

func (p *Permission) Scan(src interface{}) error {
	switch value := src.(type) {
	case int:
		p.Perm = uint64(value)
	case int32:
		p.Perm = uint64(value)
	case int64:
		p.Perm = uint64(value)
	default:
		return fmt.Errorf("invalid type")
	}

	return p.Validate()
}

func (p *Permission) Validate() error {
	return nil
}

func (p *Permission) Value() (driver.Value, error) {
	return int64(p.Perm), nil
}

func NewPermissionsError(pname string) *Error {
	return NewError("PermissionsError", fmt.Sprintf("Not authorized: %s", pname))
}
