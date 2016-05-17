package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	proto "github.com/gogo/protobuf/proto"
)

var AnyTypeRegistry = &AnyTypes{registry: make(map[string]reflect.Type)}

type AnyTypes struct {
	registry map[string]reflect.Type
	sync.RWMutex
}

func (this *AnyTypes) Get(name string) (reflect.Type, bool) {
	this.RLock()
	defer this.RUnlock()
	t, ok := this.registry[name]
	return t, ok
}

func (this *AnyTypes) RegisterAny(name string, t reflect.Type) {
	this.Lock()
	this.registry[name] = t
	this.Unlock()
}

// UnmarshalAny unmarshals an Any object based on its TypeUrl type hint.
func UnmarshalAny(any *Any) (interface{}, error) {
	class := any.TypeUrl
	bytes := any.Value

	if class, ok := AnyTypeRegistry.Get(class); !ok {
		return nil, fmt.Errorf("Couldn't find type in registry")
	} else {
		instance := reflect.New(class).Interface()
		err := proto.Unmarshal(bytes, instance.(proto.Message))
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
}

// MarshalAny uses reflection to marshal an interface{} into an Any object and
// sets up its TypeUrl type hint.

func MarshalAny(i interface{}) (*Any, error) {
	msg, ok := i.(proto.Message)
	if !ok {
		err := fmt.Errorf("Unable to convert to proto.Message: %v", i)
		return nil, err
	}
	bytes, err := proto.Marshal(msg)

	if err != nil {
		return nil, err
	}

	return &Any{
		TypeUrl: reflect.ValueOf(i).Elem().Type().Name(),
		Value:   bytes,
	}, nil
}

// marshal any to json
func (a *Any) MarshalJSON() ([]byte, error) {
	obj, err := UnmarshalAny(a)
	if err != nil {
		return []byte{}, err
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
