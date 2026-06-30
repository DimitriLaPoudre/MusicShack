package response

type StatusResponse struct {
	Status string `json:"status"`
}

var Ok = StatusResponse{
	Status: "ok",
}
