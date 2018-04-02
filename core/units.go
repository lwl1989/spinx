package core

import "encoding/binary"

//func readSize(s []byte) (uint32, int) {
//	if len(s) == 0 {
//		return 0, 0
//	}
//	size, n := uint32(s[0]), 1
//	if size&(1<<7) != 0 {
//		if len(s) < 4 {
//			return 0, 0
//		}
//		n = 4
//		size = binary.BigEndian.Uint32(s)
//		size &^= 1 << 31
//	}
//	return size, n
//}

//func readString(s []byte, size uint32) string {
//	if size > uint32(len(s)) {
//		return ""
//	}
//	return string(s[:size])
//}

func encodeSize(b []byte, size uint32) int {
	if size > 127 {
		size |= 1 << 31
		binary.BigEndian.PutUint32(b, size)
		return 4
	}
	b[0] = byte(size)
	return 1
}
