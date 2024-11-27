package main

import (
	"errors"
	"fmt"

	"github.com/okieoth/goptional"
)

type DummyEnum int64

const (
	V1 DummyEnum = iota
	V2
	V3
)

func (s DummyEnum) String() string {
	switch s {
	case V1:
		return "v1"
	case V2:
		return "v2"
	default:
		return "v3"
	}
}

func (s DummyEnum) ValueFromStr(str string) error {
	switch str {
	case "value1":
		s = V1
	case "v2":
		s = V2
	case "v3":
		s = V3
	default:
		return fmt.Errorf("input not part of the enum: %v", str)
	}
	return nil
}

type Demo struct {
	Id            int64
	List          []DummyEnum
	OptionalMap   goptional.Optional[map[string]int32]
	OptionalList  goptional.Optional[[]string]
	OptionalList2 goptional.Optional[[][]DummyEnum]
	OptionalEnum  goptional.OptionalEnum[DummyEnum]
}

type DemoBuilder struct {
	id            goptional.Optional[int64]
	list          goptional.Optional[[]DummyEnum]
	optionalMap   goptional.Optional[map[string]int32]
	optionalList  goptional.Optional[[]string]
	optionalList2 goptional.Optional[[][]DummyEnum]
	optionalEnum  goptional.OptionalEnum[DummyEnum]
}

func NewDemoBuilder() *DemoBuilder {
	var inst DemoBuilder
	return &inst
}

func (b *DemoBuilder) Id(v int64) *DemoBuilder {
	b.id.Set(v)
	return b
}

func (b *DemoBuilder) List(v []DummyEnum) *DemoBuilder {
	b.list.Set(v)
	return b
}

func (b *DemoBuilder) OptionalMap(v map[string]int32) *DemoBuilder {
	b.optionalMap.Set(v)
	return b
}

func (b *DemoBuilder) OptionalList(v []string) *DemoBuilder {
	b.optionalList.Set(v)
	return b
}

func (b *DemoBuilder) OptionalList2(v [][]DummyEnum) *DemoBuilder {
	b.optionalList2.Set(v)
	return b
}

func (b *DemoBuilder) OptionalEnum(v DummyEnum) *DemoBuilder {
	b.optionalEnum.Set(v)
	return b
}

func (b *DemoBuilder) Build() (Demo, error) {
	var r Demo
	if v, set := b.id.Get(); set {
		r.Id = v
	} else {
		return r, errors.New("Missing initialization for Demo.Id")
	}
	if v, set := b.list.Get(); set {
		r.List = v
	} else {
		return r, errors.New("Missing initialization for Demo.List")
	}
	if b.optionalMap.IsSet() {
		r.OptionalMap = b.optionalMap
	}
	if b.optionalList.IsSet() {
		r.OptionalList = b.optionalList
	}
	if b.optionalList2.IsSet() {
		r.OptionalList2 = b.optionalList2
	}
	if b.optionalEnum.IsSet() {
		r.OptionalEnum = b.optionalEnum
	}
	return r, nil
}

func main() {
	if demoInst, err := NewDemoBuilder().
		Id(1).
		List([]DummyEnum{V1, V2}).
		OptionalMap(map[string]int32{"A": 1, "B": 2}).
		Build(); err == nil {
		fmt.Printf("Demo: %v", demoInst)
	} else {
		fmt.Printf("Error while creating Demo object: %v", err)
	}
}
