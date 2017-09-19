package cart

// Make it easy to change the price type everywhere if i change my mind.
type PriceType int32

type ProductCount struct {
	product Product
	count   uint16
}

// Type representing a collection of products as they
// would appear in a cart.
// "product code" => {Product, CartCount}
type ProductCollectionType map[string]*ProductCount

type Product struct {
	Code  string
	Name  string
	Price PriceType
}
