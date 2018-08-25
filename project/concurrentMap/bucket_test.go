package cmap

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	b := newBucket()
	if b == nil {
		t.Fatalf("New bucket is fail")
	}
}

//1. Put和Get是否是同一个Pair
//2. 重复Put同一个Pair是否有问题
//3. bucket的Size是否和放置的Pair数量一致
func TestBucket_Put_Get_Size(t *testing.T) {
	number := 30
	testCases := genTestPairs(number)
	b := newBucket()
	var count uint64
	for _, p := range testCases {
		ok, err := b.Put(p, nil)
		if err != nil {
			t.Fatalf("An error occurs when putting the pair to the bucket:%s,(pair:%#v)", err, p)
		}
		if !ok {
			t.Fatalf("Cannot put the pair to the bucket(pair:%#v)", p)
		}
		actualPair := b.Get(p.Key())
		if actualPair == nil {
			t.Fatalf("Inconsistent pair:expected:%#v,actual:%#v", p.Element(), nil)
		}
		ok, err = b.Put(p, nil)
		if err != nil {
			t.Fatalf("An error occurs when putting the repeat pair to the bucket:%s,(pair:%#v)", err, p)
		}
		if ok {
			t.Fatalf("Cannot put the repeat pair to the bucket(pair:%#v)", p)
		}
		count++
		if b.Size() != count {
			t.Fatalf("Inconsistent size:expected:%d,actual:%d", count, b.Size())
		}
	}
	if b.Size() != uint64(number) {
		t.Fatalf("Inconsistent size:expected:%d,actual:%d", number, b.Size())
	}
}

func TestBucket_PutInParallel(t *testing.T) {
	number := 5
	testCases := genNoRepeatTestPairs(number)
	b := newBucket()
	lock := new(sync.Mutex)
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := b.Put(p, lock)
			if err != nil {
				t.Fatalf("An error occurs when putting the pair to the bucket:%s,(pair:%#v)", err, p)
			}
			if !ok {
				t.Fatalf("Cannot put the pair to the bucket(pair:%#v)", p)
			}
			actualPair := b.Get(p.Key())
			if actualPair == nil {
				t.Fatalf("Inconsistent pair:expected:%#v,actual:%#v", p.Element(), nil)
			}
			ok, err = b.Put(p, lock)
			if err != nil {
				t.Fatalf("An error occurs when putting the repeat pair to the bucket:%s,(pair:%#v)", err, p)
			}
			if ok {
				t.Fatalf("Cannot put the repeat pair to the bucket(pair:%#v)", p)
			}
		}
	}
	t.Run("In Parallel", func(t *testing.T) {
		for _, p := range testCases {
			t.Run(fmt.Sprintf("key=%s", p.Key()), testFunc(p, t))
		}
	})
	if b.Size() != uint64(number) {
		t.Fatalf("Inconsistent size:expected:%d,actual:%d", number, b.Size())
	}
}
func TestBucket_GetInParallel(t *testing.T) {
	number := 10
	testCases := genNoRepeatTestPairs(number)
	b := newBucket()
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			actualPair := b.Get(p.Key())
			if actualPair == nil {
				t.Fatalf("Not found pair in bucket!(key:%s)", p.Key())
			}
			if actualPair.Key() != p.Key() {
				t.Fatalf("Inconsistent key: expected: %s, actual: %s",
					p.Key(), actualPair.Key())
			}
			if actualPair.Hash() != p.Hash() {
				t.Fatalf("Inconsistent hash: expected: %d, actual: %d",
					p.Hash(), actualPair.Hash())
			}
			if actualPair.Element() != p.Element() {
				t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
					p.Element(), actualPair.Element())
			}
		}
	}
	t.Run("Get in parallel!", func(t *testing.T) {
		t.Run("Put in parallel", func(t *testing.T) {
			for _, p := range testCases {
				b.Put(p, nil)
			}
		})
		for _, p := range testCases {
			t.Run(fmt.Sprintf("Get:key=%s", p.Key()), testFunc(p, t))
		}
	})
	if b.Size() != uint64(number) {
		t.Fatalf("Inconsistent size:expected:%d,actual:%d", number, b.Size())
	}
}

