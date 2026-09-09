package main

import "fmt"

type Product struct {
	Name  string
	Price float64
	SKU   string
}

func (p *Product) CalculateShipping(weight float64) float64 {
	return weight * 2.5
}

func (p *Product) GetInfo() string {
	return fmt.Sprintf("%s - $%.2f (SKU: %s)", p.Name, p.Price, p.SKU)
}

type DigitalProduct struct {
	Product
	DownloadSize float64
}

func (d *DigitalProduct) CalculateShipping(weight float64) float64 {
	return 0
}

func (d *DigitalProduct) GetInfo() string {
	return fmt.Sprintf("%s - Digital Download (%.1f MB)", d.Product.GetInfo(), d.DownloadSize)
}

type PhysicalProduct struct {
	Product
	Weight float64
}

func (p *PhysicalProduct) CalculateShipping(weight float64) float64 {
	w := weight
	if w <= 0 {
		w = p.Weight
	}
	return p.Product.CalculateShipping(w)
}

type FragileProduct struct {
	PhysicalProduct
}

func (f *FragileProduct) CalculateShipping(weight float64) float64 {
	w := weight
	if w <= 0 {
		w = f.Weight
	}
	base := f.PhysicalProduct.Product.CalculateShipping(w)
	return base + 10.0
}

func (f *FragileProduct) GetInfo() string {
	return fmt.Sprintf("%s - FRAGILE", f.PhysicalProduct.GetInfo())
}

func main() {
	ebook := &DigitalProduct{Product: Product{Name: "Python Guide", Price: 29.99, SKU: "DIG-001"}, DownloadSize: 5.2}
	fmt.Println(ebook.GetInfo())
	fmt.Printf("Shipping: $%.2f\n", ebook.CalculateShipping(0))

	book := &PhysicalProduct{Product: Product{Name: "Python Book", Price: 39.99, SKU: "PHY-001"}, Weight: 0.5}
	fmt.Println(book.GetInfo())
	fmt.Printf("Shipping: $%.2f\n", book.CalculateShipping(0))

	vase := &FragileProduct{PhysicalProduct: PhysicalProduct{Product: Product{Name: "Ceramic Vase", Price: 49.99, SKU: "FRG-001"}, Weight: 1.2}}
	fmt.Println(vase.GetInfo())
	fmt.Printf("Shipping: $%.2f\n", vase.CalculateShipping(0))
}
