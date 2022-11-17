package main

import (
	"fmt"
)

func main() {
	fmt.Println("Org Code Review")
	// PrintOrgFile(GetOrgFile())
	items := ParseCodeReview(GetOrgFile())
	for _, item := range items {
		fmt.Println(item)
	}

	fmt.Println("Cards!")

	main_leankit()

}
