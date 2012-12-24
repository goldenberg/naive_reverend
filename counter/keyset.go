package counter

type KeySet struct {
	keyToIdx map[string]int
	idxToKey []string
}

func NewKeySet() *KeySet {
	return &KeySet{make(map[string]int), make([]string, 2)}
}

func NewKeySetFromKeys(keys []string) *KeySet {
	keyToIdx := make(map[string]int, len(keys))
	for i, k := range keys {
		keyToIdx[k] = i
	}
	return &KeySet{keyToIdx, keys}
}

func (ks *KeySet) Get(k string) (idx int, ok bool) {
	idx, ok = ks.keyToIdx[k]
	if !ok {
		ks.keyToIdx[k] = ks.Len()
	}
	return
}

func (ks *KeySet) Keys() (result []string) {
	for k := range ks.keyToIdx {
		result = append(result, k)
	}
	return
}

func (ks *KeySet) Len() int {
	return len(ks.idxToKey)
}
