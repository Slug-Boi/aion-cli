package html

import (
	"html/template"
	"log"
	"net/http"
)

func CreateHTMLTemplate() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		const templ = `
      <!DOCTYPE html>
      <html>
        <head>
          <meta charset="UTF-8">
          <title>{{.Title}}</title>
        </head>
        <body>
          {{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
        </body>
      </html>`

		t, err := template.New("webpage").Parse(templ)
		if err != nil {
			//TODO: change to Zap logger later
			log.Fatal(err)
		}

		data := struct {
			Title string
			Items []string
		}{
			Title: "My page",
			Items: []string{
				"My photos",
				"My blog",
			},
		}

		err = t.Execute(w, data)
	})
	http.ListenAndServe(":80", nil)

}
