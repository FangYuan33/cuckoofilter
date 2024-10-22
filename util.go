package cuckoo

import (
	metro "github.com/dgryski/go-metro"
)

var (
	// 指纹哈希值
	altHash = [256]uint{}
	// 掩码值
	masks = [65]uint{}
)

func init() {
	for i := 0; i < 256; i++ {
		altHash[i] = (uint(metro.Hash64([]byte{byte(i)}, 1337)))
	}
	for i := uint(0); i <= 64; i++ {
		masks[i] = (1 << i) - 1
	}
}

// 获取备用桶索引
func getAltIndex(fp fingerprint, i uint, bucketPow uint) uint {
	mask := masks[bucketPow]
	hash := altHash[fp] & mask
	return (i & mask) ^ hash
}

// getIndicesAndFingerprint 获取该元素桶的索引值，并且计算出该元素的指纹
func getIndexAndFingerprint(data []byte, bucketPow uint) (uint, fingerprint) {
	hash := defaultHasher.Hash64(data)
	fp := getFingerprint(hash)
	// 使用哈希值的高位部分和掩码计算第一个桶索引
	i1 := uint(hash>>32) & masks[bucketPow]
	return i1, fingerprint(fp)
}

// 使用哈希值的低位部分计算指纹，使得指纹在 1 - 255 之间
func getFingerprint(hash uint64) byte {
	fp := byte(hash%255 + 1)
	return fp
}

// 计算出大于或等于 n 的下一个 2 的幂
func getNextPow2(n uint64) uint {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return uint(n)
}

var defaultHasher Hasher = new(metrotHasher)

func SetDefaultHasher(hasher Hasher) {
	defaultHasher = hasher
}

type Hasher interface {
	Hash64([]byte) uint64
}

var _ Hasher = new(metrotHasher)

type metrotHasher struct{}

func (h *metrotHasher) Hash64(data []byte) uint64 {
	hash := metro.Hash64(data, 1337)
	return hash
}
