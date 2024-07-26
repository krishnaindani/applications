package main

type Repo interface {
	InsertProduct(product Product) error
	UpdateProduct(product Product) error
	DeleteProduct(product Product) error
	GetProducts() ([]Product, error)
	Transaction(func(repo Repo) error) error
}

func (pup ProductUpdatePlan) Execute(repo Repo) error {

	for _, product := range pup.Added {
		err := repo.InsertProduct(product)
		if err != nil {
			return err
		}
	}

	for _, productPair := range pup.Updated {
		err := repo.UpdateProduct(productPair.NewProduct)
		if err != nil {
			return err
		}
	}

	for _, product := range pup.Deleted {
		err := repo.DeleteProduct(product)
		if err != nil {
			return err
		}
	}

	return nil
}
