package cmap

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
)

type keyElement struct {
	key     string
	element interface{}
}

func TestNewPair(t *testing.T) {
	testCases := genKeyElementSlice(100)
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("key=%s,element=%#v", testCase.key, testCase.element), func(t *testing.T) {
			pair, err := newPair(testCase.key, testCase.element)
			if err != nil {
				t.Fatalf("An error occurs when new a pair:%s (key:%s,element:%#v)", err, testCase.key, testCase.element)
			}
			if pair == nil {
				t.Fatalf("Could new a pair!(key:%s,element:%#v)", testCase.key, testCase.element)
			}
		})
	}
}

func TestPair_Hash_Element_Key(t *testing.T) {
	testCases := genKeyElementSlice(30)
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("key:%s,element:%#v", testCase.key, testCase.element), func(t *testing.T) {
			pair, err := newPair(testCase.key, testCase.element)
			if err != nil {
				t.Fatalf("An error occurs when new a pair:%s (key:%s,element:%#v)", err, testCase.key, testCase.element)
			}
			exceptedHash := hash(testCase.key)
			if pair.Hash() != exceptedHash {
				t.Fatalf("Inconsistent hash:except:%d,actual:%d", hash(testCase.key), exceptedHash)
			}
			if pair.Key() != testCase.key {
				t.Fatalf("Inconsistent key:except:%s,actual:%#v", testCase.key, pair.Key())
			}
			if pair.Element() != testCase.element {
				t.Fatalf("Inconsistent hash:except:%v,actual:%#v", testCase.element, pair.Element())
			}
		})
	}
}

func TestPair_SetElement(t *testing.T) {
	testCases := genKeyElementSlice(30)
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("key:%s,element:%#v", testCase.key, testCase.element), func(t *testing.T) {
			pair, err := newPair(testCase.key, testCase.element)
			if err != nil {
				t.Fatalf("An error occurs when new a pair:%s (key:%s,element:%#v)", err, testCase.key, testCase.element)
			}
			element := randString()
			err = pair.SetElement(element)
			if err != nil {
				t.Fatalf("An error when set element:%s", err)
			}
			if pair.Element() != element {
				t.Fatalf("Inconsistent set element: except:%#v,actual:%#v)", element, pair.Element())
			}
		})
	}
}

func TestPair_Copy(t *testing.T) {
	testCases := genKeyElementSlice(10)
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("key:%s,element:%#v", testCase.key, testCase.element), func(t *testing.T) {
			pair, err := newPair(testCase.key, testCase.element)
			if err != nil {
				t.Fatalf("An error occurs when new a pair:%s (key:%s,element:%#v)", err, testCase.key, testCase.element)
			}
			pCopy := pair.Copy()
			if pair.Hash() != pCopy.Hash() {
				t.Fatalf("Inconsistent hash:except:%d,actual:%d", pair.Hash(), pCopy.Hash())
			}
			if pair.Key() != pCopy.Key() {
				t.Fatalf("Inconsistent key:except:%s,actual:%s", pair.Key(), pCopy.Key())
			}
			if pCopy.Element() != pair.Element() {
				t.Fatalf("Inconsistent element:except:%#v,actual:%#v", pair.Element(), pCopy.Element())
			}
		})
	}
}

func TestPair_Next(t *testing.T) {
	number := 10
	testCases := genKeyElementSlice(number)
	var current Pair
	var prev Pair
	var err error
	for _, testCase := range testCases {
		current, err = newPair(testCase.key, testCase.element)
		t.Log(current.String())
		if err != nil {
			t.Fatalf("An error occurs when new a pair:%s (key:%s,element:%#v)", err, testCase.key, testCase.element)
		}
		if prev != nil {
			current.SetNext(prev)
		}
		prev = current
	}
	for i := number - 1; i >= 0; i-- {
		next := current.Next()
		if i == 0 {
			if next != nil {
				t.Fatalf("Next is not nil!(pair:%#v,index:%d)", current, i)
			}
		} else {
			if next == nil {
				t.Fatalf("Next is nil!(pair:%#v,index:%d)", current, i)
			}
			expectNext := testCases[i-1]
			if expectNext.key != next.Key() {
				t.Fatalf("Inconsistent key:except:%v,actual:%#v", expectNext.key, next.Key())
			}
			if expectNext.element != next.Element() {
				t.Fatalf("Inconsistent element:except:%v,actual:%#v", expectNext.element, next.Element())
			}
		}
		current = next
	}
}

func genKeyElementSlice(number int) []*keyElement {
	testCases := make([]*keyElement, number)
	for i := 0; i < number; i++ {
		testCases[i] = &keyElement{randString(), randElement()}
	}
	return testCases
}

func randElement() interface{} {
	if i := rand.Int31(); i%3 != 0 {
		return i
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}

func randString() string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}
