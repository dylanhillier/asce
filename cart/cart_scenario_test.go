package cart

// This file contains required scenario tests.

import (
	"testing"
)

type ProductCodeCount struct {
	prodCode string
	count    uint16
}

var scenarioTests = []struct {
	name                 string
	itemsToAdd           []ProductCodeCount
	promoCodes           []string
	expectedBundledItems []ProductCodeCount
	expectedCartTotal    PriceType
}{
	{
		// 3 for 2 deal on Unlimited 1GB Sim.
		"Scenario 1",
		[]ProductCodeCount{{"ult_small", 3}, {"ult_large", 1}},
		[]string{},
		[]ProductCodeCount{},
		9470,
	},
	{ // Unlimited 5GB Sim Bulk Deal.
		"Scenario 2",
		[]ProductCodeCount{{"ult_small", 2}, {"ult_large", 4}},
		[]string{},
		[]ProductCodeCount{},
		20940,
	},
	{ // Unlimited 2GB, Free 1GB Data Bundle.
		"Scenario 3",
		[]ProductCodeCount{{"ult_small", 1}, {"ult_medium", 2}},
		[]string{},
		[]ProductCodeCount{{"1gb", 2}},
		8470,
	},
	{ // Promo code 10% discount on cart.
		"Scenario 4",
		[]ProductCodeCount{{"ult_small", 1}, {"1gb", 1}},
		[]string{"I<3AMAYSIM"},
		[]ProductCodeCount{},
		3132,
	},
}

func Test_Launch_Scenarios(t *testing.T) {
	catalogue := CreateDefaultCatalogue()

	for _, tt := range scenarioTests {
		c := CreateCart(CreateDefaultRules(), catalogue)

		for _, pcc := range tt.itemsToAdd {
			if prod, ok := catalogue[pcc.prodCode]; ok {
				for i := uint16(0); i < pcc.count; i++ {
					c.Add(prod)
				}
			} else {
				t.Errorf("Test: %s contains product with code '%s' which was not found in launch product map.", tt.name, pcc.prodCode)
			}
		}

		for _, promoCode := range tt.promoCodes {
			c.AddPromoCode(promoCode)
		}

		var ua, ue []string

		// Check the standard cart items.
		ua, ue = checkItemsAgainstExpectations(t, c.Items(), tt.itemsToAdd)
		if len(ua) > 0 || len(ue) > 0 {
			t.Errorf("Test: %s expectation mismatch. UnexpectedCartItems=%s, UnmatchedExpectations=%s", tt.name, ua, ue)
		}

		// Check the bundled items.
		ua, ue = checkItemsAgainstExpectations(t, c.BundledItems(), tt.expectedBundledItems)
		if len(ua) > 0 || len(ue) > 0 {
			t.Errorf("Test: %s expectation mismatch. UnexpectedCartItems=%s, UnmatchedExpectations=%s", tt.name, ua, ue)
		}

		if c.Total() != tt.expectedCartTotal {
			t.Errorf("Test: %s. Cart total %d did not match expected total %d", tt.name, c.Total(), tt.expectedCartTotal)
		}
	}
}

func checkItemsAgainstExpectations(t *testing.T, actual ProductCollectionType, expected []ProductCodeCount) (unexpectedActualItems, unmatchedExpectedItems []string) {
	// Items not found in the actual
	for _, pcc := range expected {
		if apc, ok := actual[pcc.prodCode]; !ok {
			unmatchedExpectedItems = append(unmatchedExpectedItems, pcc.prodCode)
		} else {
			// Checks the item counts for the products which match.
			if pcc.count > apc.count {
				unmatchedExpectedItems = append(unmatchedExpectedItems, pcc.prodCode)
			} else if pcc.count < apc.count {
				unexpectedActualItems = append(unexpectedActualItems, pcc.prodCode)
			}
		}
	}

	// Items in the actual, but not expected.
	for prodCode, _ := range actual {
		found := false

		for _, epcc := range expected {
			if prodCode == epcc.prodCode {
				found = true
				break
			}
		}

		if !found {
			unexpectedActualItems = append(unexpectedActualItems, prodCode)
		}
	}

	return unexpectedActualItems, unmatchedExpectedItems
}

func Test_CartExpectationMatcher(t *testing.T) {
	c := CreateCart(CreateDefaultRules(), CreateDefaultCatalogue())

	testProduct := Product{"test_prod_code", "test product", 100}
	var ua []string // Unexpected Cart Items
	var ue []string // Unmatched Expectations

	// Match - no cart items, no expected products.
	ua, ue = checkItemsAgainstExpectations(t, c.Items(), nil)
	if len(ua) != 0 || len(ue) != 0 {
		t.Errorf("Failed on zero cart items, zero expected products.")
	}

	// No Match - no cart items, expected products.
	ua, ue = checkItemsAgainstExpectations(t, c.Items(), []ProductCodeCount{ProductCodeCount{"not_in_cart", 1}})
	if len(ua) != 0 || len(ue) != 1 {
		t.Errorf("Failed on zero cart items, one expected products.")
	}

	c.Add(testProduct)

	// No Match - cart items, no expected products.
	ua, ue = checkItemsAgainstExpectations(t, c.Items(), nil)
	if len(ua) != 1 || len(ue) != 0 {
		t.Errorf("Failed on zero cart items, one expected products.")
	}

	// Match - cart items, expected products.
	ua, ue = checkItemsAgainstExpectations(t, c.Items(), []ProductCodeCount{ProductCodeCount{testProduct.Code, 1}})
	if len(ua) != 0 || len(ue) != 0 {
		t.Errorf("Failed to match cart item against expected products.")
	}

	// No Match - cart items, expected products.
	ua, ue = checkItemsAgainstExpectations(t, c.Items(), []ProductCodeCount{ProductCodeCount{"not_in_cart", 1}})
	if len(ua) != 1 || len(ue) != 1 {
		t.Errorf("Failed to detect both unexpected cart items and unmatched expected products.")
	}
}
