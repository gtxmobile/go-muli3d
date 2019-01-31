package common
func Clamp(v, minv, maxv float64) float64{
	if !(minv <= maxv) {
		panic("clamp assert")
	}
	if (v < minv) {
		return minv
	}
	if (v > maxv) {
		return maxv
	}
	return v
}
