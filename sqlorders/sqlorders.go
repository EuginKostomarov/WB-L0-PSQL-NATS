package sqlorders

import (
	"database/sql"
	"fmt"
	"model"

	_ "github.com/lib/pq"
)

type Repository interface {
	FromDb()
	FindById(order_id string) (model.Order, error)
	Create(order *model.Order) (*model.Order, error)
}

type orderRepository struct {
	All map[string]model.Order
}

func New() *orderRepository {
	return &orderRepository{
		All: make(map[string]model.Order),
	}
}

func (r *orderRepository) FromDb() {
	connStr := "user=hirazant password=1112 dbname=hirazant sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	rowsOrders, err := db.Query("select  * from orders")
	if err != nil {
		panic(err)
	}
	defer rowsOrders.Close()

	for rowsOrders.Next() {
		order := model.Order{}
		rowsOrders.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)

		rowsPayment, err := db.Query("select transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee from payments where fk_payments_order=$1", &order.OrderUid)
		if err != nil {
			panic(err)
		}
		defer rowsPayment.Close()
		payment := model.Payment{}
		for rowsPayment.Next() {
			err = rowsPayment.Scan(&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider,
				&payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
			if err != nil {
				panic(err)
			}
		}

		rowsDelivery, err := db.Query("select name, phone, zip, city, address, region, email from delivery where fk_delivery_order=$1", &order.OrderUid)
		if err != nil {
			panic(err)
		}
		defer rowsDelivery.Close()
		delivery := model.Delivery{}
		for rowsDelivery.Next() {
			err = rowsDelivery.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
				&delivery.Region, &delivery.Email)
			if err != nil {
				panic(err)
			}
		}

		rowsItems, err := db.Query("select chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status from items where fk_items_order=$1", &order.OrderUid)
		if err != nil {
			panic(err)
		}
		defer rowsItems.Close()

		var items []model.Item
		for rowsItems.Next() {
			item := model.Item{}
			err := rowsItems.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
			if err != nil {
				fmt.Println(err)
				continue
			}
			items = append(items, item)
		}

		t := model.Order{order.OrderUid, order.TrackNumber, order.Entry, delivery,
			payment, items, order.Locale, order.InternalSignature,
			order.CustomerId, order.DeliveryService, order.Shardkey,
			order.SmId, order.DateCreated, order.OofShard}
		r.All[order.OrderUid] = t
	}
}

func (r *orderRepository) FindById(order_id string) (model.Order, error) {
	order := r.All[order_id]
	return order, nil
}

func (r *orderRepository) Create(order *model.Order) (*model.Order, error) {
	connStr := "user=hirazant password=1112 dbname=hirazant sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db.QueryRow(`insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	db.QueryRow(`insert into payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, fk_payments_order) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee, order.OrderUid)
	db.QueryRow(`insert into delivery (name, phone, zip, city, address, region, email, fk_delivery_order) values ($1,$2,$3,$4,$5,$6,$7,$8)`,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email, order.OrderUid)
	for _, item := range order.Items {
		db.QueryRow(`insert into items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, fk_items_order) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId,
			item.Brand, item.Status, order.OrderUid)
	}
	r.All[order.OrderUid] = *order
	return order, nil
}
