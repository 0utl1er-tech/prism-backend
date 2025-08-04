package service

import (
	"context"

	contactv1 "github.com/0utl1er-tech/prism-backend/gen/pb/contact/v1"
	customerv1 "github.com/0utl1er-tech/prism-backend/gen/pb/customer/v1"
	db "github.com/0utl1er-tech/prism-backend/gen/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CustomerService struct {
	customerv1.UnimplementedCustomerServiceServer
	queries *db.Queries
}

func NewCustomerService(queries *db.Queries) *CustomerService {
	return &CustomerService{
		queries: queries,
	}
}

func (server *CustomerService) CreateCustomer(ctx context.Context, customer *customerv1.CreateCustomerRequest) (*customerv1.CreateCustomerResponse, error) {
	customerId := uuid.New()
	contactId := uuid.New()
	leaderId := uuid.New()
	picId := uuid.New()

	customerArg := db.CreateCustomerParams{
		ID:     customerId,
		BookID: uuid.MustParse(customer.GetBookId()),
		Name:   customer.GetName(),
		Corporation: pgtype.Text{
			String: customer.GetCorporation(),
			Valid:  customer.GetCorporation() != "",
		},
		Address: pgtype.Text{
			String: customer.GetAddress(),
			Valid:  customer.GetAddress() != "",
		},
	}

	contactArg := db.CreateContactParams{
		ID:         contactId,
		CustomerID: customerId,
		Phone:      customer.GetContact().GetPhone(),
		Mail: pgtype.Text{
			String: customer.GetContact().GetMail(),
			Valid:  customer.GetContact().GetMail() != "",
		},
	}

	leaderArg := db.CreateStaffParams{
		ID: leaderId,
		Name: pgtype.Text{
			String: customer.GetLeader(),
			Valid:  customer.GetLeader() != "",
		},
		Sex: pgtype.Text{
			String: customer.GetLeaderSex(),
			Valid:  customer.GetLeaderSex() != "",
		},
	}

	picArg := db.CreateStaffParams{
		ID: picId,
		Name: pgtype.Text{
			String: customer.GetPic(),
			Valid:  customer.GetPic() != "",
		},
		Sex: pgtype.Text{
			String: customer.GetPicSex(),
			Valid:  customer.GetPicSex() != "",
		},
	}

	// db.Queriesをポインタとして使用
	customerRes, err := server.queries.CreateCustomer(ctx, customerArg)
	if err != nil {
		return nil, err
	}

	_, err = server.queries.CreateContact(ctx, contactArg)
	if err != nil {
		return nil, err
	}

	_, err = server.queries.CreateStaff(ctx, leaderArg)
	if err != nil {
		return nil, err
	}

	_, err = server.queries.CreateStaff(ctx, picArg)
	if err != nil {
		return nil, err
	}

	// レスポンスを返す処理を追加
	return &customerv1.CreateCustomerResponse{
		Id:          customerRes.ID.String(),
		BookId:      customerRes.BookID.String(),
		Name:        customerRes.Name,
		Corporation: customerRes.Corporation.String,
		Address:     customerRes.Address.String,
		Memo:        customerRes.Memo.String,
	}, nil
}

func (server *CustomerService) SearchCustomer(ctx context.Context, customer *customerv1.SearchCustomerRequest) (*customerv1.SearchCustomerResponse, error) {
	customerArg := db.SearchCustomerParams{
		BookID: pgtype.UUID{
			Bytes: uuid.MustParse(customer.GetBookId()),
			Valid: customer.GetBookId() != "",
		},
		Name: pgtype.Text{
			String: customer.GetName(),
			Valid:  customer.GetName() != "",
		},
		Corporation: pgtype.Text{
			String: customer.GetCorporation(),
			Valid:  customer.GetCorporation() != "",
		},
		Address: pgtype.Text{
			String: customer.GetAddress(),
			Valid:  customer.GetAddress() != "",
		},
		Memo: pgtype.Text{
			String: customer.GetMemo(),
			Valid:  customer.GetMemo() != "",
		},
	}

	customers, err := server.queries.SearchCustomer(ctx, customerArg)
	if err != nil {
		return nil, err
	}

	customersRes := make([]*customerv1.Customer, len(customers))
	for i, customer := range customers {
		customersRes[i] = &customerv1.Customer{
			Id:          customer.ID.String(),
			Name:        customer.Name,
			Corporation: customer.Corporation.String,
			Address:     customer.Address.String,
			Memo:        customer.Memo.String,
		}
	}

	return &customerv1.SearchCustomerResponse{
		Customers: customersRes,
	}, nil
}

func (server *CustomerService) GetCustomer(ctx context.Context, customer *customerv1.GetCustomerRequest) (*customerv1.GetCustomerResponse, error) {
	customerId := uuid.MustParse(customer.GetId())
	customerRes, err := server.queries.GetCustomer(
		ctx,
		customerId,
	)
	if err != nil {
		return nil, err
	}

	return &customerv1.GetCustomerResponse{
		Id:          customerRes.CustomerID.String(),
		Name:        customerRes.CustomerName,
		Job:         customerRes.CustomerJob.String,
		Corporation: customerRes.CustomerCorporation.String,
		Address:     customerRes.CustomerAddress.String,
		Phone:       customerRes.ContactPhone.String,
		Mail:        customerRes.ContactMail.String,
		Fax:         customerRes.ContactFax.String,
		Memo:        customerRes.CustomerMemo.String,
		Contact: &contactv1.Contact{
			Id:    customerRes.ContactID.String(),
			Phone: customerRes.ContactPhone.String,
			Mail:  customerRes.ContactMail.String,
			Fax:   customerRes.ContactFax.String,
		},
	}, nil
}

func (server *CustomerService) GetCustomerByBookId(ctx context.Context, customer *customerv1.GetCustomerByBookIdRequest) (*customerv1.GetCustomerByBookIdResponse, error) {
	bookId := uuid.MustParse(customer.GetBookId())
	customers, err := server.queries.GetCustomerByBookId(
		ctx,
		db.GetCustomerByBookIdParams{
			BookID: bookId,
			Limit:  int32(customer.GetLimit()),
			Offset: int32(customer.GetPage()),
		},
	)
	if err != nil {
		return nil, err
	}

	customersRes := make([]*customerv1.Customer, len(customers))
	for i, customer := range customers {
		customersRes[i] = &customerv1.Customer{
			Id:          customer.ID.String(),
			Name:        customer.Name,
			Job:         customer.Job.String,
			Corporation: customer.Corporation.String,
			Address:     customer.Address.String,
			Memo:        customer.Memo.String,
		}
	}

	return &customerv1.GetCustomerByBookIdResponse{
		Customers: customersRes,
	}, nil
}
