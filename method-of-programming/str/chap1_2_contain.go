package str

// StringSetContains mark strings with bitmap
func StringSetContains(a, b string) bool {
	var bitmap uint32
	for _, r := range a {
		offset := r - 'a'
		bitmap |= uint32(1) << uint32(offset)
	}
	for _, r := range b {
		offset := r - 'a'
		mask := uint32(1) << uint32(offset)
		if bitmap&mask == 0 {
			return false
		}
	}
	return true
}
