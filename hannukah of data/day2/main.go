package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type JsonCustomer struct {
	Customerid   int     `json:"customerid"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Citystatezip string  `json:"citystatezip"`
	Birthdate    string  `json:"birthdate"`
	Phone        string  `json:"phone"`
	Timezone     string  `json:"timezone"`
	Lat          float64 `json:"lat"`
	Long         float64 `json:"long"`
}
type JsonOrder struct {
	Orderid    int    `json:"orderid"`
	Customerid int    `json:"customerid"`
	Ordered    string `json:"ordered"`
	Shipped    string `json:"shipped"`
	Items      []struct {
		Sku       string  `json:"sku"`
		Qty       int     `json:"qty"`
		UnitPrice float64 `json:"unit_price"`
	} `json:"items"`
	Total float64 `json:"total"`
}
type Order struct {
	OrderObj    *JsonOrder
	CustomerObj *JsonCustomer
}

func readFileByLine(filename string) []string {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.TrimSpace(string(dat)), "\n")
	return lines
}

var phoneDict = initPhoneDict()

func initPhoneDict() map[rune]string {
	return map[rune]string{
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}
}
func isLastNameValid(lastName, phoneNumber string) bool {
	rLastName := []rune(strings.ToLower(lastName))
	if len(rLastName) != 10 {
		return false
	}
	numbers := []rune(strings.Join(strings.Split(phoneNumber, "-"), ""))
	for i, r := range numbers {
		padRunes := phoneDict[r]
		if !strings.ContainsRune(padRunes, rLastName[i]) {
			return false
		}
	}
	return true

}
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
func main() {
	customers := make(map[int]JsonCustomer, 0)
	for _, line := range readFileByLine("./noahs-jsonl/5784/noahs-customers.jsonl") {
		var customer JsonCustomer
		byt := []byte(line)
		if err := json.Unmarshal(byt, &customer); err != nil {
			panic(err)
		}
		customers[customer.Customerid] = customer
		lastName := strings.Split(customer.Name, " ")[1]
		// Part 1
		if isLastNameValid(lastName, customer.Phone) {
			println(customer.Customerid, customer.Name, customer.Address, customer.Citystatezip, customer.Birthdate, customer.Phone, customer.Timezone, customer.Lat, customer.Long)
			println("1. Phone Number:", customer.Phone, "\n\n")
		}
	}
	// Part 2
	type Product struct {
		Sku           string    `json:"sku"`
		Desc          string    `json:"desc"`
		WholesaleCost float64   `json:"wholesale_cost"`
		DimsCm        []float64 `json:"dims_cm"`
	}
	products := make(map[string]Product, 0)
	for _, line := range readFileByLine("./noahs-jsonl/5784/noahs-products.jsonl") {
		var product Product
		byt := []byte(line)
		if err := json.Unmarshal(byt, &product); err != nil {
			panic(err)
		}
		products[product.Sku] = product
	}
	orders := make(map[int]Order, 0)
	for _, line := range readFileByLine("./noahs-jsonl/5784/noahs-orders.jsonl") {
		var order JsonOrder
		byt := []byte(line)
		if err := json.Unmarshal(byt, &order); err != nil {
			panic(err)
		}
		customer, ok := customers[order.Customerid]
		if !ok {
			panic("Customer not found for order: " + string(byt))
		}
		orders[order.Orderid] = Order{
			OrderObj:    &order,
			CustomerObj: &customer,
		}
		//part 2 logic
		splitCustomerName := strings.Split(customer.Name, " ")
		firstName, lastName := []rune(splitCustomerName[0]), []rune(splitCustomerName[1])
		if firstName[0] == 'J' && lastName[0] == 'P' && strings.Split(order.Shipped, "-")[0] == "2017" {
			for _, item := range order.Items {
				if strings.Contains(strings.ToLower(products[item.Sku].Desc), "coffee") {
					println("Order ID:", order.Orderid, "Customer ID:", customer.Customerid, "Customer Name:", customer.Name, "Product SKU:", item.Sku, "Qty:", item.Qty, "Unit Price:", item.UnitPrice)
					println("2. Phone Number:", customer.Phone)
					println("Zipcode for 3.", strings.Split(customer.Citystatezip, " ")[2], "\n\n")
				}
			}
		}
	}
	//part 3 logic
	//what we know:
	//cancer birthday
	//year of the rabbit
	//lives in neighborhood
	zipcode := "11435"
	//cancer birthday is between 6/21 and 7/22
	for _, customer := range customers {
		yearOfDogBase := 1951
		customerZip := strings.Split(customer.Citystatezip, " ")[2]
		if customerZip != zipcode {
			continue
		}
		splitBirthday := strings.Split(customer.Birthdate, "-")
		birthYear, birthMonth, birthDay := parseInt(splitBirthday[0]), parseInt(splitBirthday[1]), parseInt(splitBirthday[2])
		if birthMonth > 7 || birthMonth == 7 && birthDay > 22 {
			continue
		}
		if birthMonth < 6 || birthMonth == 6 && birthDay < 21 {
			continue
		}
		for yearOfDogBase < 2025 {
			if birthYear == yearOfDogBase {
				println("Customer ID:", customer.Customerid, "Name:", customer.Name, "Address:", customer.Address, "City/State/Zip:", customer.Citystatezip, "Birthdate:", customer.Birthdate)
				println("3. Phone Number:", customer.Phone, "\n\n")
			}
			yearOfDogBase += 12
		}

	}

}
