package gen

type RadixTree[T any] struct {
	children map[string]*RadixTree[T]
	value    *T
}

func NewRadixTree[T any]() *RadixTree[T] {
	return &RadixTree[T]{
		children: make(map[string]*RadixTree[T]),
	}
}

func (t *RadixTree[T]) Add(path []string, value T) {
	dst := t
	for _, p := range path {
		if _, ok := dst.children[p]; !ok {
			dst.children[p] = NewRadixTree[T]()
		}
		dst = dst.children[p]
	}
	dst.value = &value
}

func (t *RadixTree[T]) Child(item string) *RadixTree[T] {
	if t == nil {
		return nil
	}

	return t.children[item]
}

// Children returns a tree with the same children and no value of its own.
func (t *RadixTree[T]) Children() *RadixTree[T] {
	if t == nil {
		return nil
	}
	return &RadixTree[T]{children: t.children}
}

func (t *RadixTree[T]) Value() (*T, bool) {
	if t == nil {
		return nil, false
	}
	return t.value, t.value != nil
}
