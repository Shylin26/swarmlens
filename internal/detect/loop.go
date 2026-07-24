package detect

func LoopDetected(window []string, minRepeats int) bool {
	requiredLength := minRepeats * 2
	if len(window) < requiredLength {
		return false
	}
	a := window[0]
	b := window[1]
	for i := 0; i < requiredLength; i += 2 {
		if window[i] != a || window[i+1] != b {
			return false
		}
	}
	return true

}
