package html

import (
	"html/template"
	"log"
	"net/http"
)

type WebData struct {
	GroupNumber string
	Timeslot    string
	Day         string
	Date        string
	WishLevel   string
	//Path        []int
}

// HTML template code inspired by https://gowebexamples.com/templates/
func GenerateHTML(input []WebData) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		const templ = `
      <!DOCTYPE html>
      <html>
        <head>
          <meta charset="UTF-8">
          <title>Aion Timeslot Management</title>
        </head>
		<body>
		<h1>Timeslot Management</h1>
		 <button>Re-Solve</button> 
		 <button>Advanced View</button>
		 <button>Download PDF</button>
        <ol>
          {{range .WebData}}<li>{{.GroupNumber}} {{.Timeslot}} {{.Day}} {{.Date}} {{.WishLevel}}</li>{{end}}
        </ol>
		<script>
			function downloadPDF() { 

		}
		</script>
		</body>
      </html>`

		t, err := template.New("webpage").Parse(templ)
		if err != nil {
			//TODO: change to Zap logger later
			log.Fatal(err)
		}

		data := struct {
			WebData []WebData
		}{
			WebData: input,
		}

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		// TODO: Use zaplogger to log the error.
		return
	}

}
