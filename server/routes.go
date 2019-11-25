package server

// SetupRoutes configures all the routes for this service
func (s *Server) SetupRoutes() {

	s.router.Get("/version", s.GetVersion())

	s.router.Post("/invoices", s.CreateInvoice())

	s.router.Get("/invoices/{id}", s.GetInvoice())

}
