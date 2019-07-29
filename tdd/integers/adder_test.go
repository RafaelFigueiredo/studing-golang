package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {

	resultCheck := func(t *testing.T, sum, expected int) {
		t.Helper()
		if sum != expected {
			t.Errorf("Adder testing error, expected: '%d' , got: '%d'", expected, sum)
		}
	}

	t.Run("Two plus two", func(t *testing.T) {
		sum := Add(2, 2)
		expected := 4
		resultCheck(t, sum, expected)
	})
	t.Run("Zero plus zero", func(t *testing.T) {
		sum := Add(0, 0)
		expected := 0
		resultCheck(t, sum, expected)
	})
	t.Run("Minus two plus two", func(t *testing.T) {
		sum := Add(-2, 2)
		expected := 0
		resultCheck(t, sum, expected)
	})

}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
