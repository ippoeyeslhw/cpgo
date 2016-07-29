package cpgo

import (
	//	"fmt"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// 사이보스플러스 GetDibStatus
func (c *CpClass) GetDibStatus() int16 {
	r := oleutil.MustCallMethod(c.obj, "GetDibStatus")
	if ret, ok := r.Value().(int16); ok {
		return int16(ret)
	}
	panic(r)
}

// 사이보스플러스 SetInputValue함수 Wrapper
func (c *CpClass) SetInputValue(typ int, val interface{}) {
	_ = oleutil.MustCallMethod(c.obj, "SetInputValue", typ, val)
}

// 사이보스플러스 BlockRequest 함수 Wrapper
func (c *CpClass) BlockRequest() {
	_ = oleutil.MustCallMethod(c.obj, "BlockRequest")
}

// 사이보스플러스 BlockRequest 함수 Wrapper
func (c *CpClass) Request() {
	if c.evnt == nil {
		panic("err")
	}
	_ = oleutil.MustCallMethod(c.obj, "Request")
}

// 사이보스플러스 Subscribe 함수 Wrapper
func (c *CpClass) Subscribe() {
	if c.evnt == nil {
		panic("err")
	}
	_ = oleutil.MustCallMethod(c.obj, "Subscribe")
}

// 사이보스플러스 SubscribeLastest 함수 Wrapper
func (c *CpClass) SubscribeLastest() {
	if c.evnt == nil {
		panic("err")
	}
	_ = oleutil.MustCallMethod(c.obj, "SubscribeLastest")
}

// 사이보스플러스 Unsubscribe 함수 Wrapper
func (c *CpClass) Unsubscribe() {
	_ = oleutil.MustCallMethod(c.obj, "Unsubscribe")
}

// 사이보스플러스 GetHeaderValue Wrapper
func (c *CpClass) GetHeaderValue(typ int) (result *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetHeaderValue", typ)
}

// 사이보스플러스 GetDataValue Wrapper
func (c *CpClass) GetDataValue(typ int, idx int) (result *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetDataValue", typ, idx)
}
