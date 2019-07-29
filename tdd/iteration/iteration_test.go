package iteration

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {

	repeated := Repeat("a", 5)
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func ExampleRepeat() {
	str := Repeat("b", 5)
	fmt.Print(str)
	//Output: bbbbb
}
func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}
