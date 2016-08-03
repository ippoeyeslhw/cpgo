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
 * 사이보스플러스 [도움말(비공식)](http://cybosplus.github.io/)

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
반드시 cointialize, uncoinitialize 호출
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

### Property
Property값을 가져오려면 getter를 사용하여야 합니다.
도움말의 Property명에 앞에 Get을 붙인 메서드명을 사용합니다.
```go
	tmp := &CpClass{}
	tmp.Create("CpUtil.CpCybos")
	defer tmp.Release()
	fmt.Println(tmp)

	fmt.Println("isconnect: ", tmp.GetIsConnect().Value())
	fmt.Println("servertype: ", tmp.GetServerType().Value())
	fmt.Println("remain time: ", tmp.GetLimitRequestRemainTime().Value())
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
이를 사용하여 Received이벤트를 수신할때까지 대기할수 있습니다.
```go
for  {
	PumpWaitingMessages()
	time.Sleep(1)
}
```

# 예제

통신방식을 기준으로 작성한 몇가지 예제입니다.

### BlockRequest사용
```go
package main

import (
	"fmt"

	ole "github.com/go-ole/go-ole"
	"github.com/ippoeyeslhw/cpgo"
)

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	stkmst := &cpgo.CpClass{}
	stkmst.Create("CpDib.StockMst")
	defer stkmst.Release()

	stkmst.SetInputValue(0, "A000270")
	stkmst.BlockRequest()
	fmt.Println(stkmst.GetHeaderValue(1).Value())
}
// 결과:
// 기아차
```

### Request사용
```go
package main

import (
	"fmt"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/ippoeyeslhw/cpgo"
)

type ContEvnt struct {
	isDone bool
}

func (s *ContEvnt) Received(c *cpgo.CpClass) {
	count := c.GetHeaderValue(2).Value().(int16) // 수신개수
	fmt.Println("response count: ", count)
	for i := 0; i < int(count); i++ {
		fmt.Println(
			c.GetDataValue(1, i).Value(), // 종목코드
			c.GetDataValue(4, i).Value()) // 내용
	}
	if c.GetContinue().Value() == int32(1) { // 연속데이터 있음
		fmt.Println("next request")
		c.Request() // 재요청
	} else {
		s.isDone = true
	}
}

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	cpmw := &cpgo.CpClass{}
	cpmw.Create("CpSysDib.CpMarketWatch")
	defer cpmw.Release()

	evnt := &ContEvnt{false}
	cpmw.BindEvent(evnt)

	// 연속조회  request 갯수제한 유의
	cpmw.SetInputValue(0, "*") // 전종목
	cpmw.SetInputValue(1, "2") // 공시정보
	cpmw.Request()

	for evnt.isDone == false {
		cpgo.PumpWaitingMessages()
		time.Sleep(1)
	}
	cpmw.UnbindEvent()
}

// 결과:
//response count:  20
//A151910 나노스(주) 주권매매거래정지기간변경(개선기간 부여)
//A043290 케이맥(주) 단일판매ㆍ공급계약체결
//A151910 나노스(주) 기타시장안내(기업심사위원회 심의결과 및 개선기간 부여 안내)
//A900050 중국원양자원유한공사 기타 경영사항(자율공시)(자회사의 주요경영사항)
//A043710 (주)서울리거 기타시장안내(상장적격성 실질심사사유 추가 관련 )
//A065420 (주)에스아이리소스 주권매매거래정지(불성실공시법인 지정)
//A065420 (주)에스아이리소스 불성실공시법인지정(공시불이행)
//A011200 현대상선(주) 전환가액ㆍ신주인수권행사가액ㆍ교환가액의 조정(안내공시)
//A011200 현대상선(주) 전환가액ㆍ신주인수권행사가액ㆍ교환가액의 조정(안내공시)
// ...
// ...
```

### Subscribe 예제

```go
package main

import (
	"fmt"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/ippoeyeslhw/cpgo"
)

type SubTestStruct struct {
	cont bool
	cnt  int
}

func (s *SubTestStruct) Received(c *cpgo.CpClass) {

	fmt.Printf("(%f)%d , (%f)%d\n",
		c.GetHeaderValue(14).Value(), // 1차 매수호가
		c.GetHeaderValue(15).Value(), // 1차 매수잔량
		c.GetHeaderValue(25).Value(), // 1차 매도호가
		c.GetHeaderValue(26).Value()) // 1차 매도잔량

	if s.cnt > 100 {
		// 100건이 넘을시 중단
		s.cont = false
	}
	s.cnt++
}

func main() {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	tmp := &cpgo.CpClass{}
	tmp.Create("CpSysDib.CmeCurr")
	defer tmp.Release()
	fmt.Println(tmp)

	evnt := &SubTestStruct{true, 0}
	tmp.BindEvent(evnt)
	fmt.Println(tmp)

	// 야간 CME 선물시장 밤에 테스트할것
	tmp.SetInputValue(0, "101L9")
	tmp.Subscribe()

	fmt.Println("sub/pub start")

	for evnt.cont == true {
		cpgo.PumpWaitingMessages()
		time.Sleep(1)
	}
	tmp.Unsubscribe()
	tmp.UnbindEvent()
}

// 결과:
//sub/pub start
//(251.449997)30 , (251.500000)2
//(251.449997)35 , (251.500000)2
//(251.449997)35 , (251.550003)34
//(251.449997)35 , (251.550003)34
//(251.449997)35 , (251.550003)34
//(251.449997)35 , (251.500000)1
// ...
// ...


```