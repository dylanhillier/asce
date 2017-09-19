package cart

type BundledProduct struct {
	code  string
	count uint16
}

type Rule interface {
	Evaluate(Cart) (discount PriceType, bundledProduct BundledProduct)
}

func CreateDefaultRules() []Rule {
	return []Rule{
		// 3 for 2 deal on Unlimited 1GB Sim.
		CreateXForYRule("ult_small", 3, 2),
		// Unlimited 5GB Sim Bulk Deal.
		CreateBulkDiscountRule("ult_large", 3, 500),
		// Unlimited 2GB, Free 1GB Data Bundle.
		CreateBundleRule("ult_medium", 1, "1gb", 1),
		// Promo code 10% discount on cart.
		CreatePromoRule("I<3AMAYSIM", 10),
	}
}

func CreateXForYRule(prodCode string, x, y uint16) Rule {
	return &xForYRule{prodCode, x, y}
}

func CreateBulkDiscountRule(prodCode string, countToExceed uint16, discountAbs PriceType) Rule {
	return &bulkDiscountRule{prodCode, countToExceed, discountAbs}
}

func CreateBundleRule(buyProdCode string, itemsToBuy uint16, getProdCode string, itemsToGet uint16) Rule {
	return &bundleRule{buyProdCode, itemsToBuy, getProdCode, itemsToGet}
}

func CreatePromoRule(code string, discountPct int8) Rule {
	return &promoRule{code, discountPct}
}

func percentageOfPrice(price PriceType, pct int8) PriceType {
	return PriceType(float32(price) * (float32(pct) / float32(100)))
}

type xForYRule struct {
	prodCode string
	x        uint16
	y        uint16
}

func (r *xForYRule) Evaluate(c Cart) (discount PriceType, bundledProduct BundledProduct) {

	if v, ok := c.Items()[r.prodCode]; !ok {
		return 0, BundledProduct{}
	} else {
		if v.count >= r.x {
			timesToApplyDiscount := uint16(v.count / r.x)
			discount := PriceType((r.x-r.y)*timesToApplyDiscount) * v.product.Price
			return discount, BundledProduct{}
		}
	}

	return 0, BundledProduct{}
}

type bulkDiscountRule struct {
	prodCode      string
	countToExceed uint16
	discountAbs   PriceType
}

func (r *bulkDiscountRule) Evaluate(c Cart) (discount PriceType, bundledProduct BundledProduct) {
	if v, ok := c.Items()[r.prodCode]; !ok {
		return 0, BundledProduct{}
	} else {
		if v.count >= r.countToExceed {
			discount := PriceType(v.count) * r.discountAbs
			return discount, BundledProduct{}
		}
	}

	return 0, BundledProduct{}
}

type promoRule struct {
	code        string
	discountPct int8
}

func (r *promoRule) Evaluate(c Cart) (discount PriceType, bundledProduct BundledProduct) {
	var cartTotal PriceType = 0

	codeFound := false
	for _, code := range c.PromoCodes() {
		if code == r.code {
			codeFound = true
			break
		}
	}

	if !codeFound {
		return 0, BundledProduct{}
	}

	for _, v := range c.Items() {
		cartTotal += PriceType(v.count) * v.product.Price
	}

	discount = percentageOfPrice(cartTotal, r.discountPct)

	return discount, BundledProduct{}
}

type bundleRule struct {
	buyProdCode string
	itemsToBuy  uint16
	getProdCode string
	itemsToGet  uint16
}

func (r *bundleRule) Evaluate(c Cart) (discount PriceType, bundledProduct BundledProduct) {
	if v, ok := c.Items()[r.buyProdCode]; !ok {
		return 0, BundledProduct{}
	} else {
		if v.count >= r.itemsToBuy {
			items := (v.count / r.itemsToBuy) * r.itemsToGet
			return 0, BundledProduct{r.getProdCode, items}
		}
	}

	return 0, BundledProduct{}
}
