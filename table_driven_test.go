package main

import (
	"fmt"
	"sync"
	"testing"
)

type RequestLimitExceeded struct {
	limit int
	Err   error
}

func (e *RequestLimitExceeded) Error() string {
	return fmt.Sprintf("Limit: %d was exceed", e.limit)
}

type StubType struct {
	m          sync.Mutex
	response   string
	reqCounter int
}

const RequestLimit = 100

func (st *StubType) Read(p []byte) (n int, err error) {
	if st.reqCounter == RequestLimit {
		err = &RequestLimitExceeded{limit: RequestLimit}
		return
	}
	st.m.Lock()
	defer st.m.Unlock()
	copy(p, st.response)
	st.reqCounter++
	return st.reqCounter, nil
}

func TestMultipleReader(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"less than request limit": TestLessThanRequestLimit,
		"equal to request limit":  TestEqualToRequestLimit,
		"more than request limit": TestEqualToRequestLimit,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

func testMultipleReader(t *testing.T, tt testTpe) {
	tt.wg.Add(tt.currReqPool)
	for i := 0; i < tt.currReqPool; i++ {
		go func(i int) {
			var resp []byte
			if _, err := tt.st.Read(resp); err != nil {
				assertNotEqualRequestLimit(t, tt.st, RequestLimit, tt.wg)
			} else {
				tt.wg.Done()
			}
		}(i)
	}
	tt.wg.Wait()
}

type testTpe struct {
	currReqPool int
	wg          *sync.WaitGroup
	st          *StubType
}

func TestLessThanRequestLimit(t *testing.T) {
	tt := testTpe{
		currReqPool: RequestLimit - 50,
		wg:          &sync.WaitGroup{},
		st:          &StubType{reqCounter: RequestLimit, response: "my-dummy-response"},
	}
	testMultipleReader(t, tt)
}

func TestEqualToRequestLimit(t *testing.T) {
	tt := testTpe{
		currReqPool: RequestLimit,
		wg:          &sync.WaitGroup{},
		st:          &StubType{reqCounter: RequestLimit, response: "my-dummy-response"},
	}
	testMultipleReader(t, tt)
}

func TestMoreThanRequestLimit(t *testing.T) {
	tt := testTpe{
		currReqPool: RequestLimit + 50,
		wg:          &sync.WaitGroup{},
		st:          &StubType{reqCounter: RequestLimit, response: "my-dummy-response"},
	}
	testMultipleReader(t, tt)
}

func assertNotEqualRequestLimit(t *testing.T, got *StubType, want int, wg *sync.WaitGroup) {
	t.Helper()
	if got.reqCounter != want {
		t.Errorf("got %q, want %q", got.reqCounter, want)
	} else {
		wg.Done()
	}
}
