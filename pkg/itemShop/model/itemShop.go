package model

type (
	Item struct {
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Picture     string `json:"picture"`
		Price       uint   `json:"price"`
	}

	// Filter
	ItemFilter struct {
		Name        string `query:"name" validate:"omitempty,max=64"`
		Description string `query:"description" validate:"omitempty,max=128"`
		Paginate
	}

	// Pagination
	Paginate struct {
		Page int64 `query:"page" validate:"omitempty,min=1"`
		Size int64 `query:"size" validate:"omitempty,min=1,max=20"`
	}

	PaginateResult struct {
		Page      int64 `json:"page"`
		TotalPage int64 `json:"totalPage"`
	}

	ItemResult struct {
		Items    []*Item        `json:"items"`
		Paginate PaginateResult `json:"paginate"`
	}

	BuyingReq struct {
		PlayerID string
		ItemID   uint64 `json:"itemID" validation:"required,gt=0"`
		Quantity uint   `json:"quantity" validation:"required,gt=0"`
	}

	SellingReq struct {
		PlayerID string
		ItemID   uint64 `json:"itemID" validation:"required,gt=0"`
		Quantity uint   `json:"quantity" validation:"required,gt=0"`
	}
)
