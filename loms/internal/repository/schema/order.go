package schema

type Order struct {
	Id int64 `db:"id"`
}

type OrderItems struct {
	Sku   uint32 `db:"item_sku"`
	Count uint16 `db:"reserved_count"`
}

type ListOrder struct {
	ItemSku       uint32 `db:"item_sku"`
	ReservedCount uint32 `db:"reserved_count"`
}

type Id int
