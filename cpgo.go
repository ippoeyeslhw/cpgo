package cpgo

import (
	//	"fmt"
	"sync"
	"syscall"
	"unsafe"

	ole "github.com/go-ole/go-ole"
	//"github.com/go-ole/go-ole/oleutil"
)

// peekmessage 로드
var (
	user32, _       = syscall.LoadLibrary("user32.dll")
	pPeekMessage, _ = syscall.GetProcAddress(user32, "PeekMessageW")
	//pDispatchMessage, _ = syscall.GetProcAddress(user32, "DispatchMessage")
)

// 사이보스플러스의 콜백메서드 인터페이스
type Receiver interface {
	Received(*CpClass)
}

// 사이보스플러스 객체를 구성하는 데이터묶음
type CpClass struct {
	unk  *ole.IUnknown
	obj  *ole.IDispatch
	evnt *dispCpEvent

	// for event
	cb     Receiver
	point  *ole.IConnectionPoint
	cookie uint32
}

// 이벤트 수신을 위한 구조체
type dispCpEvent struct {
	lpVtbl *dispCpEventVtbl
	ref    int32
	host   *CpClass
}

// 가상함수 테이블
type dispCpEventVtbl struct {
	// IUnknown
	pQueryInterface uintptr
	pAddRef         uintptr
	pRelease        uintptr
	// IDispatch
	pGetTypeInfoCount uintptr
	pGetTypeInfo      uintptr
	pGetIDsOfNames    uintptr
	pInvoke           uintptr
}

// 사이보스플러스 객체 생성
func (c *CpClass) Create(name string) {
	// clsid 구함
	clsid, err := ole.CLSIDFromString(name)
	if err != nil {
		panic(err)
	}
	// _IDib interface
	iid, err := ole.CLSIDFromString("{33518a10-0931-11d4-8231-00105a7c4f8c}")
	if err != nil {
		panic(err)
	}
	// unknown
	c.unk, err = ole.CreateInstance(clsid, ole.IID_IUnknown)
	if err != nil {
		panic(err)
	}
	// get obj
	//c.obj, err = c.unk.QueryInterface(ole.IID_IDispatch)
	c.obj, err = c.unk.QueryInterface(iid)
	if err != nil {
		panic(err)
	}
}

// 객체 헤제
func (c *CpClass) Release() {
	if c.unk != nil {
		c.unk.Release()
		c.unk = nil
	}
	if c.obj != nil {
		c.obj.Release()
		c.obj = nil
	}
	if c.evnt != nil {
		//c.evnt.Release()
		dispRelease((*ole.IUnknown)(unsafe.Pointer(c.evnt)))
		c.evnt = nil
		if c.point != nil {
			c.UnbindEvent()
		}
	}
}

// 이벤트 지정
func (c *CpClass) BindEvent(callback Receiver) {

	if c.evnt == nil {
		// Callback method binding
		evnt := (new(dispCpEvent))
		evnt.lpVtbl = (new(dispCpEventVtbl))
		evnt.lpVtbl.pQueryInterface = syscall.NewCallback(dispQueryInterface)
		evnt.lpVtbl.pAddRef = syscall.NewCallback(dispAddRef)
		evnt.lpVtbl.pRelease = syscall.NewCallback(dispRelease)
		evnt.lpVtbl.pGetTypeInfoCount = syscall.NewCallback(dispGetTypeInfoCount)
		evnt.lpVtbl.pGetTypeInfo = syscall.NewCallback(dispGetTypeInfo)
		evnt.lpVtbl.pGetIDsOfNames = syscall.NewCallback(dispGetIDsOfNames)
		evnt.lpVtbl.pInvoke = syscall.NewCallback(dispInvoke)
		evnt.host = c
		// assign event
		c.evnt = evnt
	}
	c.cb = callback

	// get event iid
	dibevnt_iid, err := ole.CLSIDFromString("{B8944520-09C3-11D4-8232-00105A7C4F8C}")
	if err != nil {
		panic(err)
	}

	if c.point != nil {
		// 이미 포인트가 지정되어 있었으면?
		c.UnbindEvent()
	}
	// connectionpoint container
	unknown_con, err := c.obj.QueryInterface(ole.IID_IConnectionPointContainer)
	if err != nil {
		panic(err)
	}

	// get point
	container := (*ole.IConnectionPointContainer)(unsafe.Pointer(unknown_con))
	var point *ole.IConnectionPoint
	err = container.FindConnectionPoint(dibevnt_iid, &point)
	if err != nil {
		panic(err)
	}

	// Advise
	cookie, err := point.Advise((*ole.IUnknown)(unsafe.Pointer(c.evnt)))
	//_, err = point.Advise((*ole.IUnknown)(unsafe.Pointer(evnt)))
	container.Release()
	if err != nil {
		point.Release()
		panic(err)
	}
	c.point = point
	c.cookie = cookie
}

