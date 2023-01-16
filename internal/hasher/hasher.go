package hasher

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
	pow2     = float64(int64(1) << 31)
)

func SumWithNum(key string, shardsAmount int) (sum uint64, num int) {
	sum = Sum(key)
	num = JCH(sum, shardsAmount)

	return sum, num
}

func Sum(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}

// https://arxiv.org/pdf/1406.2294.pdf
// dgryski golang implementation https://github.com/dgryski/go-jump/blob/master/jump.go
func JCH(hashedKey uint64, shardsAmount int) int {
	var b int64 = -1
	var j int64
	key := hashedKey

	for j < int64(shardsAmount) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (pow2 / float64((key>>33)+1)))
	}

	return int(b)
}
