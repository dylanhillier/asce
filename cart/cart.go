package cart

import (
	"fmt"
)

type Cart interface {
	Add(Product)
	Remove(Product)
	AddPromoCode(string)
	RemovePromoCode(string)
	Clear()
	Items() ProductCollectionType
	BundledItems() ProductCollectionType
	PromoCodes() []string
	Total() PriceType
}

func CreateCart(rules []Rule, catalogue Catalogue) Cart {
	return &defaultCart{
		catalogue:         catalogue,
		products:          make(ProductCollectionType),
		bundleProducts:    make(ProductCollectionType),
		promoCodes:        make(map[string]bool),
		rules:             rules,
		undiscountedTotal: 0,
		discount:          0,
	}
}

type defaultCart struct {
	catalogue         Catalogue
	products          ProductCollectionType
	bundleProducts    ProductCollectionType // These are imutable by interface methods Add/Remove.
	promoCodes        map[string]bool       // This is a set.
	rules             []Rule
	undiscountedTotal PriceType // Total of Products in cart without offers/promotions applied.
	discount          PriceType // Discount applied due to triggered rules.
}

func (c *defaultCart) Add(p Product) {
	if v, ok := c.products[p.Code]; !ok {
		//fmt.Printf("Adding %s to the cart. Count=1\n", p.Code)
		c.products[p.Code] = &ProductCount{p, 1}
	} else {
		v.count++
		//fmt.Printf("Adding %s to the cart. Count=%d\n", p.Code, v.count)
	}

	c.undiscountedTotal += p.Price
	c.evaluateRules()
}

func (c *defaultCart) Remove(p Product) {
	if v, ok := c.products[p.Code]; ok {
		v.count--
		if v.count == 0 {
			//fmt.Printf("Removed last %s from the cart.\n", p.Code)
			delete(c.products, p.Code)
		} else {
			//fmt.Printf("Removing %s from the cart. Count=%d\n", p.Code, v.count)
		}

		// TODO: Should put validation here to ensure it doesn't ever go negative.
		c.undiscountedTotal -= p.Price
		c.evaluateRules()
	}
}

func (c *defaultCart) AddPromoCode(code string) {
	c.promoCodes[code] = true
	c.evaluateRules()
}

func (c *defaultCart) RemovePromoCode(code string) {
	delete(c.promoCodes, code)
	c.evaluateRules()
}

func (c *defaultCart) PromoCodes() []string {
	codes := []string{}
	for k, _ := range c.promoCodes {
		codes = append(codes, k)
	}

	return codes
}

func (c *defaultCart) Clear() {
	c.products = make(ProductCollectionType)
	c.bundleProducts = make(ProductCollectionType)
	c.promoCodes = make(map[string]bool)
}

func (c *defaultCart) Items() ProductCollectionType {
	return c.products
}

func (c *defaultCart) BundledItems() ProductCollectionType {
	return c.bundleProducts
}

func (c *defaultCart) Total() PriceType {
	return c.undiscountedTotal - c.discount
}

func (c *defaultCart) evaluateRules() {
	c.discount = 0
	c.bundleProducts = make(ProductCollectionType)

	for _, rule := range c.rules {
		discount, bp := rule.Evaluate(c)
		c.discount += discount

		if bp.count != 0 && bp.code != "" {

			if v, ok := c.bundleProducts[bp.code]; !ok {
				if product, ok := c.catalogue[bp.code]; !ok {
					fmt.Errorf("Failed to find %v in product map. This implies a rule is setup for a product which doesnt exist.\n", bp.code)
				} else {
					c.bundleProducts[bp.code] = &ProductCount{product, bp.count}
				}
			} else {
				v.count += bp.count
			}
		}
	}
}
