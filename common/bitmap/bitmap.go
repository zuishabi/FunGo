package bitmap

// ...existing code...

// IsBitSet 返回位图 bitmap 中 pos 位置是否为 1。
// pos 从 0 开始，按每字节低位为 bit 0 处理。
// 若 pos 超出 bitmap 长度，返回 false（表示未点赞）。
// ...existing code...
func IsBitSet(bitmap []byte, pos uint64) bool {
	byteIndex := pos / 8
	bitOffset := pos % 8
	if int(byteIndex) >= len(bitmap) {
		return false
	}
	return bitmap[byteIndex]&(1<<bitOffset) != 0
}

// SetBit 在位图 bitmap 的 pos 位置设置为 1，必要时会扩展切片并返回新的切片引用。
// ...existing code...
func SetBit(bitmap []byte, pos uint64) []byte {
	byteIndex := pos / 8
	bitOffset := pos % 8
	if int(byteIndex) >= len(bitmap) {
		newBuf := make([]byte, byteIndex+1)
		copy(newBuf, bitmap)
		bitmap = newBuf
	}
	bitmap[byteIndex] |= 1 << bitOffset
	return bitmap
}

// ClearBit 在位图 bitmap 的 pos 位置清零（设置为 0）。若超出长度则不做任何事。
// ...existing code...
func ClearBit(bitmap []byte, pos uint64) []byte {
	byteIndex := pos / 8
	bitOffset := pos % 8
	if int(byteIndex) >= len(bitmap) {
		return bitmap
	}
	bitmap[byteIndex] &^= 1 << bitOffset
	return bitmap
}
