package code

//go:generate codeg -type=int -doc -output ./error_code_generated.md
const (
	// ErrGoodsNotFound - 404: Goods not found.
	ErrGoodsNotFound int = iota + 100501

	// ErrCategoryNotFound - 404: Category not found.
	ErrCategoryNotFound

	// ErrBrandNotFound - 500: Es unmarshal error.
	ErrEsUnmarshal
)
