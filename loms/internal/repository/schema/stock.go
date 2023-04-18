package schema

type Stock struct {
	ItemSku     uint32 `db:"item_sku"`
	WarehouseId int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}
