# cpgo
=====

cybosplus(증권API) + golang 을 이용하여 
투자 분석 및 시스템트레이딩을 할수 있도록 만든 간단한 
Wrapper 라이브러리 입니다.

# 필요한것들
 * 윈도우 운영체제
 * 사이보스플러스
 * golang (32bits)
 * [go-ole패키지](https://github.com/go-ole/go-ole)

# 설치

```
go get github.com/ippoeyeslhw/cpgo

```

# 설명
cpgo 는 사이보스플러스를 go언어에서 사용할수 있도록 Wrapping 한 것입니다.
[go-ole패키지](https://github.com/go-ole/go-ole)를 기반으로 작성되었습니다.


### Import
go-ole패키지를 기반으로 작성되었으므로 많은 부분을 의존합니다.
반드시 임포트하여야합니다.
```go
import (
	ole "github.com/go-ole/go-ole"
	"github.com/ippoeyeslhw/cpgo"
)
```

### 프로그램 시작
COM Object 프로그래밍을 할 것이므로 반드시 cointialize, uncoinitialize 호출
해야합니다.
```go
	ole.CoInitialize(0)
	defer ole.CoUninitialize()
```

### 객체생성
도움말을 열어봤을때 좌측에 보이는 주요 네가지 라이브러리인
 * CpDib : 혹은 DSCBO1
 * CpSysDib
 * CpTrade
 * CpUtil

각각에 속하는 coclass들을 "."으로 연결하여 객체를 생성할수 있습니다.
Create 메서드로 COM객체 생성을 Release 메서드로 헤제를 할수 있습니다.
```go
	stkmst := &cpgo.CpClass{}
	stkmst.Create("CpDib.StockMst")
	defer stkmst.Release()
```

### 주요 인터페이스
기본적인 동작을 위해 구현된 주요 인터페이스 메서드들입니다.
 * SetInputValue : 입력데이터 세팅
 * Request  :  Non-blocking 요청
 * BlockRequest : Blocking 요청
 * Subscribe : 실시간 수신 요청
 * SubscribeLastest : 실시간 수신 (스냅샷성 데이터) 요청
 * Unsubscribe : 실시간 수신 해지
 * GetHeaderValue : 수신 헤더데이터
 * GetDataValue : 수신 데이터

```go
	stkmst.SetInputValue(0, "A000270")
	stkmst.BlockRequest()
	fmt.Println(stkmst.GetHeaderValue(1).Value())
```

### 이벤트처리
Received 이벤트는 Receiver 인터페이스를 구현하면 됩니다.
```go
type Receiver interface {
	Received(*CpClass)
}
```
이를 BindEvent 메서드를 사용하여 이벤트를 지정할수 있습니다.
이벤트를 해지하려면 UnbindEvent 메서드를 사용합니다.

```go
type RqTestStruct struct {
}

func (t *RqTestStruct) Received(c *CpClass) {
	fmt.Println(c.GetHeaderValue(1).Value())
	c.UnbindEvent() // 이벤트 바로 헤제
}
//... 생략

	evnt := &RqTestStruct{}
	cpobj.BindEvent(evnt)

```

PeekMessage 를 기반으로 동작하는 PumpWaitingMessages 함수를 제공합니다.


# 예제
