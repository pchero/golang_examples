package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

var XMLdata = `<urlset>
<url>
   <loc>http://www.example.com/xml-element-golang</loc>
   <lastmod>2015-06-14</lastmod>
   <changefreq>daily</changefreq>
   <priority>0.5</priority>
</url>
<url>
   <loc>http://www.example.com/extract-element-by-for-loop</loc>
   <lastmod>2015-06-14</lastmod>
   <changefreq>daily</changefreq>
   <priority>0.5</priority>
</url>
<url>
   <loc>http://www.example.com/extract-element-by-for-test</loc>
   <lastmod>2015-06-14</lastmod>
   <changefreq>dailyamp&amp;</changefreq>
   <priority>0.5</priority>
</url>
<urlset>`

type XMLQuery struct {
	Loc string `xml:",chardata"`
}

var l XMLQuery

func main() {
	// example on handling XML chardata(string)
	decoder := xml.NewDecoder(strings.NewReader(string(XMLdata)))

	for {
		// err is ignore here. If you reading from a XML file
		// do not ignore err and also check for io.EOF
		token, err := decoder.Token()
		if err != nil {
			fmt.Println(err)
		}

		if token == nil {
			break
		}

		switch Element := token.(type) {
		case xml.StartElement:
			fmt.Println("Element name: ", Element.Name.Local)

			// if Element.Name.Local == "loc" {
			// 	fmt.Println("Element name is : ", Element.Name.Local)

			// 	err := decoder.DecodeElement(&l, &Element)
			// 	if err != nil {
			// 		fmt.Println(err)
			// 	}
			// 	fmt.Println("Element value is : ", l.Loc)
			// }

		case xml.CharData:
			s := string(Element)
			// s := string([]byte(Element))
			// if s != html.UnescapeString(s) {
			// 	fmt.Println("Something was wrong: ", s)
			// }

			// print out the element data
			// convert to []byte slice and cast to string type
			fmt.Println("Data: ", s)

			// s = string([]byte("test&"))
			// fmt.Println(s)
		}

	}
}
