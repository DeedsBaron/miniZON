package schema

type CartItem struct {
	ItemSku uint32 `db:"item_sku"`
	Count   uint64 `db:"count"`
}
