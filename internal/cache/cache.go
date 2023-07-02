package cache

type CacheContainer interface {
	GetKeymaps() struct{}
	GetPolicyAlgorithm() string
}
