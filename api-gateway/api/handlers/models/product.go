package models

type Product struct {
	Id          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Amount      int32   `json:"amount"`
	Created_at  string  `json:"created_at"`
	Updated_at  string  `json:"updated_at"`
}

type ListProducts struct {
	Count    int64      `json:"count"`
	Products []*Product `json:"products"`
}

type PurchasedProductsList struct {
	Products []*Product `json:"products"`
}

type BuyProductRequest struct {
	UserId    string `json:"user_id"`
	ProductId int32  `json:"product_id"`
	Amount    int32  `json:"amount"`
}

type BuyProductResponse struct {
	Message     string `json:"message"`
	UserId      string `json:"user_id"`
	ProductId   int32  `json:"product_id"`
	ProductName string `json:"product_name"`
	Amount      int32  `json:"amount"`
}

type ProductAmountRequest struct {
	ProductId int32 `json:"product_id"`
	Amount    int32 `json:"amount"`
}
