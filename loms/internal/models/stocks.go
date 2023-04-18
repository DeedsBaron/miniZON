package models

type StocksItem struct {
	WarehouseID WarehouseID
	Count       uint64
}

type Stocks struct {
	Stocks []StocksItem
}

type ItemStocks map[WarehouseID]uint64
