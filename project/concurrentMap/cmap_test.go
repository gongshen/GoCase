package cmap

import (
	"fmt"
	"testing"
)

func TestConcurrentMap_New(t *testing.T) {
	var concurrency int
	var pairRedistributor PairRedistributor
	cmap, err := NewConcurrenctMap(concurrency, pairRedistributor)
	if err == nil {
		t.Fatalf("No error occurs when new concurrent map with concurrency:%d,shouldn't be the case", concurrency)
	}
	concurrency = MAX_CONCURRENCY + 1
	cmap, err = NewConcurrenctMap(concurrency, pairRedistributor)
	if err == nil {
		t.Fatalf("No error occurs when new concurrent map with concurrency:%d,shouldn't be the case", concurrency)
	}
	concurrency = 16
	cmap, err = NewConcurrenctMap(concurrency, pairRedistributor)
	if err != nil {
		t.Fatalf("An error occurs when new a concurrent map: %s (concurrency=%d, pairRedistributor=%#v)",
			err, concurrency, pairRedistributor)
	}
	if cmap == nil {
		t.Fatalf("Couldn't a new concurrent map! (concurrency= %d, pairRedistributor=%#v)",
			concurrency, pairRedistributor)
	}
	if cmap.Concurrency() != concurrency {
		t.Fatalf("Inconsistent concurrency: expected= %d, actual= %d",
			concurrency, cmap.Concurrency())
	}
}
func TestConcurrentMap_Put(t *testing.T) {
	number := 30
	TestCases := genTestPairs(number)
	concurrency := 10
	var pairRedistributor PairRedistributor
	cmap, _ := NewConcurrenctMap(concurrency, pairRedistributor)
	var count uint64
	for _, p := range TestCases {
		key := p.Key()
		element := p.Element()
		ok, err := cmap.Put(key, element)
		if err != nil {
			t.Fatalf("An error occurs when put a key-element in concurrent map:%s!(key=%s,element=%#v)", err, key, element)
		}
		if !ok {
			t.Fatalf("Couldn't put key element to the concurrent map!(key=%s,element=%#v)", key, element)
		}
		actualElement := cmap.Get(key)
		if actualElement == nil {
			t.Fatalf("Inconsistent element!(expected= %#v, actual= %#v)",
				element, nil)
		}
		ok, err = cmap.Put(key, element)
		if err != nil {
			t.Fatalf("An error occurs when putting a repeated key element to the cmap:%s (key= %s, element= %#v)",
				err, key, element)
		}
		if ok {
			t.Fatalf("Couldn't put key element to the cmap! (key= %s, element= %#v)",
				key, element)
		}
		count++
		if cmap.Len() != count {
			t.Fatalf("Inconsistent size!(expected= %d, actual= %d)",
				count, cmap.Len())
		}
	}
}
func TestConcurrentMap_PutInParallel(t *testing.T) {
	number := 30
	TestCases := genNoRepeatTestPairs(number)
	concurrency := number / 2
	cmap, _ := NewConcurrenctMap(concurrency, nil)
	testFunc := func(key string, element interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := cmap.Put(key, element)
			if err != nil {
				t.Fatalf("An error occurs when put a key-element in concurrent map:%s!(key=%s,element=%#v)", err, key, element)
			}
			if !ok {
				t.Fatalf("Couldn't put key element to the concurrent map!(key=%s,element=%#v)", key, element)
			}
			actualElement := cmap.Get(key)
			if actualElement == nil {
				t.Fatalf("Inconsistent element!(expected= %#v, actual= %#v)",
					element, nil)
			}
			ok, err = cmap.Put(key, element)
			if err != nil {
				t.Fatalf("An error occurs when putting a repeated key element to the cmap:%s (key= %s, element= %#v)",
					err, key, element)
			}
			if ok {
				t.Fatalf("Couldn't put key element to the cmap! (key= %s, element= %#v)",
					key, element)
			}
		}
	}
	t.Run("Put in parallel", func(t *testing.T) {
		for _, p := range TestCases {
			key := p.Key()
			element := p.Element()
			t.Run(fmt.Sprintf("Put:(key=%s,element=%#v)", key, element), testFunc(key, element, t))
		}
	})
	if cmap.Len() != uint64(number) {
		t.Fatalf("Inconsistent size!(expected= %d, actual= %d)",
			number, cmap.Len())
	}
}
func TestConcurrentMap_GetInParallel(t *testing.T) {
	number := 30
	concurrency := number / 2
	testCases := genNoRepeatTestPairs(number)
	cmap, _ := NewConcurrenctMap(concurrency, nil)
	testFunc := func(key string, element interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			actualElement := cmap.Get(key)
			if actualElement == nil {
				t.Fatalf("Inconsistent element!(expected= %#v, actual= %#v)",
					element, nil)
			}
			if actualElement != element {
				t.Fatalf("Inconsistent element!(expected= %#v, actual= %#v)",
					element, actualElement)
			}
		}
	}
	t.Run("Get in parallel", func(t *testing.T) {
		t.Run("Put in parallel", func(t *testing.T) {
			for _, p := range testCases {
				cmap.Put(p.Key(), p.Element())
			}
		})
		for _, p := range testCases {
			t.Run(fmt.Sprintf("Get:(key=%s,element=%#v)", p.Key(), p.Element()), testFunc(p.Key(), p.Element(), t))
		}
	})
	if cmap.Len() != uint64(number) {
		t.Fatalf("Inconsistent size!(expected= %d, actual= %d)",
			number, cmap.Len())
	}
}

