//并发安全字典与普通字典性能比较
package cmap_test

import (
	cmap "book/并发编程/同步/并发安全的字典"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

//!+   Put

func BenchmarkConcurrentMap_PutVariable(b *testing.B) {
	var number = 20
	testCases := cmap.GenNoRepeatTestPairs(number)
	concurrency := number / 4
	cmap, _ := cmap.NewConcurrenctMap(concurrency, nil)
	b.ResetTimer()
	for _, p := range testCases {
		key := p.Key()
		element := p.Element()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cmap.Put(key, element)
			}
		})
	}
}

func BenchmarkConcurrentMap_PutInvariable(b *testing.B) {
	number := 20
	concurrency := number / 4
	cmap, _ := cmap.NewConcurrenctMap(concurrency, nil)
	key := "invariable key"
	b.ResetTimer()
	for i := 0; i < number; i++ {
		element := strconv.Itoa(i)
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cmap.Put(key, element)
			}
		})
	}
}
func BenchmarkMap_Put(b *testing.B) {
	number := 10
	testCases := cmap.GenNoRepeatTestPairs(number)
	m := make(map[string]interface{})
	b.ResetTimer()
	for _, p := range testCases {
		key := p.Key()
		element := p.Element()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m[key] = element
			}
		})
	}
}

//!+  Get
func BenchmarkConcurrentMap_Get(b *testing.B) {
	number := 100000
	testCases := cmap.GenNoRepeatTestPairs(number)
	concurrency := number / 4
	cmap, _ := cmap.NewConcurrenctMap(concurrency, nil)
	for _, p := range testCases {
		cmap.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cmap.Get(key)
			}
		})
	}
}
func BenchmarkMap_Get(b *testing.B) {
	number := 100000
	m := make(map[string]interface{})
	testCases := cmap.GenNoRepeatTestPairs(number)
	for _, p := range testCases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_ = m[key]
			}
		})
	}
}

//!+   Delete
func BenchmarkConcurrentMap_Delete(b *testing.B) {
	number := 100000
	concurrency := number / 4
	testCases := cmap.GenNoRepeatTestPairs(number)
	cmap, _ := cmap.NewConcurrenctMap(concurrency, nil)
	for _, p := range testCases {
		cmap.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cmap.Delete(key)
			}
		})
	}
}
func BenchmarkMap_Delete(b *testing.B) {
	number := 100000
	testCases := cmap.GenNoRepeatTestPairs(number)
	m := make(map[string]interface{})
	for _, p := range testCases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				delete(m, key)
			}
		})
	}
}

//!+   Len

func BenchmarkConcurrentMap_Len(b *testing.B) {
	number := 100000
	concurrency := number / 4
	testCases := cmap.GenNoRepeatTestPairs(number)
	cmap, _ := cmap.NewConcurrenctMap(concurrency, nil)
	for _, p := range testCases {
		cmap.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("Len(%d)", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cmap.Len()
			}
		})
	}
}

func BenchmarkMap_Len(b *testing.B) {
	number := 100000
	testCases := cmap.GenNoRepeatTestPairs(number)
	m := make(map[string]interface{})
	for _, p := range testCases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("Len(%d)", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_ = len(m)
			}
		})
	}
}
