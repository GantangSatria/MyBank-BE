package response

type Base struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

//helpers

func Success(message string, data interface{}) *Base {
	return &Base{Success: true, Message: message, Data: data}
}

func SuccessWithMeta(message string, data interface{}, meta *Meta) *Base {
	return &Base{Success: true, Message: message, Data: data, Meta: meta}
}

func Error(message string) *Base {
	return &Base{Success: false, Message: message}
}
