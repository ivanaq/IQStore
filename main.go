package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func connectionDB() (connection *sql.DB) {

	Driver := "mysql"
	User := "root"
	Password := "123456"
	DBName := "stock_db"

	connection, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+DBName)
	if err != nil {
		panic(err.Error())

	}
	return connection
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", Init)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	log.Println("Server running...")
	http.ListenAndServe(":8080", nil)

}

type Product struct {
	SKU      string
	Name     string
	Brand    string
	Size     int
	Price    float64
	Imageurl string
}

//type Products []Product

func Init(w http.ResponseWriter, r *http.Request) {
	establishedConnection := connectionDB()

	records, err := establishedConnection.Query("SELECT * FROM productsstock")

	if err != nil {
		panic(err.Error())
	}

	product := Product{}
	arrayProduct := []Product{}

	for records.Next() {

		var sku, name, brand, imageurl string
		var size int
		var price float64

		err = records.Scan(&sku, &name, &brand, &size, &price, &imageurl)
		if err != nil {

			panic(err.Error())
		}

		product.SKU = sku
		product.Name = name
		product.Brand = brand
		product.Size = size
		product.Price = price
		product.Imageurl = imageurl

		arrayProduct = append(arrayProduct, product)
	}
	fmt.Println(arrayProduct)
	templates.ExecuteTemplate(w, "start", arrayProduct)
}

func Edit(w http.ResponseWriter, r *http.Request) {

	skuProduct := r.URL.Query().Get("SKU")
	fmt.Println(skuProduct)
	establishedConnection := connectionDB()
	record, err := establishedConnection.Query("SELECT * FROM productsstock WHERE SKU=?", skuProduct)

	product := Product{}

	for record.Next() {

		var sku, name, brand, imageurl string
		var size int
		var price float64

		err = record.Scan(&sku, &name, &brand, &size, &price, &imageurl)
		if err != nil {

			panic(err.Error())
		}

		product.SKU = sku
		product.Name = name
		product.Brand = brand
		product.Size = size
		product.Price = price
		product.Imageurl = imageurl

	}
	fmt.Println(product)
	templates.ExecuteTemplate(w, "edit", product)
}

func Create(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		sku := r.FormValue("sku")
		name := r.FormValue("name")
		brand := r.FormValue("brand")
		size := r.FormValue("size")
		price := r.FormValue("price")
		imageurl := r.FormValue("imageurl")

		establishedConnection := connectionDB()
		insertReg, err := establishedConnection.Prepare("INSERT INTO productsstock VALUES(?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertReg.Exec(sku, name, brand, size, price, imageurl)
		http.Redirect(w, r, "/", 301)
	}

	templates.ExecuteTemplate(w, "insert", nil)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		sku := r.FormValue("sku")
		name := r.FormValue("name")
		brand := r.FormValue("brand")
		size := r.FormValue("size")
		price := r.FormValue("price")
		imageurl := r.FormValue("imageurl")

		establishedConnection := connectionDB()

		UpdateReg, err := establishedConnection.Prepare("UPDATE productsstock SET name=?, brand=?, size=?, price=?, imageurl=? WHERE sku=?")

		if err != nil {
			panic(err.Error())
		}

		UpdateReg.Exec(name, brand, size, price, imageurl, sku)
		http.Redirect(w, r, "/", 301)
	}

	templates.ExecuteTemplate(w, "insert", nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	skuProduct := r.URL.Query().Get("SKU")
	fmt.Println(skuProduct)

	establishedConnection := connectionDB()
	deleteReg, err := establishedConnection.Prepare("DELETE FROM productsstock WHERE SKU=?")
	if err != nil {
		panic(err.Error())
	}
	deleteReg.Exec(skuProduct)
	http.Redirect(w, r, "/", 301)
}
