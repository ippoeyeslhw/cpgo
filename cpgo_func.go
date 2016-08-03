package cpgo

import (
	//	"fmt"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// 사이보스플러스 Property Getter 메서드입니다.
// Continue 프로퍼티 값을 얻어옵니다.
func (c *CpClass) GetContinue() (r *ole.VARIANT) {
	return oleutil.MustGetProperty(c.obj, "Continue")
}

// 사이보스플러스 GetDibStatus 메서드 Wrapper
func (c *CpClass) GetDibStatus() int16 {
	r := oleutil.MustCallMethod(c.obj, "GetDibStatus")
	if ret, ok := r.Value().(int16); ok {
		return int16(ret)
	}
	panic(r)
}

// 사이보스플러스 SetInputValue 메서드 Wrapper
func (c *CpClass) SetInputValue(typ int, val interface{}) {
	_ = oleutil.MustCallMethod(c.obj, "SetInputValue", typ, val)
}

// 사이보스플러스 BlockRequest 메서드 Wrapper
func (c *CpClass) BlockRequest() {
	_ = oleutil.MustCallMethod(c.obj, "BlockRequest")
}

// 사이보스플러스 BlockRequest 메서드 Wrapper
func (c *CpClass) Request() {
	if c.evnt == nil {
		panic("err")
	}
	_ = oleutil.MustCallMethod(c.obj, "Request")
}

// 사이보스플러스 Subscribe 메서드 Wrapper
func (c *CpClass) Subscribe() {
	if c.evnt == nil {
		panic("err")
	}
	_ = oleutil.MustCallMethod(c.obj, "Subscribe")
}

// 사이보스플러스 SubscribeLastest 메서드 Wrapper
func (c *CpClass) SubscribeLastest() {
	if c.evnt == nil {
		panic("err")
	}
	_ = oleutil.MustCallMethod(c.obj, "SubscribeLastest")
}

// 사이보스플러스 Unsubscribe 메서드 Wrapper
func (c *CpClass) Unsubscribe() {
	_ = oleutil.MustCallMethod(c.obj, "Unsubscribe")
}

// 사이보스플러스 GetHeaderValue 메서드 Wrapper
func (c *CpClass) GetHeaderValue(typ int) (result *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetHeaderValue", typ)
}

// 사이보스플러스 GetDataValue  메서드 Wrapper
func (c *CpClass) GetDataValue(typ int, idx int) (result *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetDataValue", typ, idx)
}

// CpUtil

// 사이보스플러스 CpUtil.CpCybos ServerType Property getter
func (c *CpClass) GetIsConnect() (r *ole.VARIANT) {
	return oleutil.MustGetProperty(c.obj, "IsConnect")
}

// 사이보스플러스 CpUtil.CpCybos ServerType Property getter
func (c *CpClass) GetServerType() (r *ole.VARIANT) {
	return oleutil.MustGetProperty(c.obj, "ServerType")
}

// 사이보스플러스 CpUtil.CpCybos LimitRequestRemainTime Property getter
func (c *CpClass) GetLimitRequestRemainTime() (r *ole.VARIANT) {
	return oleutil.MustGetProperty(c.obj, "LimitRequestRemainTime")
}

// 사이보스플러스 CpUtil.CpCybos GetLimitRemainCount 메서드 Wrapper
func (c *CpClass) GetLimitRemainCount(typ int) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetLimitRemainCount", typ)
}

// 사이보스플러스 CpUtil.CpStockCode  CodeToName 메서드 Wrapper
func (c *CpClass) CodeToName(cod string) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "CodeToName", cod)
}

// 사이보스플러스 CpUtil.CpStockCode  NameToCode 메서드 Wrapper
func (c *CpClass) NameToCode(nm string) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "NameToCode", nm)
}

// 사이보스플러스 CpUtil.CpStockCode  CodeToFullCode 메서드 Wrapper
func (c *CpClass) CodeToFullCode(cod string) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "CodeToFullCode", cod)
}

// 사이보스플러스 CpUtil.CpStockCode  FullCodeToName 메서드 Wrapper
func (c *CpClass) FullCodeToName(fullcod string) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "FullCodeToName", fullcod)
}

// 사이보스플러스 CpUtil.CpStockCode  FullCodeToCode 메서드 Wrapper
func (c *CpClass) FullCodeToCode(fullcod string) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "FullCodeToCode", fullcod)
}

// 사이보스플러스 CpUtil.CpStockCode  CodeToIndex 메서드 Wrapper
func (c *CpClass) CodeToIndex(cod string) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "CodeToIndex", cod)
}

// 사이보스플러스 CpUtil.CpStockCode  GetCount 메서드 Wrapper
func (c *CpClass) GetCount() (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetCount")
}

// 사이보스플러스 CpUtil.CpStockCode  GetData 메서드 Wrapper
func (c *CpClass) GetData(typ int, idx int) (r *ole.VARIANT) {
	return oleutil.MustCallMethod(c.obj, "GetData", typ, idx)
}
