package cart

import (
	"testing"
)

// TODO:
// 1. Test that given a pricing rule, it triggers as expected when the correct items are added to the cart.
// 2. Test that a triggered pricing rule reverts if an item is removed which triggered the discount.
// 3. Test promocodes.

func checkCartContainsNProductsWithCode(t *testing.T, cart Cart, prodCode string, expectedCount uint16) {
	items := cart.Items()

	if v, ok := items[prodCode]; !ok {
		if expectedCount > 0 {
			t.Errorf("Cart didn't contain the expected product code: %s", prodCode)
		}
	} else {
		if v.count != expectedCount {
			t.Errorf("Cart contained %d products with code %s. Expected %d",
				v.count,
				prodCode,
				expectedCount)
		}
	}
}

func Test_Cart_GIVEN_EmptyCart_WHEN_ProductAdded_EXPECT_CartHasItem(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	p := Product{"nin_year_zero", "Nine Inch Nails - Year Zero", 3000}

	checkCartContainsNProductsWithCode(t, c, p.Code, 0)

	c.Add(p)

	checkCartContainsNProductsWithCode(t, c, p.Code, 1)
}

func Test_Cart_GIVEN_EmptyCart_WHEN_ProductAddedMultipleTimes_EXPECT_CartHasItems(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	p := Product{"nin_the_slip", "Nine Inch Nails - The Slip", 3000}

	checkCartContainsNProductsWithCode(t, c, p.Code, 0)

	c.Add(p)
	c.Add(p)

	checkCartContainsNProductsWithCode(t, c, p.Code, 2)
}

func Test_Cart_GIVEN_EmptyCart_WHEN_ProductAdded_THEN_ProductRemoved_EXPECT_CartIsEmpty(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	p := Product{"nin_hesitation_marks", "Nine Inch Nails - Hesitation Marks", 3000}

	checkCartContainsNProductsWithCode(t, c, p.Code, 0)

	c.Add(p)

	checkCartContainsNProductsWithCode(t, c, p.Code, 1)

	c.Remove(p)

	checkCartContainsNProductsWithCode(t, c, p.Code, 0)
}

func Test_Cart_GIVEN_NonEmptyCart_WHEN_Cleared_EXPECT_CartIsEmpty(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	p1 := Product{"nin_downward_spiral", "Nine Inch Nails - Downward Spiral", 3000}
	p2 := Product{"nin_with_teeth", "Nine Inch Nails - With Teeth", 3000}

	c.Add(p1)
	c.Add(p2)

	checkCartContainsNProductsWithCode(t, c, p1.Code, 1)
	checkCartContainsNProductsWithCode(t, c, p2.Code, 1)

	c.Clear()
	itemCount := len(c.Items())
	if itemCount > 0 {
		t.Errorf("Cart clear failed. %d items remaining.", itemCount)
	}
}

func Test_Cart_WHEN_ProductsAdded_EXPECT_TotalAndItemsToIncrease(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	p1 := Product{"nin_downward_spiral", "Nine Inch Nails - Downward Spiral", 3000}
	p2 := Product{"nin_with_teeth", "Nine Inch Nails - With Teeth", 5000}

	if c.Total() != 0 {
		t.Errorf("CartTotal=%d, Expected=0", c.Total())
	}

	if len(c.Items()) != 0 {
		t.Errorf("CartItemCount=%d, Expected=0", len(c.Items()))
	}

	c.Add(p1)

	if c.Total() != p1.Price {
		t.Errorf("CartTotal=%d, Expected=%d", c.Total(), p1.Price)
	}

	if len(c.Items()) != 1 {
		t.Errorf("CartItemCount=%d, Expected=1", len(c.Items()))
	}

	c.Add(p2)
	if c.Total() != p1.Price+p2.Price {
		t.Errorf("CartTotal=%d, Expected=%d", c.Total(), p1.Price+p2.Price)
	}

	if len(c.Items()) != 2 {
		t.Errorf("CartItemCount=%d, Expected=2", len(c.Items()))
	}

}

