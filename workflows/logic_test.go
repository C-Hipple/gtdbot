package workflows

import (
	"gtdbot/org"
	"testing"
)

func makeSection() org.Section{
	var items []org.OrgTODO
	for _, val := range []string{"1", "2", "3"} {
		items = append(items, makeItem(val))
	}

	section := org.Section{
		Name: "Section Desc",
		StartLine: 10,
		IndentLevel: 2,
		Items: items,
	}
	return section

}

func makeItem(i string) org.OrgTODO {
	return org.NewOrgItem("header "+i,
		[]string{"detail 1"+i, "detail 2"+i, "detail 3"+i},
		"TODO",
		[]string{"tag1", "tag2"},
	)
}

func Test_CheckTODOInSection(t *testing.T) {
	section := makeSection()

	at_line := CheckTODOInSection(makeItem("1"), section)
	if at_line != 11 {
		t.Fatalf("Incorrect starting line found.  Expected %v, found %v", 11, at_line)
	}

	at_line = CheckTODOInSection(makeItem("2"), section)
	if at_line != 15 {
		t.Fatalf("Incorrect starting line found.  Expected %v, found %v", 15, at_line)
	}

	at_line = CheckTODOInSection(makeItem("4"), section)
	if at_line != -1 {
		t.Fatalf("Incorrect starting line found.  Expected %v, found %v", -1, at_line)
	}
}
