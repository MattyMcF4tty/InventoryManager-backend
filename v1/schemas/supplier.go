package schemas

type Supplier struct {
	Id          int8                  `json:"id"`
	Name        string                `json:"name"`
	ContactInfo []SupplierContactInfo `json:"contact_info"`
	Website     string                `json:"website"`
	Address     string                `json:"address"`
	VatNumber   string                `json:"vat_number"`

	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}
