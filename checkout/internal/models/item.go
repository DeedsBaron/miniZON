package models

type Sku uint32

type SkuInfo struct {
	Sku       Sku
	CartIndex int
}

type Count uint32

type Item struct {
	Sku   Sku
	Count Count
	Name  string
	Price uint32
}

type ProductInfo struct {
	Name      string
	Price     uint32
	CartIndex int
}
