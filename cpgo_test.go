package cpgo

import (
	"fmt"
	"testing"
	"time"

	ole "github.com/go-ole/go-ole"
)

func TestCoinitialize(t *testing.T) {
	ole.CoInitialize(0)
}

// 객체 생성 헤제 테스트
func TestCpClassStruct(t *testing.T) {
	tmp := &CpClass{}
	if tmp == nil {
		t.Error("CpClass struct new error")
	}
	fmt.Println(tmp)
	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()

	fmt.Println(tmp)

}

// block request 테스트
func TestCpBlockRequest(t *testing.T) {
	tmp := &CpClass{}
	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()

	tmp.SetInputValue(0, "A000270")
	fmt.Println("bf blockrequest")
	tmp.BlockRequest()
	fmt.Println("bf getheadervalue")
	v := tmp.GetHeaderValue(1)
	fmt.Println(v.Value())
}

type RqTestStruct struct {
}

func (t *RqTestStruct) Received(c *CpClass) {
	fmt.Println(c.GetHeaderValue(1).Value())
	c.UnbindEvent() // 이벤트 바로 헤제
}

// request 테스트 (비동기)
func TestCpRequest(t *testing.T) {
	tmp := &CpClass{}
	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()
	fmt.Println(tmp)

	evnt := &RqTestStruct{}
	tmp.BindEvent(evnt)
	fmt.Println(tmp)

	tmp.SetInputValue(0, "A000660")
	tmp.Request()

	fmt.Println("wait event")
	var m ole.Msg
	for tmp.evnt.ref != 0 {
		time.Sleep(1)
		ole.GetMessage(&m, 0, 0, 0)
		ole.DispatchMessage(&m)
	}
}

// PumpWaitingMessage 텟트
func TestPumpWaitingMessage(t *testing.T) {
	tmp := &CpClass{}
	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()
	fmt.Println(tmp)

	evnt := &RqTestStruct{}
	tmp.BindEvent(evnt)
	fmt.Println(tmp)

	tmp.SetInputValue(0, "A000660")
	tmp.Request()

	fmt.Println("wait event")
	fmt.Println(pPeekMessage)
	for tmp.evnt.ref != 0 {
		time.Sleep(1)
		PumpWaitingMessage()
	}

}

func TestCoUninitialize(t *testing.T) {
	ole.CoUninitialize()
}
