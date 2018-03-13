# myMux



NOT FOR PRODUCTION USE - WORK IN PROGRESS

Simple http router for go


	handler := mymux.NewRouterTemplateHandler(myErrorHandler)

Custom path types

	handler.RegisterType("caps", "[A-Z]{1,64}")
	
Add routes
	
	handler.AppendRoute("GET", "/echo/Text:caps", mymux.Adapt(caps))
    
Define handlers
    
    func caps(w http.ResponseWriter, r *http.Request) {
        vars := mymux.GetVars(r)
        text := vars["Text"]
        w.Write([]byte(text))
    }
    
    func myErrorHandler(w http.ResponseWriter, status int, detail string) {
        w.WriteHeader(status)
        w.Write([]byte(detail))
    }
