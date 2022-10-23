package main

import (
	"fmt"

	"github.com/ttacon/libphonenumber"
)

func main() {
	fmt.Printf("Hello world\n")

	// Parse the number.
	testNumber := "+18663472423"
	num, err := libphonenumber.Parse(testNumber, "US")
	if err != nil {
		// Handle error appropriately.
	}

	// Get the cleaned number and the length of the area code.
	natSigNumber := libphonenumber.GetNationalSignificantNumber(num)
	geoCodeLength := libphonenumber.GetLengthOfGeographicalAreaCode(num)
	fmt.Printf("test number: %s, natSigNumber: %s, geoCodeLength: %d, country_code: %d\n", testNumber, natSigNumber, geoCodeLength, *num.CountryCode)

	// Extract the area code.
	areaCode := ""
	if geoCodeLength > 0 {
		areaCode = natSigNumber[0:geoCodeLength]
	}
	fmt.Println(areaCode)

	tt, err := libphonenumber.Parse("+16502530000", "US")
	fmt.Printf("res: %v\n", tt)

	tt, err = libphonenumber.Parse("+8210216565211122112231", "US")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Printf("res: %v\n", tt)

	// libphonenumber.

	// number := "+358401231234"
	// defaultRegion := "FI"
	// possibleNumber := phonenumber.IsPossibleNumber(number, defaultRegion)
	// fmt.Printf("The number %s is possible phone number in region %s: %v\n", number, defaultRegion, possibleNumber)

	// fmt.Printf("+82: %s\n", phonenumber.CountryCode("81"))

	// number := "+358401231234"
	// defaultRegion := "FI"
	// possibleNumber := phonenumber.IsPossibleNumber(number, defaultRegion)
	// fmt.Printf("The number %s is possible phone number in region %s: %v", number, defaultRegion, possibleNumber)

}
