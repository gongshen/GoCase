package cmap

import (
	"fmt"
	"testing"
)

func TestSegment_New(t *testing.T) {
	s := newSegment(-1, nil)
	if s == nil {
		t.Fatal("Couldn't new segment!")
	}
}

func TestSegment_Put(t *testing.T) {
	number := 10
	testCases := genTestPairs(number)
	s := newSegment(-1, nil)
	var count uint64
	for _, p := range testCases {
		ok, err := s.Put(p)
		if err != nil {
			t.Fatalf("An error occurs when put a pair to segment:%s!(key=%#v)", err, p.Key())
		}
		if !ok {
			t.Fatalf("Couldn't put a pair to segment!(key=%#v)", p.Key())
		}
		actualPair := s.Get(p.Key())
		if actualPair == nil {
			t.Fatalf("Couldn't get a pair from segment!(expected=%#v,actual=%#v)", p.Key(), nil)
		}
		ok, err = s.Put(p)
		if err != nil {
			t.Fatalf("An error occurs when repeat put a pair to segment:%s!(key=%#v)", err, p.Key())
		}
		if ok {
			t.Fatalf("Couldn't repeat put a pair to segment!(key=%#v)", p.Key())
		}
		count++
		if s.Size() != count {
			t.Fatalf("Inconsistent size!(expected=%#v,actual=%#v)",
				count, s.Size())
		}
	}
	if s.Size() != uint64(number) {
		t.Fatalf("Inconsistent size!(expected=%#v,actual=%#v)", number, s.Size())
	}
}

func TestSegment_PutInParallel(t *testing.T) {
	number := 10
	testCases := genTestPairs(number)
	s := newSegment(-1, nil)
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := s.Put(p)
			if err != nil {
				t.Fatalf("An error occurs when put a pair to segment:%s!(key=%#v)", err, p.Key())
			}
			if !ok {
				t.Fatalf("Couldn't put a pair to segment!(key=%#v)", p.Key())
			}
			actualPair := s.Get(p.Key())
			if actualPair == nil {
				t.Fatalf("Couldn't get a pair from segment!(expected=%#v,actual=%#v)", p.Key(), nil)
			}
			ok, err = s.Put(p)
			if err != nil {
				t.Fatalf("An error occurs when repeat put a pair to segment:%s!(key=%#v)", err, p.Key())
			}
			if ok {
				t.Fatalf("Couldn't repeat put a pair to segment!(key=%#v)", p.Key())
			}
		}
	}
	t.Run("Put in parallel", func(t *testing.T) {
		for _, p := range testCases {
			t.Run(fmt.Sprintf("Put:key=%#v", p.Key()), testFunc(p, t))
		}
	})
	if s.Size() != uint64(number) {
		t.Fatalf("Inconsistent size!(expected=%#v,actual=%#v)", number, s.Size())
	}
}

func TestSegment_GetInParallel(t *testing.T) {
	number := 10
	testCases := genTestPairs(number)
	s := newSegment(-1, nil)
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			actualPair := s.Get(p.Key())
			if actualPair == nil {
				t.Fatalf("Couldn't get a pair from segment!(expected=%#v,actual=%#v)", p.Key(), nil)
			}
			if actualPair.Key() != p.Key() {
				t.Fatalf("Inconsistent key!(expected=%#v,actual=%#v)", p.Key(), actualPair.Key())
			}
			if actualPair.Hash() != p.Hash() {
				t.Fatalf("Inconsistent hash!(expected=%d, actual=%d)",
					p.Hash(), actualPair.Hash())
			}
			if actualPair.Element() != p.Element() {
				t.Fatalf("Inconsistent element!(expected=%#v,actual=%#v)", p.Element(), actualPair.Element())
			}
		}
	}
	t.Run("Get in parallel", func(t *testing.T) {
		t.Run("Put in parallel", func(t *testing.T) {
			for _, p := range testCases {
				s.Put(p)
			}
		})
		for _, p := range testCases {
			t.Run(fmt.Sprintf("Get:key=%#v", p.Key()), testFunc(p, t))
		}
	})
	if s.Size() != uint64(number) {
		t.Fatalf("Inconsistent size!(expected=%#v,actual=%#v)", number, s.Size())
	}
}
func TestSegment_Delete(t *testing.T) {
	number := 10
	testCases := genTestPairs(number)
	s := newSegment(-1, nil)
	for _, p := range testCases {
		s.Put(p)
	}
	count := uint64(number)
	for _, p := range testCases {
		ok := s.Delete(p.Key())
		if !ok {
			t.Fatalf("Couldn't delete a pair from segment! (key=%#v)", p.Key())
		}
		actualPair := s.Get(p.Key())
		if actualPair != nil {
			t.Fatalf("Inconsistent pair!(expected=%#v, actual=%#v)",
				nil, actualPair)
		}
		ok = s.Delete(p.Key())
		if ok {
			t.Fatalf("Couldn't delete a pair from segment again! (key=%#v)", p.Key())
		}
		if count > 0 {
			count--
		}
		if s.Size() != count {
			t.Fatalf("Inconsistent size!(expected=%d, actual=%d)",
				count, s.Size())
		}
	}
	if s.Size() != 0 {
		t.Fatalf("Inconsistent size!(expected=%#v,actual=%#v)", 0, s.Size())
	}
}

