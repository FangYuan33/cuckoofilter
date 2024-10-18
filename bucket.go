package cuckoo

// 定义指纹类型，它实际上是 byte 的别名
type fingerprint byte

// 定义 bucket 类型，它是一个固定大小为 `bucketSize` 的 `fingerprint` 数组
type bucket [bucketSize]fingerprint

const (
	// 表示空指纹，用于标识桶中的空槽
	nullFp = 0
	// 桶的大小，即每个桶可以存储的指纹数量为 4
	bucketSize = 4
)

// 插入新的指纹
func (b *bucket) insert(fp fingerprint) bool {
	// 找到空指纹，并将该新指纹写入
	for i, tfp := range b {
		if tfp == nullFp {
			b[i] = fp
			return true
		}
	}
	return false
}

// 删除指纹
func (b *bucket) delete(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == fp {
			b[i] = nullFp
			return true
		}
	}
	return false
}

// 获取指纹的索引
func (b *bucket) getFingerprintIndex(fp fingerprint) int {
	for i, tfp := range b {
		if tfp == fp {
			return i
		}
	}
	return -1
}

// 重置指纹
func (b *bucket) reset() {
	for i := range b {
		b[i] = nullFp
	}
}
