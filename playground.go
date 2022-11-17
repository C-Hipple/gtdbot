package main

import (
	"encoding/json"
	"fmt"
)

type PoliticsJson struct {
	Results []struct {
		Multimedia []struct {
			URL string `json:"url"`
		} `json:"multimedia"`
		Title string `json:"title"`
	} `json:"results"`
}

func main_dep() {
	s := `{
"results": [
    {
      "section": "N.Y. / Region",
      "subsection": "",
      "title": "Christie Spins His Version of Security Record on Trail",
      "abstract": "An examination of Gov. Chris Christie’s record as New Jersey’s top federal prosecutor shows that he has, at times, overstated the significance of the terrorism prosecutions he oversaw.",
      "url": "http://www.nytimes.com/2015/12/27/nyregion/Christie-markets-himself-as-protector-to-gain-in-polls.html",
      "byline": "By ALEXANDER BURNS and CHARLIE SAVAGE",
      "item_type": "Article",
      "updated_date": "2015-12-26T18:04:19-5:00",
      
      "multimedia": [
        {
          "url": "http://static01.nyt.com/images/2015/12/27/nyregion/27CHRISTIE1/27CHRISTIE1-thumbStandard.jpg",
          "format": "Standard Thumbnail",
          "height": 75,
          "width": 75,
          "type": "image",
          "subtype": "photo",
          "caption": "Gov. Chris Christie of New Jersey spoke about the Sept. 11, 2001, attacks at a Republican conference last month.",
          "copyright": "Stephen Crowley/The New York Times"
        }
        
      ]
    }
  ]
}`

	var p PoliticsJson

	err := json.Unmarshal([]byte(s), &p)
	if err != nil {
		panic(err)
	}
	fmt.Println(p.Results[0].Title)
	fmt.Println(p.Results[0].Multimedia[0].URL)
	fmt.Println(p.Results[0])
}
