package leankit

import "strconv"

func Fixture_LeankitCardOrgLineItem(num int) LeankitCardOrgLineItem {
	str_value_of_num := strconv.FormatInt(int64(num), 10)
	notes := []string{"Assigned User: Chris", "PR: github.com/user/repo/" + str_value_of_num}
	return LeankitCardOrgLineItem{"Title" + str_value_of_num, "TODO", "abc.com/" + str_value_of_num, notes}
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

func Fixture_Card(id string) Card {
	if id == "" {
		id = "123"
	}
	return Card{Id: id, Title: "Title"}
}
