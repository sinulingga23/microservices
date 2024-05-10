package payload

type ResponseGeneral struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func NewResponseGeneral(statusCode int, message string, data ...interface{}) ResponseGeneral {
	response := ResponseGeneral{StatusCode: statusCode, Message: message}
	if len(data) > 0 {
		response.Data = data[0]
	}
	return response
}
