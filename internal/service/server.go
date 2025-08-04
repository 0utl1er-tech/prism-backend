package service

import (
	bookv1 "github.com/0utl1er-tech/prism-backend/gen/pb/book/v1"
	contactv1 "github.com/0utl1er-tech/prism-backend/gen/pb/contact/v1"
	customerv1 "github.com/0utl1er-tech/prism-backend/gen/pb/customer/v1"
	db "github.com/0utl1er-tech/prism-backend/gen/sqlc"
)

type Server struct {
	customerv1.UnimplementedCustomerServiceServer
	bookv1.UnimplementedBookServiceServer
	contactv1.UnimplementedContactServiceServer
	queries db.Queries
}

func NewServer(queries db.Queries) *Server {
	return &Server{queries: queries}
}
