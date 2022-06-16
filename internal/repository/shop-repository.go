package repository

import (
	"time"

	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"gorm.io/gorm"
)

type ShopRepository interface {
	InsertProduct(product entity.Product) error
	DeleteProduct(productID int) error
	UpdateProduct(product entity.Product, userID int) error
	ListBuys() ([]*entity.Buy, error)
	ListProducts() ([]*entity.Product, error)
	ProductDetail(productID int) (*entity.Product, error)
	GetProduct(name string) (*entity.Product, error)
	DeleteProductByName(id int, name string) error
	AddProductCart(idcart int, product dto.ProductToCart) error
	StockProduct(name string) (bool, error)
	// RemoveProductCart(idCart int, product dto.ProductToCart) error
	AddCreditCard(cc string, idcart int) error
	DeleteCart(idcart int) error
	CartExist(idcart int) (bool, error)
	ConfirmPayment(idUser interface{}) error
}

type shopConnection struct {
	connection *gorm.DB
}

func NewShopRepository(dbConn *gorm.DB) ShopRepository {
	return &shopConnection{
		connection: dbConn}
}

func (db *shopConnection) InsertProduct(p entity.Product) error {
	if err := db.connection.Raw("INSERT INTO products(detail, price, createdat, updatedat) VALUES(?,?,?,?);", p.Detail, p.Price, time.Now(), time.Now()).Scan(p).Error; err != nil {
		return err
	}

	return nil
}

func (db *shopConnection) DeleteProduct(productID int) error {
	p := entity.Product{}
	if err := db.connection.Raw("DELETE FROM products WHERE id = ?", productID).Scan(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *shopConnection) DeleteProductByName(id int, name string) error {
	p := entity.Product{}
	if err := db.connection.Raw("DELETE FROM products WHERE detail = ? AND id = ?", name, id).Scan(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *shopConnection) UpdateProduct(product entity.Product, productID int) error {
	p := entity.Product{}
	if err := db.connection.Raw("UPDATE products SET detail = ? , price = ? , updatedat = ? WHERE id = ?", product.Detail, product.Price, time.Now(), productID).Scan(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *shopConnection) ListBuys() ([]*entity.Buy, error) {
	buys := []*entity.Buy{}
	if err := db.connection.Raw("SELECT * FROM buys ORDER BY createdat DESC").Scan(&buys).Error; err != nil {
		return nil, err
	}
	return buys, nil
}

func (db *shopConnection) ListProducts() ([]*entity.Product, error) {
	p := []*entity.Product{}
	if err := db.connection.Raw("SELECT * FROM products").Scan(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (db *shopConnection) ProductDetail(productID int) (*entity.Product, error) {
	detail := entity.Product{}
	if err := db.connection.Select("detail").Where("id = ?", productID).Find(&detail).Error; err != nil {
		return nil, err
	}

	return &detail, nil
}

func (db *shopConnection) GetProduct(name string) (*entity.Product, error) {
	product := entity.Product{}
	if err := db.connection.Raw("SELECT * FROM products WHERE detail = ?", name).Scan(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (db *shopConnection) AddProductCart(idcart int, product dto.ProductToCart) error {
	cart := entity.Carts{}
	if err := db.connection.Raw("INSERT INTO carts(idcart,product,creditcard) VALUES(?,?,?);", idcart, product, "").Scan(&cart).Error; err != nil {
		return err
	}

	return nil
}

func (db *shopConnection) StockProduct(name string) (bool, error) {
	products := entity.Product{}
	var err error
	if err = db.connection.Select("id").Where("detail = ?", name).Find(&products).Error; err != nil {
		return false, err
	}
	if products.ID <= 0 {
		return false, err
	}
	return true, nil
}

func (db *shopConnection) CartExist(idcart int) (bool, error) {
	products := entity.Carts{}
	var err error
	if err = db.connection.Select("product").Where("idcart = ?", idcart).Find(&products).Error; err != nil {
		return false, err
	}
	if products.Product == nil {
		return false, err
	}
	return true, nil
}

func (db *shopConnection) DeleteProductFromCart(id int, name string) error {
	p := entity.Carts{}
	if err := db.connection.Raw("DELETE FROM carts WHERE detail = ? AND id = ?", name, id).Scan(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *shopConnection) AddCreditCard(cc string, idcart int) error {
	p := entity.Carts{}
	if err := db.connection.Raw("UPDATE carts SET creditcard = ? WHERE idcart = ?", cc, idcart).Scan(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *shopConnection) DeleteCart(idcart int) error {
	p := entity.Carts{}
	if err := db.connection.Raw("DELETE FROM carts WHERE idcart = ?", idcart).Scan(p).Error; err != nil {
		return err
	}
	return nil
}

func (db *shopConnection) ConfirmPayment(idUser interface{}) error {
	buys := entity.Buy{}
	if err := db.connection.Raw("INSERT INTO buys(iduser,createdat) VALUES(?,?);", idUser, time.Now()).Scan(&buys).Error; err != nil {
		return err
	}

	return nil
}
