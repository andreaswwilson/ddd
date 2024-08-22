package mock

import "ddd"

var JiraFormInput = map[string]string{
	"online-1": `[
   {
      "label":"Kostnadsoppfølger",
      "questionKey":"kostnadsoppfolger",
      "answer":"Angrboða@skatteetaten.no"
   }
]`,
}

var JiraForm = map[string]*ddd.JiraForm{
	"online-1": {
		Kostnadsoppfolger: "Angrboða@skatteetaten.no",
	},
}
