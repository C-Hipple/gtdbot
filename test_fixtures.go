package main

import "strconv"

func Fixture_LeankitCardOrgLineItem(num int) LeankitCardOrgLineItem {
	str_value_of_num := strconv.FormatInt(int64(num), 10)
	return LeankitCardOrgLineItem{"Title" + str_value_of_num, "TODO", "abc.com/" + str_value_of_num}
}

func Fixture_Section() Section {
	items := make([]LeankitCardOrgLineItem, 4)
	for i := range [4]int{} {
		items = append(items, Fixture_LeankitCardOrgLineItem(i))
	}

	section := Section{Description: "Cards", StartLine: 5, IndentLevel: 2}
	for _, item := range items {
		section.Items = append(section.Items, item)
	}
	return section

}

func Fixture_Card() Card {
	return Card{Id: "123", Title: "Title"}
}
