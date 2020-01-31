package approval

import (
	"flag"
)

var (
	listenPortFlag       = flag.Int("ListenPort", 8111, "Port to Listen on")
	requestNumberFlag    = flag.String("RequestNumber", "REQ-1111", "Request #")
	requestDetailsFlag   = flag.String("RequestrDetails", "Sample Request", "Request Details")
	approvalTemplateFlag = flag.String("ApprovalTemplate", "approval.html", "Path to approval template")
)

/*
// this is relative to the current directory, is there a way to embed this?
	tmpl := template.Must(template.ParseFiles(*approvalTemplateFlag))

	// TODO: inject the request #, request description into the template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// dump request info
		requestDump, err := httputil.DumpRequest(r, true)
		logging.Panic(err)

		//fmt.Println(string(requestDump))
		log.Info().Msgf("request:\n%v", string(requestDump))

		if r.Method != http.MethodPost {

		}


	})

	logging.Panic(http.ListenAndServe(fmt.Sprintf(":%v", *listenPortFlag), nil))
*/

func showRequest() {
	// no template variables
	///logging.Panic(tmpl.Execute(w, nil))
	logging.Panic(tmpl.Execute(w, struct {
		RequestNumber  string
		RequestDetails string
		Success        bool
	}{*requestNumberFlag, *requestDetailsFlag, false}))

	return
}

func processResponse() {
	log.Info().Msgf("Action: %v", r.FormValue("action"))
	log.Info().Msgf("Comments: %v", r.FormValue("comments"))

	write(*requestNumberFlag, r.FormValue("action"), r.FormValue("comments"), "nil/TBD")

	logging.Panic(tmpl.Execute(w, struct{ Success bool }{true}))
}
