package test


func TestLess (t *testing.T) {
	res := []int{1, 2, 3, 4, 5, 6, 7}
	ex := []int{1, 2, 3 , 4, 5}
	if res != ex {
		t.Errorf("Expected %v, but res %v", ex, res)
	}

	}
}
func BenchmarkXxx(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}