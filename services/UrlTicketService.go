package services

type UrlTicketService interface {
	TicketService
	GetUrl(string) string
	Incr(string)
}

type urlTicketServiceImpl struct {
	ticketServiceImpl
}

func UrlTicketServiceOf() UrlTicketService  {
	var service = new(urlTicketServiceImpl)
	service.Init()
	return service
}

func (this *urlTicketServiceImpl)GetUrl(ticket string)string {
	var data ,err =this.GetTicketInfo(ticket)
	if err!=nil {
		return ""
	}
	return data["url"].(string)
}

func (this *urlTicketServiceImpl)Incr(ticket string)  {

}