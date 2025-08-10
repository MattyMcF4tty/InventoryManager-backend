package schemas

type SupplierContactInfo struct {
	Id          int8   `json:"id"`
	SupplierId  int8   `json:"supplier_id"`
	ContactName string `json:"contact_name"`
	Role        string `json:"role"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}
