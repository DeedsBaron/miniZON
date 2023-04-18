package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"route256/libs/logger"
	"route256/loms/internal/models"
	"route256/loms/internal/repository/schema"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *LomsRepository) GetStocks(ctx context.Context, sku uint32) (*models.Stocks, error) {
	db := r.GetQueryEngine(ctx)
	childSpan, _ := opentracing.StartSpanFromContext(ctx, "repo: get_stocks")
	defer childSpan.Finish()

	sq := squirrel.Expr(
		`select s.item_sku, s.warehouse_id, (count - coalesce(t1.reservered_count, 0)) as count
		from stocks s
		left join (select warehouse_id, o.status, sum(reserved_count) as reservered_count
					from reservation
					join orders o on reservation.order_id = o.id
					where item_sku = $1 and o.status = $2
					group by warehouse_id, o.status) as t1
		on s.warehouse_id = t1.warehouse_id
		where s.item_sku = $1 and count - coalesce(t1.reservered_count, 0) != 0;`, sku, models.StatusAwaitingPayment)

	rawQuery, args, err := sq.ToSql()
	if err != nil {
		logger.Errorf("tosql error %v", err)
		return nil, models.ErrInternal
	}
	var st []schema.Stock
	if err = pgxscan.Select(ctx, db, &st, rawQuery, args...); err != nil {
		logger.Errorf("stocks database error sku=%d, err:%v\n", sku, err)
		return nil, models.ErrInternal
	}
	if st == nil {
		return nil, models.ErrStocksNotFound
	}

	stocks := &models.Stocks{
		Stocks: make([]models.StocksItem, 0, len(st)),
	}
	for _, stock := range st {
		stocks.Stocks = append(stocks.Stocks, models.StocksItem{
			WarehouseID: models.WarehouseID(stock.WarehouseId),
			Count:       stock.Count,
		})
	}

	return stocks, nil
}
