package cuckoo

import (
	"fmt"
	"math/bits"
	"math/rand"
)

const maxCuckooCount = 500

type Filter struct {
	// 用于存储指纹
	buckets []bucket
	// 当前过滤器中存储的元素数量
	count uint
	// 用于计算桶索引的幂
	bucketPow uint
}

// NewFilter 创建指定容量的过滤器
// 1,000,000 是常用的默认值大小，在 64bit 计算机上占用 1MB 左右内存
func NewFilter(capacity uint) *Filter {
	capacity = getNextPow2(uint64(capacity)) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	buckets := make([]bucket, capacity)
	return &Filter{
		buckets:   buckets,
		count:     0,
		bucketPow: uint(bits.TrailingZeros(capacity)),
	}
}

// Lookup 检查数据是否在过滤器中
func (cf *Filter) Lookup(data []byte) bool {
	// 检查第一个桶
	i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
	if cf.buckets[i1].getFingerprintIndex(fp) > -1 {
		return true
	}
	// 检查备用桶
	i2 := getAltIndex(fp, i1, cf.bucketPow)
	return cf.buckets[i2].getFingerprintIndex(fp) > -1
}

// Reset 重置
func (cf *Filter) Reset() {
	for i := range cf.buckets {
		cf.buckets[i].reset()
	}
	cf.count = 0
}

// 随机返回两个索引中的一个
func randi(i1, i2 uint) uint {
	if rand.Intn(2) == 0 {
		return i1
	}
	return i2
}

// Insert 插入数据
func (cf *Filter) Insert(data []byte) bool {
	// 尝试插入第一个桶
	i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
	if cf.insert(fp, i1) {
		return true
	}
	// 插入失败尝试插入备用桶
	i2 := getAltIndex(fp, i1, cf.bucketPow)
	if cf.insert(fp, i2) {
		return true
	}
	// 仍然失败则尝试重新插入
	return cf.reinsert(fp, randi(i1, i2))
}

func (cf *Filter) insert(fp fingerprint, i uint) bool {
	if cf.buckets[i].insert(fp) {
		cf.count++
		return true
	}
	return false
}

// 重新插入
func (cf *Filter) reinsert(fp fingerprint, i uint) bool {
	// 在最大尝试次数内
	for k := 0; k < maxCuckooCount; k++ {
		// 随机将桶中现有指纹踢出去，并将该指纹放在该位置上
		j := rand.Intn(bucketSize)
		oldfp := fp
		fp = cf.buckets[i][j]
		cf.buckets[i][j] = oldfp

		// 将该被踢出去的元素在备用桶中寻找位置
		i = getAltIndex(fp, i, cf.bucketPow)
		if cf.insert(fp, i) {
			return true
		}
	}
	return false
}

// InsertUnique 插入不存在的元素
func (cf *Filter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	return cf.Insert(data)
}

// Delete 删除过滤器中的指纹
func (cf *Filter) Delete(data []byte) bool {
	i1, fp := getIndexAndFingerprint(data, cf.bucketPow)
	if cf.delete(fp, i1) {
		return true
	}
	// 删除失败，则尝试从备用桶删除
	i2 := getAltIndex(fp, i1, cf.bucketPow)
	return cf.delete(fp, i2)
}

// 删除指定桶的指纹
func (cf *Filter) delete(fp fingerprint, i uint) bool {
	if cf.buckets[i].delete(fp) {
		if cf.count > 0 {
			cf.count--
		}
		return true
	}
	return false
}

// Count 过滤器中元素数量
func (cf *Filter) Count() uint {
	return cf.count
}

// Encode 编码过滤器
func (cf *Filter) Encode() []byte {
	bytes := make([]byte, len(cf.buckets)*bucketSize)
	for i, b := range cf.buckets {
		for j, f := range b {
			index := (i * len(b)) + j
			bytes[index] = byte(f)
		}
	}
	return bytes
}

// Decode 将过滤器解码
func Decode(bytes []byte) (*Filter, error) {
	var count uint
	if len(bytes)%bucketSize != 0 {
		return nil, fmt.Errorf("expected bytes to be multiple of %d, got %d", bucketSize, len(bytes))
	}
	if len(bytes) == 0 {
		return nil, fmt.Errorf("bytes can not be empty")
	}
	buckets := make([]bucket, len(bytes)/4)
	for i, b := range buckets {
		for j := range b {
			index := (i * len(b)) + j
			if bytes[index] != 0 {
				buckets[i][j] = fingerprint(bytes[index])
				count++
			}
		}
	}
	return &Filter{
		buckets:   buckets,
		count:     count,
		bucketPow: uint(bits.TrailingZeros(uint(len(buckets)))),
	}, nil
}