func TestBucket_GetFirstPair(t *testing.T) {
	number := 10
	testCases := genNoRepeatTestPairs(number)
	b := newBucket()
	//版本一：将Pair全部放入Bucket中，然后循环遍历其FirstPair

	for _, p := range testCases {
		b.Put(p, nil)
	}
	size := b.Size()
	if size != uint64(number) {
		t.Fatalf("Inconsistent size!(expected:%d,actual%d)", number, size)
	}
	currentFirstPair := b.GetFirstPair()
	for i := int(size - 1); i >= 0; i-- {
		expectedPair := testCases[i]
		if currentFirstPair.Key() != expectedPair.Key() {
			t.Fatalf("Inconsistent key: expected: %s, actual: %s",
				expectedPair.Key(), currentFirstPair.Key())
		}
		if currentFirstPair.Element() != expectedPair.Element() {
			t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
				expectedPair.Element(), currentFirstPair.Element())
		}
		currentFirstPair = currentFirstPair.Next()
	}
	if currentFirstPair != nil {
		t.Fatal("The next of the last pair in bucket isn't nil")
	}
	//版本二：在将Pair循环放入Bucket中时，同时测试其FirstPair
	/*
		testFunc:= func(p Pair,t *testing.T)func(t *testing.T) {
			return func(t *testing.T) {
				t.Parallel()
				b.Put(p,nil)
				firstPair:=b.GetFirstPair()
				if firstPair.Key()!=p.Key(){
					t.Fatalf("Get the incorrect first pair!(expected:%#v,actual:%#v)",p,firstPair)
				}
			}
		}
		t.Run("GetFirstPair in parallel!", func(t *testing.T) {
			for _,p:=range testCases{
				t.Run(fmt.Sprintf("key:%#v,element:%#v",p.Key(),p.Element()),testFunc(p,t))
			}
		})
	*/
	if b.Size() != uint64(number) {
		t.Fatalf("Inconsistent size:expected:%d,actual:%d", number, b.Size())
	}
}
func TestBucket_Delete(t *testing.T) {
	number := 2
	testCases := genTestPairs(number)
	b := newBucket()
	for _, p := range testCases {
		b.Put(p, nil)
	}
	count := uint64(number)
	for _, p := range testCases {
		ok := b.Delete(p.Key(), nil)
		if !ok {
			t.Fatalf("Couldn't delete a pair from a bucket!(pair:%#v)", p)
		}
		actualPair := b.Get(p.Key())
		if actualPair != nil {
			t.Fatalf("Inconsistent pair!(expected:%#v,actualPair:%#v)", nil, actualPair)
		}
		ok = b.Delete(p.Key(), nil)
		if ok {
			t.Fatalf("An error occurs when delete the pair again!(pair:%#v)", p)
		}
		if count > 0 {
			count--
		}
		if b.Size() != count {
			t.Fatalf("Inconsistent size:expected:%d,actual:%d", count, b.Size())
		}
	}
	if b.Size() != 0 {
		t.Fatalf("Inconsistent size:expected:%d,actual:%d", 0, b.Size())
	}
}

func TestBucket_DeleteInParallel(t *testing.T) {
	number := 30
	testCases := genNoRepeatTestPairs(number)
	b := newBucket()
	for _, p := range testCases {
		b.Put(p, nil)
	}
	lock := new(sync.Mutex)
	testFunc := func(p Pair, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok := b.Delete(p.Key(), lock)
			if !ok {
				t.Fatalf("Couldn't delete a pair from a bucket!(pair:%#v)", p)
			}
			actualPair := b.Get(p.Key())
			if actualPair != nil {
				t.Fatalf("Inconsistent pair!(expected:%#v,actualPair:%#v)", nil, actualPair)
			}
			ok = b.Delete(p.Key(), nil)
			if ok {
				t.Fatalf("An error occurs when delete the pair again!(pair:%#v)", p)
			}
		}
	}
	t.Run("Delete in parallel", func(t *testing.T) {
		for _, p := range testCases {
			t.Run(fmt.Sprintf("key:%#v", p.Key()), testFunc(p, t))
		}
	})
	if b.Size() != 0 {
		t.Fatalf("Inconsistent size:expected:%d,actual:%d", 0, b.Size())
	}
}

func TestBucket_Clear(t *testing.T) {
	number := 10
	testCases := genTestPairs(number)
	b := newBucket()
	for _, p := range testCases {
		b.Put(p, nil)
	}
	b.Clear(nil)
	if b.Size() != 0 {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			0, b.Size())
	}
}

