package cart

type Catalogue map[string]Product

func CreateDefaultCatalogue() Catalogue {
	return Catalogue{
		"ult_small":  Product{"ult_small", "Unlimited 1GB", 2490},
		"ult_medium": Product{"ult_medium", "Unlimited 2GB", 2990},
		"ult_large":  Product{"ult_large", "Unlimited 5GB", 4490},
		"1gb":        Product{"1gb", "1GB Data-pack", 990},
	}
}
