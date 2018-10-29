package stack

import "testing"

func BenchmarkSlicePush(b *testing.B) {
	stack := NewSliceInplement()
	b.ResetTimer()
	b.Run("push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Push(i)
		}
	})

}

func BenchmarkListPush(b *testing.B) {
	stack := NewHeader()
	b.ResetTimer()
	b.Run("push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stack.Push(i)
		}
	})

}