func TestSegment_DeleteInParallel(t *testing.T) {
	number := 10
	testCases := genNoRepeatTestPairs(number)
	s := newSegment(-1, nil)
	for _, p := range testCases {
		s.Put(p)
	}
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok := s.Delete(p.Key())
			if !ok {
				t.Fatalf("Couldn't delete a pair from segment! (key=%#v)", p.Key())
			}
			actualPair := s.Get(p.Key())
			if actualPair != nil {
				t.Fatalf("Inconsistent pair!(expected=%#v, actual=%#v)",
					nil, actualPair)
			}
			ok = s.Delete(p.Key())
			if ok {
				t.Fatalf("Couldn't delete a pair from segment again! (key=%#v)", p.Key())
			}
		}
	}
	t.Run("Delete in parallel", func(t *testing.T) {
		for _, p := range testCases {
			t.Run(fmt.Sprintf("Delete:key=%#v", p.Key()), testFunc(p, t))
		}
	})
	if s.Size() != 0 {
		t.Fatalf("Inconsistent size!(expected=%#v,actual=%#v)", 0, s.Size())
	}
}

var testCaseNumberForSegment = 200000
var testCasesForSegment = genNoRepeatTestPairs(testCaseNumberForSegment)
var testCases1ForSegment = testCasesForSegment[:testCaseNumberForSegment/2]
var testCases2ForSegment = testCasesForSegment[testCaseNumberForSegment/2:]

func TestSegment_AllInParallel(t *testing.T) {
	testCases1 := testCases1ForBucketTest
	testCases2 := testCases2ForBucketTest
	s := newSegment(-1, nil)
	t.Run("All in parallel", func(t *testing.T) {
		t.Run("Put1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				_, err := s.Put(p)
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the segment:%s (key=%#v)",
						err, p.Key())
				}
			}
		})
		t.Run("Put2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				_, err := s.Put(p)
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the segment:%s (key=%#v)",
						err, p.Key())
				}
			}
		})
		t.Run("Get1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				actualPair := s.Get(p.Key())
				if actualPair == nil {
					continue
				}
				if actualPair.Key() != p.Key() {
					t.Fatalf("Inconsistent key!(expected=%#v,actual=%#v)", p.Key(), actualPair.Key())
				}
				if actualPair.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash!(expected=%d, actual=%d)",
						p.Hash(), actualPair.Hash())
				}
				if actualPair.Element() != p.Element() {
					t.Fatalf("Inconsistent element!(expected=%#v,actual=%#v)", p.Element(), actualPair.Element())
				}
			}
		})
		t.Run("Get2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				actualPair := s.Get(p.Key())
				if actualPair == nil {
					continue
				}
				if actualPair.Key() != p.Key() {
					t.Fatalf("Inconsistent key!(expected=%#v,actual=%#v)", p.Key(), actualPair.Key())
				}
				if actualPair.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash!(expected=%d, actual=%d)",
						p.Hash(), actualPair.Hash())
				}
				if actualPair.Element() != p.Element() {
					t.Fatalf("Inconsistent element!(expected=%#v,actual=%#v)", p.Element(), actualPair.Element())
				}
			}
		})
		t.Run("Delete1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				s.Delete(p.Key())
			}
		})
		t.Run("Delete2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				s.Delete(p.Key())
			}
		})
		t.Log(s.Size())
	})
}
