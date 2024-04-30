package main

import "net/http"

func main() {
	// mongodb+srv://challengeUser:WUMgIwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true

	// server
	sm := http.NewServeMux()
	http.ListenAndServe(":3000", sm)

	// routes
	// /GET records from Mongo
	// /GET from map
	// /POST to map

}