func Test_Cart_WHEN_ProductsRemoved_EXPECT_TotalAndItemsToDecrease(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	p1 := Product{"nin_downward_spiral", "Nine Inch Nails - Downward Spiral", 3000}
	p2 := Product{"nin_with_teeth", "Nine Inch Nails - With Teeth", 5000}

	if c.Total() != 0 {
		t.Errorf("CartTotal=%d, Expected=0", c.Total())
	}

	if len(c.Items()) != 0 {
		t.Errorf("CartItemCount=%d, Expected=0", len(c.Items()))
	}

	c.Add(p1)
	c.Add(p2)

	if c.Total() != p1.Price+p2.Price {
		t.Errorf("CartTotal=%d, Expected=%d", c.Total(), p1.Price+p2.Price)
	}

	if len(c.Items()) != 2 {
		t.Errorf("CartItemCount=%d, Expected=2", len(c.Items()))
	}

	c.Remove(p1)
	if c.Total() != p2.Price {
		t.Errorf("CartTotal=%d, Expected=%d", c.Total(), p2.Price)
	}

	if len(c.Items()) != 1 {
		t.Errorf("CartItemCount=%d, Expected=1", len(c.Items()))
	}
}

func Test_Cart_WHEN_PromoCodeAdded_EXPECT_CartRetainsPromoCode(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	expectedPromoCode := "1337"
	if len(c.PromoCodes()) != 0 {
		t.Errorf("ActualPromoCodes=%d ExpectedPromoCodes=0", len(c.PromoCodes()))
	}

	c.AddPromoCode(expectedPromoCode)

	if len(c.PromoCodes()) != 1 {
		t.Errorf("ActualPromoCodes=%v ExpectedPromoCodes=0", c.PromoCodes())
	} else {
		actualPromoCode := c.PromoCodes()[0]
		if actualPromoCode != expectedPromoCode {
			t.Errorf("ActualPromoCode=%v ExpectedPromoCode=%v", actualPromoCode, expectedPromoCode)
		}
	}
}

func Test_Cart_WHEN_DuplicatePromoCodeAdded_EXPECT_CartRetainsOnlyUniquePromoCode(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	expectedPromoCode := "1337"
	if len(c.PromoCodes()) != 0 {
		t.Errorf("ActualPromoCodes=%d ExpectedPromoCodes=0", len(c.PromoCodes()))
	}

	c.AddPromoCode(expectedPromoCode)

	if len(c.PromoCodes()) != 1 {
		t.Errorf("ActualPromoCodes=%v ExpectedPromoCodes=0", c.PromoCodes())
	} else {
		actualPromoCode := c.PromoCodes()[0]
		if actualPromoCode != expectedPromoCode {
			t.Errorf("ActualPromoCode=%v ExpectedPromoCode=%v", actualPromoCode, expectedPromoCode)
		}
	}

	c.AddPromoCode(expectedPromoCode)

	if len(c.PromoCodes()) != 1 {
		t.Errorf("ActualPromoCodes=%v ExpectedPromoCodes=0", c.PromoCodes())
	} else {
		actualPromoCode := c.PromoCodes()[0]
		if actualPromoCode != expectedPromoCode {
			t.Errorf("ActualPromoCode=%v ExpectedPromoCode=%v", actualPromoCode, expectedPromoCode)
		}
	}
}

func Test_Cart_WHEN_PromoCodeRemoved_EXPECT_CartRemovesPromoCode(t *testing.T) {
	c := CreateCart(nil, CreateDefaultCatalogue())
	expectedPromoCode := "1337"
	if len(c.PromoCodes()) != 0 {
		t.Errorf("ActualPromoCodes=%d ExpectedPromoCodes=0", len(c.PromoCodes()))
	}

	c.AddPromoCode(expectedPromoCode)

	if len(c.PromoCodes()) != 1 {
		t.Errorf("ActualPromoCodes=%v ExpectedPromoCodes=0", c.PromoCodes())
	} else {
		actualPromoCode := c.PromoCodes()[0]
		if actualPromoCode != expectedPromoCode {
			t.Errorf("ActualPromoCode=%v ExpectedPromoCode=%v", actualPromoCode, expectedPromoCode)
		}
	}

	c.RemovePromoCode(expectedPromoCode)

	if len(c.PromoCodes()) != 0 {
		t.Errorf("ActualPromoCodes=%v ExpectedPromoCodes=0", c.PromoCodes())
	}
}