func TestBucket_ClearInParallel(t *testing.T) {
	number := 1000
	testCases := genTestPairs(number)
	b := newBucket()
	lock := new(sync.Mutex)
	t.Run("Clear in parallel!", func(t *testing.T) {
		t.Run("Put", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases {
				ok, err := b.Put(p, lock)
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the bucket: %s (pair: %#v)",
						err, p)
				}
				if !ok {
					t.Fatalf("Couldn't put pair to the bucket! (pair: %#v)",
						p)
				}
			}
		})
		t.Run("Clear", func(t *testing.T) {
			t.Parallel()
			for i := number; i > 0; i-- {
				b.Clear(lock)
			}
		})
	})
	if b.Size() > 0 {
		t.Log("Not Clean.Clear again!")
		b.Clear(nil)
	}
	if b.Size() != 0 {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			0, b.Size())
	}
}

var testCaseNumberForBucketTest = 200000
var testCasesForBucketTest = genNoRepeatTestPairs(testCaseNumberForBucketTest)
var testCases1ForBucketTest = testCasesForBucketTest[:testCaseNumberForBucketTest/2]
var testCases2ForBucketTest = testCasesForBucketTest[testCaseNumberForBucketTest/2:]

func TestBucket_AllInParallel(t *testing.T) {
	testCases1 := testCases1ForBucketTest
	testCases2 := testCases2ForBucketTest
	b := newBucket()
	lock := new(sync.Mutex)
	t.Run("All in parallel", func(t *testing.T) {
		t.Run("Put1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				existingPair := b.Get(p.Key())
				if existingPair != nil {
					b.Delete(p.Key(), lock)
				}
				ok, err := b.Put(p, lock)
				if !ok {
					t.Fatalf("Couldn't put a pair to the bucket! (pair: %#v)", p)
				}
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the bucket: %s (pair: %#v)",
						err, p)
				}
			}
		})
		t.Run("Put2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				existingPair := b.Get(p.Key())
				if existingPair != nil {
					b.Delete(p.Key(), lock)
				}
				ok, err := b.Put(p, lock)
				if !ok {
					t.Fatalf("Couldn't put a pair to the bucket! (pair: %#v)", p)
				}
				if err != nil {
					t.Fatalf("An error occurs when putting a pair to the bucket: %s (pair: %#v)",
						err, p)
				}
			}
		})
		t.Run("Get1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				actualPair := b.Get(p.Key())
				if actualPair == nil {
					continue
				}
				if actualPair.Key() != p.Key() {
					t.Fatalf("Inconsistent key: expected: %s, actual: %s",
						p.Key(), actualPair.Key())
				}
				if actualPair.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash: expected: %d, actual: %d",
						p.Hash(), actualPair.Hash())
				}
				if actualPair.Element() != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						p.Element(), actualPair.Element())
				}
			}
		})
		t.Run("Get2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				actualPair := b.Get(p.Key())
				if actualPair == nil {
					continue
				}
				if actualPair.Key() != p.Key() {
					t.Fatalf("Inconsistent key: expected: %s, actual: %s",
						p.Key(), actualPair.Key())
				}
				if actualPair.Hash() != p.Hash() {
					t.Fatalf("Inconsistent hash: expected: %d, actual: %d",
						p.Hash(), actualPair.Hash())
				}
				if actualPair.Element() != p.Element() {
					t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
						p.Element(), actualPair.Element())
				}
			}
		})
		t.Run("Delete1", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases1 {
				b.Delete(p.Key(), lock)
			}
		})
		t.Run("Delete2", func(t *testing.T) {
			t.Parallel()
			for _, p := range testCases2 {
				b.Delete(p.Key(), lock)
			}
		})
		t.Run("Clear", func(t *testing.T) {
			t.Parallel()
			go func() {
				for _ = range time.Tick(time.Millisecond*10) {
					b.Clear(lock)
				}
			}()
			<-time.Tick(10*time.Millisecond)
		})
	})
}
func genTestPairs(number int) []Pair {
	testCases := make([]Pair, number)
	for i := 0; i < number; i++ {
		testCases[i], _ = newPair(randString(), randElement())
	}
	return testCases
}

func genNoRepeatTestPairs(number int) []Pair {
	testCases := make([]Pair, number)
	m := make(map[string]struct{})
	for i := 0; i < number; i++ {
		for {
			p, _ := newPair(randString(), randElement())
			if _, ok := m[p.Key()]; !ok {
				testCases[i] = p
				m[p.Key()] = struct{}{}
				break
			}
		}
	}
	return testCases
}
