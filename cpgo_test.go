package cpgo

import (
	"fmt"
	"testing"
	"time"

	ole "github.com/go-ole/go-ole"
)

// 객체 생성 헤제 테스트
func TestCpClassStruct(t *testing.T) {
	tmp := &CpClass{}
	if tmp == nil {
		t.Error("CpClass struct new error")
	}
	fmt.Println(tmp)

	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()

	fmt.Println(tmp)

}

type RqTestStruct struct {
}

func (t *RqTestStruct) Received(c *CpClass) {
	fmt.Println(c.GetHeaderValue(1).Value())
	c.UnbindEvent() // 이벤트 바로 헤제
}

// request 테스트
func TestCpClassMethod(t *testing.T) {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	tmp := &CpClass{}
	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()

	tmp.SetInputValue(0, "A000270")
	fmt.Println("bf blockrequest")
	tmp.BlockRequest()
	fmt.Println("bf getheadervalue")
	v := tmp.GetHeaderValue(1)
	fmt.Println(v.Value())

	// event instance, and binding events
	fmt.Println("bind event")
	evnt := &RqTestStruct{}
	tmp.BindEvent(evnt)

	tmp.SetInputValue(0, "A000660")
	fmt.Println("rq")
	tmp.Request()

	fmt.Println("wait event")
	var m ole.Msg
	for tmp.evnt.ref != 0 {
		time.Sleep(1)
		ole.GetMessage(&m, 0, 0, 0)
		ole.DispatchMessage(&m)
	}
}

func TestCpRequest(t *testing.T) {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	tmp := &CpClass{}
	tmp.Create("DSCBO1.StockMst")
	defer tmp.Release()

	evnt := &RqTestStruct{}
	tmp.BindEvent(evnt)

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

/*
=== RUN   TestCpClassStruct
&{<nil> <nil> <nil> <nil> <nil> 0}
&{0x3276b20c 0x3276b20c <nil> <nil> <nil> 0}
--- PASS: TestCpClassStruct (0.01s)
=== RUN   TestCpClassMethod
bf blockrequest
bf getheadervalue
기아차
bind event
rq
wait event
SK하이닉스
--- PASS: TestCpClassMethod (0.08s)
PASS
*/