// 이벤트 헤제
func (c *CpClass) UnbindEvent() {
	if c.point != nil {
		c.point.Unadvise(c.cookie)
		c.point.Release()
		c.point = nil
		c.cookie = 0
	}
}

// 이하 콜백 이벤트 바인딩하기 위한 함수 선언들
func dispQueryInterface(this *ole.IUnknown, iid *ole.GUID, punk **ole.IUnknown) uint32 {
	s, _ := ole.StringFromCLSID(iid)
	*punk = nil
	if ole.IsEqualGUID(iid, ole.IID_IUnknown) ||
		ole.IsEqualGUID(iid, ole.IID_IDispatch) {
		dispAddRef(this)
		*punk = this
		return ole.S_OK
	}
	// _DibEvents GUID
	if s == "{B8944520-09C3-11D4-8232-00105A7C4F8C}" {
		dispAddRef(this)
		*punk = this
		return ole.S_OK
	}

	return ole.E_NOINTERFACE
}

func dispAddRef(this *ole.IUnknown) int32 {
	pthis := (*dispCpEvent)(unsafe.Pointer(this))
	pthis.ref++
	return pthis.ref
}

func dispRelease(this *ole.IUnknown) int32 {
	pthis := (*dispCpEvent)(unsafe.Pointer(this))
	pthis.ref--
	return pthis.ref
}
func dispGetIDsOfNames(args *uintptr) uint32 {
	p := (*[6]int32)(unsafe.Pointer(args))
	//this := (*ole.IDispatch)(unsafe.Pointer(uintptr(p[0])))
	//iid := (*ole.GUID)(unsafe.Pointer(uintptr(p[1])))
	wnames := *(*[]*uint16)(unsafe.Pointer(uintptr(p[2])))
	namelen := int(uintptr(p[3]))
	//lcid := int(uintptr(p[4]))
	pdisp := *(*[]int32)(unsafe.Pointer(uintptr(p[5])))
	for n := 0; n < namelen; n++ {
		s := ole.UTF16PtrToString(wnames[n])
		println(s)
		pdisp[n] = int32(n)
	}
	return ole.S_OK
}
func dispGetTypeInfoCount(this *ole.IUnknown, pcount *int) uint32 {
	if pcount != nil {
		*pcount = 0
	}
	return ole.S_OK
}

func dispGetTypeInfo(this *ole.IUnknown, namelen int, lcid int) uint32 {
	return ole.E_NOTIMPL
}
func dispInvoke(this *ole.IDispatch, dispid int, riid *ole.GUID, lcid int, flags int16, dispparams *ole.DISPPARAMS, result *ole.VARIANT, pexcepinfo *ole.EXCEPINFO, nerr *uint) uintptr {
	pthis := (*dispCpEvent)(unsafe.Pointer(this))
	if dispid == 1 {
		if pthis.host.cb != nil {
			// instance callback
			pthis.host.cb.Received(pthis.host)
			return ole.S_OK
		}
	}
	return ole.E_NOTIMPL
}

//

func PeekMessage(msg *ole.Msg, hwnd uint32, MsgFilterMin uint32, MsgFilterMax uint32, RemoveMsg uint32) (ret int32, err error) {
	r0, _, err := syscall.Syscall6(uintptr(pPeekMessage), 5,
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(MsgFilterMin),
		uintptr(MsgFilterMax),
		uintptr(RemoveMsg),
		0)

	ret = int32(r0)
	return
}

func PumpWaitingMessage() int32 {
	ret := int32(0)

	var msg ole.Msg

	mutex := &sync.Mutex{}
	mutex.Lock()
	for {
		r, _ := PeekMessage(&msg, 0, 0, 0, 1)
		if r == 0 {
			break
		}
		if msg.Message == 0x0012 { // WM_QUIT
			ret = int32(1)
			break
		}
		ole.DispatchMessage(&msg)
	}
	mutex.Unlock()
	return ret
}
