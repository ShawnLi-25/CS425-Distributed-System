package main

import (
	"fmt"
	"testing"
	"strconv"
)

func TestLogCnt(t *testing.T){
	frequentPtrn := "-E a"

	extendedPtrn := "-E ab"

	infrequentPtrn := "-E abc"

	//unitTesting result below
	fmt.Println("Testing logCnt function")
	frequentRes := logCnt(frequentPtrn)
	
	extendedRes := logCnt(extendedPtrn)
	
	infrequentRes := logCnt(infrequentPtrn)
	
	fmt.Println("Result")
	if frequentRes == 939979 {
		fmt.Print("=== The Result of Grep for Frequent Pattern is True ===\n")

	} else{
		fmt.Print("=== The Result of Grep for Frequent Pattern is Wrong ===\n === Expected: 939979, Result:" , strconv.Itoa(frequentRes), "===\n")
	}	
	if extendedRes == 7725 {
		fmt.Print("=== The Result of Grep for Extended Pattern is True ===\n")

	} else{
		fmt.Print("=== The Result of Grep for Extended Pattern is Wrong ===\n === Expected: 7725, Result:" , strconv.Itoa(extendedRes), "===\n")
	}	
	if infrequentRes == 91 {
		fmt.Print("=== The Result of Grep for Infrequent Pattern is True ===\n")

	} else{
		fmt.Print("=== The Result of Grep for Infrequent Pattern is Wrong ===\n === Expected: 91, Result:" , strconv.Itoa(infrequentRes), "===\n")
	}	
}