func TestConcurrentMap_Delete(t *testing.T) {
	number := 30
	concurrency := number / 2
	testCases := genNoRepeatTestPairs(number)
	cmap, _ := NewConcurrenctMap(concurrency, nil)
	for _, p := range testCases {
		cmap.Put(p.Key(), p.Element())
	}
	count := uint64(number)
	for _, p := range testCases {
		ok := cmap.Delete(p.Key())
		if !ok {
			t.Fatalf("Couldn't delete a key-element from cmap! (key= %s, element= %#v)",
				p.Key(), p.Element())
		}
		actualElement := cmap.Get(p.Key())
		if actualElement != nil {
			t.Fatalf("Inconsistent key-element!(expected= %#v, actual= %#v)",
				nil, actualElement)
		}
		ok = cmap.Delete(p.Key())
		if ok {
			t.Fatalf("Couldn't delete a key-element from cmap again! (key= %s, element= %#v)",
				p.Key(), p.Element())
		}
		if count > 0 {
			count--
		}
		if count != cmap.Len() {
			t.Fatalf("Inconsistent size!(expected= %d, actual= %d)",
				count, cmap.Len())
		}
	}
	if cmap.Len() != 0 {
		t.Fatalf("Inconsistent size!(expected= %d, actual= %d)",
			0, cmap.Len())
	}
}
func TestConcurrentMap_DeleteInParallel(t *testing.T) {
	number := 30
	concurrency := number / 2
	testCases := genNoRepeatTestPairs(number)
	cmap, _ := NewConcurrenctMap(concurrency, nil)
	testFunc := func(key string, element interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			ok := cmap.Delete(key)
			if !ok {
				t.Fatalf("Couldn't delete a key-element from cmap! (key= %s, element= %#v)",
					key, element)
			}
			actualElement := cmap.Get(key)
			if actualElement != nil {
				t.Fatalf("Inconsistent key-element!(expected= %#v, actual= %#v)",
					nil, actualElement)
			}
			ok = cmap.Delete(key)
			if ok {
				t.Fatalf("Couldn't delete a key-element from cmap again! (key= %s, element= %#v)",
					key, element)
			}
		}
	}
	t.Run("Delete in parallel", func(t *testing.T) {
		t.Run("Put in parallel", func(t *testing.T) {
			for _, p := range testCases {
				cmap.Put(p.Key(), p.Element())
			}
		})
		for _, p := range testCases {
			t.Run(fmt.Sprintf("Delete:(key=%s,element=%#v)", p.Key(), p.Element()), testFunc(p.Key(), p.Element(), t))
		}
	})
	if cmap.Len() != 0 {
		t.Fatalf("Inconsistent size!(expected= %d, actual= %d)",
			0, cmap.Len())
	}
}

var testCaseNumberForCmapTest = 200000
var testCasesForCmapTest = genNoRepeatTestPairs(testCaseNumberForCmapTest)
var testCases1ForCmapTest = testCasesForCmapTest[:testCaseNumberForCmapTest/2]
var testCases2ForCmapTest = testCasesForCmapTest[testCaseNumberForCmapTest/2:]

func TestConcurrentMap_AllInParallel(t *testing.T) {
	testCases1 := testCases1ForCmapTest
	testCases2 := testCases2ForCmapTest
	concurrency := testCaseNumberForCmapTest / 4
	cmap, _ := NewConcurrenctMap(concurrency, nil)
	t.Run("All in parallel", func(t *testing.T) {
		t.Run("Put1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				_, err := cmap.Put(p.Key(), p.Element())
				if err != nil {
					t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key= %s, element= %#v)",
						err, p.Key(), p.Element())
				}
			}
		})
		t.Run("Put2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				_, err := cmap.Put(p.Key(), p.Element())
				if err != nil {
					t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key= %s, element= %#v)",
						err, p.Key(), p.Element())
				}
			}
		})
		t.Run("Get1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				actualElement := cmap.Get(p.Key())
				if actualElement == nil {
					continue
				}
				if actualElement != p.Element() {
					t.Fatalf("Inconsistent element!(expected=%#v,actual=%#v)", p.Element(), actualElement)
				}
			}
		})
		t.Run("Get2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				actualElement := cmap.Get(p.Key())
				if actualElement == nil {
					continue
				}
				if actualElement != p.Element() {
					t.Fatalf("Inconsistent element!(expected=%#v,actual=%#v)", p.Element(), actualElement)
				}
			}
		})
		t.Run("Delete1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				cmap.Delete(p.Key())
			}
		})
		t.Run("Delete2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				cmap.Delete(p.Key())
			}
		})
		//这里执行删除操作，最后是删不光的，因为前期的一些删除操作都失败了（数据未放入）
	})

}
