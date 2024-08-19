package shared

func (b ByteSlice) IsEmpty() bool {
	return len(b) == 0 || b == nil
}
