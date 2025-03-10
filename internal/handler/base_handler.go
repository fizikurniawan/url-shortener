package handler

type PageData struct {
	BaseURL     string
	CurrentTime string
	User        interface{}
	Error       string
	Template    string                 // nama template yang akan digunakan
	Data        map[string]interface{} // data tambahan untuk template
}

func NewPageData(template string) *PageData {
	return &PageData{
		Template: template,
		Data:     make(map[string]interface{}),
	}
}
