package model

type PurchaseOrder struct {
	ID              int    `json:"id,omitempty"`
	NamaSuplier     string `json:"nama_suplier,omitempty"`
	NomorPO         string `json:"nomor_po,omitempty"`
	Tanggal         string `json:"tanggal,omitempty"`
	CatatanPO       string `json:"catatan_po,omitempty"`
	PreparedBy      string `json:"prepared_by,omitempty"`
	PreparedJabatan string `json:"prepared_jabatan,omitempty"`
	ApprovedBy      string `json:"approved_by,omitempty"`
	ApprovedJabatan string `json:"approved_jabatan,omitempty"`
	Divisi          string `json:"divisi,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	CreatedBy       string `json:"created_by,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	UpdatedBy       string `json:"updated_by,omitempty"`
}

// ItemBuyer represents the item_buyer table
type ItemBuyer struct {
	ID       int    `json:"id,omitempty"`
	POID     int    `json:"po_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Quantity string `json:"quantity,omitempty"`
	Price    string `json:"price,omitempty"`
	Discount string `json:"discount,omitempty"`
	Kat      string `json:"kat,omitempty"`
}
