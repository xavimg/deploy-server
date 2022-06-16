package service

import (
	"errors"

	"github.com/xavimg/Turing/apituringserver/internal/dto"
	"github.com/xavimg/Turing/apituringserver/internal/entity"
	"github.com/xavimg/Turing/apituringserver/internal/repository"
)

type ShopService interface {
	InsertProduct(product entity.Product) error
	DeleteProduct(productID int) error
	UpdateProduct(product entity.Product, productID int) error
	ListBuys() ([]*entity.Buy, error)
	ListProducts() ([]*entity.Product, error)
	ProductDetail(productID int) (*entity.Product, error)
	GetProduct(name string) (*dto.ProductToCart, error)
	DeleteProductByName(id int, name string) error
	AddProductCart(idCart int, product dto.ProductToCart) error
	StockProduct(name string) (bool, error)
	AddCreditCard(cc string, idcart int) error
	DeleteCart(idcart int) error
	CartExist(idcart int) (bool, error)
	ConfirmPayment(idUser interface{}) error
}

type shopService struct {
	ShopRepository repository.ShopRepository
}

func NewShopService(ss repository.ShopRepository) ShopService {
	return &shopService{
		ShopRepository: ss,
	}
}

func (ss *shopService) InsertProduct(p entity.Product) error {
	if err := ss.ShopRepository.InsertProduct(p); err != nil {
		return err
	}
	return nil
}

func (ss *shopService) DeleteProduct(productID int) error {
	if err := ss.ShopRepository.DeleteProduct(productID); err != nil {
		return err
	}
	return nil
}

func (ss *shopService) DeleteProductByName(id int, name string) error {
	if err := ss.ShopRepository.DeleteProductByName(id, name); err != nil {
		return err
	}
	return nil
}

func (ss *shopService) UpdateProduct(product entity.Product, productID int) error {
	if err := ss.ShopRepository.UpdateProduct(product, productID); err != nil {
		return err
	}
	return nil
}

func (ss *shopService) ListBuys() ([]*entity.Buy, error) {
	buys, err := ss.ShopRepository.ListBuys()
	if err != nil {
		return nil, err
	}
	return buys, nil
}

func (ss *shopService) ListProducts() ([]*entity.Product, error) {
	products, err := ss.ShopRepository.ListProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ss *shopService) ProductDetail(productID int) (*entity.Product, error) {
	detail, err := ss.ShopRepository.ProductDetail(productID)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (ss *shopService) GetProduct(name string) (*dto.ProductToCart, error) {
	product, err := ss.ShopRepository.GetProduct(name)
	if err != nil {
		return nil, err
	}

	pdc := dto.ProductToCart{
		ID:       int(product.ID),
		Name:     product.Detail,
		Price:    product.Price,
		Quantity: 1,
	}

	return &pdc, nil
}

func (ss *shopService) AddProductCart(idCart int, product dto.ProductToCart) error {
	if err := ss.ShopRepository.AddProductCart(idCart, product); err != nil {
		return err
	}

	return nil
}

func (s *shopService) StockProduct(name string) (bool, error) {
	exists, err := s.ShopRepository.StockProduct(name)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("no more stock of this product")
	}

	return true, nil
}

func (s *shopService) CartExist(idcart int) (bool, error) {
	exists, err := s.ShopRepository.CartExist(idcart)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("cart doesn't exists")
	}

	return true, nil
}

func (ss *shopService) AddCreditCard(cc string, idcart int) error {
	if err := ss.ShopRepository.AddCreditCard(cc, idcart); err != nil {
		return err
	}
	return nil
}

func (ss *shopService) DeleteCart(idcart int) error {
	if err := ss.ShopRepository.DeleteCart(idcart); err != nil {
		return err
	}
	return nil
}

func (ss *shopService) ConfirmPayment(idUser interface{}) error {
	if err := ss.ShopRepository.ConfirmPayment(idUser); err != nil {
		return err
	}
	return nil
}
