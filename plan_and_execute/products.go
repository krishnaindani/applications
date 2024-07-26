package main

type Product struct {
	ID    string
	Price int
}

type ProductPair struct {
	OldProduct Product
	NewProduct Product
}

type ProductUpdatePlan struct {
	Added     []Product
	NoChanges []Product
	Updated   []ProductPair
	Deleted   []Product
}

func CreatePlan(
	oldProducts []Product,
	newProducts []Product,
) ProductUpdatePlan {

	plan := ProductUpdatePlan{}

	oldProductMap := make(map[string]Product)
	for _, product := range oldProducts {
		oldProductMap[product.ID] = product
	}

	newProductMap := make(map[string]Product)
	for _, product := range newProducts {
		oldProduct, oldProductExists := oldProductMap[product.ID]
		if !oldProductExists {
			plan.Added = append(plan.Added, product)
		} else if oldProduct.Price != product.Price {
			plan.Updated = append(plan.Updated, ProductPair{
				OldProduct: oldProduct,
				NewProduct: product,
			})
		} else {
			plan.NoChanges = append(plan.NoChanges, product)
		}
		newProductMap[product.ID] = product
	}

	for _, product := range oldProducts {
		if _, ok := newProductMap[product.ID]; !ok {
			plan.Deleted = append(plan.Deleted, product)
		}
	}

	return plan
}

func UpdateProducts(repo Repo,
	csv string,
	preview bool,
) (ProductUpdatePlan, error) {

	newProducts, err := parseCSV()
	if err != nil {
		return ProductUpdatePlan{}, err
	}

	var plan ProductUpdatePlan
	err = repo.Transaction(func(repo Repo) error {
		oldProducts, err := repo.GetProducts()
		if err != nil {
			return err
		}

		plan := CreatePlan(oldProducts, newProducts)
		if preview {
			return nil
		}

		return plan.Execute(repo)
	})
	if err != nil {
		return ProductUpdatePlan{}, err
	}

	return plan, nil
}
