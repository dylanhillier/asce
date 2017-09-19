package cart

import (
	"testing"
)

func compareActualAgainstExpectation(t *testing.T, actualDiscount PriceType, actualBundleProduct BundledProduct, expectedDiscount PriceType, expectedBundleProduct BundledProduct) {
	if actualDiscount != expectedDiscount {
		t.Errorf("ActualDiscount: %v ExpectedDiscount: %v", actualDiscount, expectedDiscount)
	}

	if actualBundleProduct != expectedBundleProduct {
		t.Errorf("ActualBundleProduct: %v ExpectedBundledProduct: %v", actualBundleProduct, expectedBundleProduct)
	}
}

func Test_XForYRule_WHEN_InsufficientProductsAdded_EXPECT_NoDiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_small"]
	rule := CreateXForYRule(product.Code, 2, 1)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := PriceType(0)
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_XForYRule_WHEN_SufficientProductsAdded_EXPECT_DiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_small"]
	rule := CreateXForYRule(product.Code, 2, 1)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := product.Price
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_XForYRule_WHEN_MultiplesOfSufficientProductsAdded_EXPECT_DiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_small"]
	rule := CreateXForYRule(product.Code, 2, 1)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	cart.Add(product)
	cart.Add(product)
	cart.Add(product)
	cart.Add(product) // Odd Number - No discount on this item.
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := product.Price * 2
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_BulkDiscountRule_WHEN_InsufficientProductsAdded_EXPECT_NoDiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_large"]
	discount := PriceType(500)
	rule := CreateBulkDiscountRule(product.Code, 3, discount)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := PriceType(0)
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_BulkDiscountRule_WHEN_SufficientProductsAdded_EXPECT_DiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_large"]
	discount := PriceType(500)
	rule := CreateBulkDiscountRule(product.Code, 3, discount)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	cart.Add(product)
	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := discount * 3
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_BulkDiscountRule_WHEN_MoreThanSufficientProductsAdded_EXPECT_DiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_large"]
	discount := PriceType(500)
	rule := CreateBulkDiscountRule(product.Code, 3, discount)
	cart := CreateCart([]Rule{rule}, catalogue)

	itemsToAdd := 20
	for i := 0; i < itemsToAdd; i++ {
		cart.Add(product)
	}
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := discount * PriceType(itemsToAdd)
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_BundleRule_WHEN_InsufficientProductsAdded_EXPECT_NoDiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_small"]
	bundleProduct := catalogue["1gb"]
	rule := CreateBundleRule(product.Code, 2, bundleProduct.Code, 1)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := PriceType(0)
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_BundleRule_WHEN_ExactNumProductsAdded_EXPECT_NoDiscountAndBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_small"]
	bundleProduct := catalogue["1gb"]
	rule := CreateBundleRule(product.Code, 2, bundleProduct.Code, 1)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := PriceType(0)
	expectedBundleProduct := BundledProduct{bundleProduct.Code, 1}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_BundleRule_WHEN_MultiplesOfTiggerProductCountAdded_EXPECT_NoDiscountAndBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_small"]
	bundleProduct := catalogue["1gb"]
	rule := CreateBundleRule(product.Code, 2, bundleProduct.Code, 1)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	cart.Add(product)
	cart.Add(product)
	cart.Add(product)
	cart.Add(product) // This one does not trigger an additional bundle deal.
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := PriceType(0)
	expectedBundleProduct := BundledProduct{bundleProduct.Code, 2}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_PromoRule_WHEN_NoPromoCodeAdded_EXPECT_NoDiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_large"]
	promoCode := "MoreCowbell"
	discountPct := int8(20)
	rule := CreatePromoRule(promoCode, discountPct)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := PriceType(0)
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}

func Test_PromoRule_WHEN_PromoCodeAdded_EXPECT_DiscountAndNoBundledItems(t *testing.T) {
	catalogue := CreateDefaultCatalogue()
	product := catalogue["ult_large"]
	promoCode := "MoreCowbell"
	discountPct := int8(20)
	rule := CreatePromoRule(promoCode, discountPct)
	cart := CreateCart([]Rule{rule}, catalogue)

	cart.Add(product)
	cart.AddPromoCode(promoCode)
	actualDiscount, actualBundleProduct := rule.Evaluate(cart)

	expectedDiscount := percentageOfPrice(product.Price, discountPct)
	expectedBundleProduct := BundledProduct{}

	compareActualAgainstExpectation(t, actualDiscount, actualBundleProduct, expectedDiscount, expectedBundleProduct)
}
