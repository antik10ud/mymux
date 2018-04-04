# myMux



NOT FOR PRODUCTION USE - WORK IN PROGRESS

Simple http router for go


	handler := mymux.NewRouterTemplateHandler(myErrorHandler)

Custom path types

	handler.RegisterType("caps", "[A-Z]{1,64}")
	
Add routes
	
	echoRoute = handler.AppendRoute("GET", "/echo/Text:caps", mymux.Adapt(caps))
    
Define handlers
    
    func caps(w http.ResponseWriter, r *http.Request) {
        //get route vars
        vars := mymux.GetVars(r)
        text := vars["Text"]
    	
        //build resource urls with params
        url := echoRoute.URL(mymux.URLVars{"Text": "OTHER"})

        //return content
    	w.Write([]byte(fmt.Sprintf("You say %s, you can try %s ", text, url)))
    }
    
    func myErrorHandler(w http.ResponseWriter, status int, detail string) {
        w.WriteHeader(status)
        w.Write([]byte(detail))
    }
