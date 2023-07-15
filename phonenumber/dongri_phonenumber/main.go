package main

import (
	"fmt"

	"github.com/dongri/phonenumber"
)

func main() {
	number := phonenumber.Parse("090-6135-4467", "JP")
	fmt.Println(number)
	// Output: 819061354467

	number2 := phonenumber.Parse("+821121656521", "")
	fmt.Println(number2)

	country := phonenumber.GetISO3166ByNumber(number2, false)
	fmt.Printf("alpha2: %s, alpha3: %s, country_cdoe: %s, country_name: %s, mobile_begin_with: %v, lenth: %v\n", country.Alpha2, country.Alpha3, country.CountryCode, country.CountryName, country.MobileBeginWith, country.PhoneNumberLengths)
	fmt.Println(country.CountryName)
}
