package schemas

type Item struct {
	Id          int8    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int8    `json:"quantity"`
	Category    string  `json:"category"`
	ImageUrl    *string `json:"image_url,omitempty"`
	SupplierId  int8    `json:"supplier_id"`

	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}
