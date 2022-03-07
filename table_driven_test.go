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
	if st.reqCounter == st.reqCounter {
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
	type testCases struct {
		description string
		currReqPool int
		wg          *sync.WaitGroup
		st          *StubType
	}

	for _, s := range []*testCases{
		{
			description: "request number is less than the request limit of type",
			currReqPool: RequestLimit - 50,
			wg:          &sync.WaitGroup{},
			st:          &StubType{reqCounter: RequestLimit, response: "my-dummy-response"},
		},
		{
			description: "request number is equal the request limit of type",
			currReqPool: RequestLimit,
			wg:          &sync.WaitGroup{},
			st:          &StubType{reqCounter: RequestLimit, response: "my-dummy-response"},
		},
		{
			description: "request number is more than the request limit of type",
			currReqPool: RequestLimit + 50,
			wg:          &sync.WaitGroup{},
			st:          &StubType{reqCounter: RequestLimit, response: "my-dummy-response"},
		},
	} {
		t.Run("concurrent access to reader", func(t *testing.T) {
			s.wg.Add(s.currReqPool)
			for i := 0; i < s.currReqPool; i++ {
				go func(i int) {
					var resp []byte
					if _, err := s.st.Read(resp); err != nil {
						assertNotEqualRequestLimit(t, s.st, RequestLimit, s.wg)
					} else {
						s.wg.Done()
					}
				}(i)
			}
			s.wg.Wait()
		})
	}
}

func assertNotEqualRequestLimit(t *testing.T, got *StubType, want int, wg *sync.WaitGroup) {
	t.Helper()
	if got.reqCounter != want {
		t.Errorf("got %q, want %q", got.reqCounter, want)
	} else {
		wg.Done()
	}
}
