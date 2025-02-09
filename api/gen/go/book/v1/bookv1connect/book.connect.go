// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: book/v1/book.proto

package bookv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/FotiadisM/mock-microservice/api/gen/go/book/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// BookServiceName is the fully-qualified name of the BookService service.
	BookServiceName = "book.v1.BookService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// BookServiceGetBookProcedure is the fully-qualified name of the BookService's GetBook RPC.
	BookServiceGetBookProcedure = "/book.v1.BookService/GetBook"
	// BookServiceListBookProcedure is the fully-qualified name of the BookService's ListBook RPC.
	BookServiceListBookProcedure = "/book.v1.BookService/ListBook"
	// BookServiceCreateBookProcedure is the fully-qualified name of the BookService's CreateBook RPC.
	BookServiceCreateBookProcedure = "/book.v1.BookService/CreateBook"
	// BookServiceUpdateBookProcedure is the fully-qualified name of the BookService's UpdateBook RPC.
	BookServiceUpdateBookProcedure = "/book.v1.BookService/UpdateBook"
	// BookServiceDeleteBookProcedure is the fully-qualified name of the BookService's DeleteBook RPC.
	BookServiceDeleteBookProcedure = "/book.v1.BookService/DeleteBook"
	// BookServiceGetAuthorProcedure is the fully-qualified name of the BookService's GetAuthor RPC.
	BookServiceGetAuthorProcedure = "/book.v1.BookService/GetAuthor"
	// BookServiceListAuthorProcedure is the fully-qualified name of the BookService's ListAuthor RPC.
	BookServiceListAuthorProcedure = "/book.v1.BookService/ListAuthor"
	// BookServiceCreateAuthorProcedure is the fully-qualified name of the BookService's CreateAuthor
	// RPC.
	BookServiceCreateAuthorProcedure = "/book.v1.BookService/CreateAuthor"
	// BookServiceUpdateAuthorProcedure is the fully-qualified name of the BookService's UpdateAuthor
	// RPC.
	BookServiceUpdateAuthorProcedure = "/book.v1.BookService/UpdateAuthor"
	// BookServiceDeleteAuthorProcedure is the fully-qualified name of the BookService's DeleteAuthor
	// RPC.
	BookServiceDeleteAuthorProcedure = "/book.v1.BookService/DeleteAuthor"
	// BookServiceThrowPanicProcedure is the fully-qualified name of the BookService's ThrowPanic RPC.
	BookServiceThrowPanicProcedure = "/book.v1.BookService/ThrowPanic"
	// BookServiceThrowServiceErrorProcedure is the fully-qualified name of the BookService's
	// ThrowServiceError RPC.
	BookServiceThrowServiceErrorProcedure = "/book.v1.BookService/ThrowServiceError"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	bookServiceServiceDescriptor                 = v1.File_book_v1_book_proto.Services().ByName("BookService")
	bookServiceGetBookMethodDescriptor           = bookServiceServiceDescriptor.Methods().ByName("GetBook")
	bookServiceListBookMethodDescriptor          = bookServiceServiceDescriptor.Methods().ByName("ListBook")
	bookServiceCreateBookMethodDescriptor        = bookServiceServiceDescriptor.Methods().ByName("CreateBook")
	bookServiceUpdateBookMethodDescriptor        = bookServiceServiceDescriptor.Methods().ByName("UpdateBook")
	bookServiceDeleteBookMethodDescriptor        = bookServiceServiceDescriptor.Methods().ByName("DeleteBook")
	bookServiceGetAuthorMethodDescriptor         = bookServiceServiceDescriptor.Methods().ByName("GetAuthor")
	bookServiceListAuthorMethodDescriptor        = bookServiceServiceDescriptor.Methods().ByName("ListAuthor")
	bookServiceCreateAuthorMethodDescriptor      = bookServiceServiceDescriptor.Methods().ByName("CreateAuthor")
	bookServiceUpdateAuthorMethodDescriptor      = bookServiceServiceDescriptor.Methods().ByName("UpdateAuthor")
	bookServiceDeleteAuthorMethodDescriptor      = bookServiceServiceDescriptor.Methods().ByName("DeleteAuthor")
	bookServiceThrowPanicMethodDescriptor        = bookServiceServiceDescriptor.Methods().ByName("ThrowPanic")
	bookServiceThrowServiceErrorMethodDescriptor = bookServiceServiceDescriptor.Methods().ByName("ThrowServiceError")
)

// BookServiceClient is a client for the book.v1.BookService service.
type BookServiceClient interface {
	// Book rpcs
	GetBook(context.Context, *connect.Request[v1.GetBookRequest]) (*connect.Response[v1.GetBookResponse], error)
	ListBook(context.Context, *connect.Request[v1.ListBookRequest]) (*connect.Response[v1.ListBookResponse], error)
	CreateBook(context.Context, *connect.Request[v1.CreateBookRequest]) (*connect.Response[v1.CreateBookResponse], error)
	UpdateBook(context.Context, *connect.Request[v1.UpdateBookRequest]) (*connect.Response[v1.UpdateBookResponse], error)
	DeleteBook(context.Context, *connect.Request[v1.DeleteBookRequest]) (*connect.Response[v1.DeleteBookResponse], error)
	// Authro rpcs
	GetAuthor(context.Context, *connect.Request[v1.GetAuthorRequest]) (*connect.Response[v1.GetAuthorResponse], error)
	ListAuthor(context.Context, *connect.Request[v1.ListAuthorRequest]) (*connect.Response[v1.ListAuthorResponse], error)
	CreateAuthor(context.Context, *connect.Request[v1.CreateAuthorRequest]) (*connect.Response[v1.CreateAuthorResponse], error)
	UpdateAuthor(context.Context, *connect.Request[v1.UpdateAuthorRequest]) (*connect.Response[v1.UpdateAuthorResponse], error)
	DeleteAuthor(context.Context, *connect.Request[v1.DeleteAuthorRequest]) (*connect.Response[v1.DeleteAuthorResponse], error)
	// Exampl rpcs
	ThrowPanic(context.Context, *connect.Request[v1.ThrowPanicRequest]) (*connect.Response[v1.ThrowPanicResponse], error)
	ThrowServiceError(context.Context, *connect.Request[v1.ThrowServiceErrorRequest]) (*connect.Response[v1.ThrowServiceErrorResponse], error)
}

// NewBookServiceClient constructs a client for the book.v1.BookService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewBookServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) BookServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &bookServiceClient{
		getBook: connect.NewClient[v1.GetBookRequest, v1.GetBookResponse](
			httpClient,
			baseURL+BookServiceGetBookProcedure,
			connect.WithSchema(bookServiceGetBookMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listBook: connect.NewClient[v1.ListBookRequest, v1.ListBookResponse](
			httpClient,
			baseURL+BookServiceListBookProcedure,
			connect.WithSchema(bookServiceListBookMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createBook: connect.NewClient[v1.CreateBookRequest, v1.CreateBookResponse](
			httpClient,
			baseURL+BookServiceCreateBookProcedure,
			connect.WithSchema(bookServiceCreateBookMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateBook: connect.NewClient[v1.UpdateBookRequest, v1.UpdateBookResponse](
			httpClient,
			baseURL+BookServiceUpdateBookProcedure,
			connect.WithSchema(bookServiceUpdateBookMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteBook: connect.NewClient[v1.DeleteBookRequest, v1.DeleteBookResponse](
			httpClient,
			baseURL+BookServiceDeleteBookProcedure,
			connect.WithSchema(bookServiceDeleteBookMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getAuthor: connect.NewClient[v1.GetAuthorRequest, v1.GetAuthorResponse](
			httpClient,
			baseURL+BookServiceGetAuthorProcedure,
			connect.WithSchema(bookServiceGetAuthorMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listAuthor: connect.NewClient[v1.ListAuthorRequest, v1.ListAuthorResponse](
			httpClient,
			baseURL+BookServiceListAuthorProcedure,
			connect.WithSchema(bookServiceListAuthorMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createAuthor: connect.NewClient[v1.CreateAuthorRequest, v1.CreateAuthorResponse](
			httpClient,
			baseURL+BookServiceCreateAuthorProcedure,
			connect.WithSchema(bookServiceCreateAuthorMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateAuthor: connect.NewClient[v1.UpdateAuthorRequest, v1.UpdateAuthorResponse](
			httpClient,
			baseURL+BookServiceUpdateAuthorProcedure,
			connect.WithSchema(bookServiceUpdateAuthorMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteAuthor: connect.NewClient[v1.DeleteAuthorRequest, v1.DeleteAuthorResponse](
			httpClient,
			baseURL+BookServiceDeleteAuthorProcedure,
			connect.WithSchema(bookServiceDeleteAuthorMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		throwPanic: connect.NewClient[v1.ThrowPanicRequest, v1.ThrowPanicResponse](
			httpClient,
			baseURL+BookServiceThrowPanicProcedure,
			connect.WithSchema(bookServiceThrowPanicMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		throwServiceError: connect.NewClient[v1.ThrowServiceErrorRequest, v1.ThrowServiceErrorResponse](
			httpClient,
			baseURL+BookServiceThrowServiceErrorProcedure,
			connect.WithSchema(bookServiceThrowServiceErrorMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// bookServiceClient implements BookServiceClient.
type bookServiceClient struct {
	getBook           *connect.Client[v1.GetBookRequest, v1.GetBookResponse]
	listBook          *connect.Client[v1.ListBookRequest, v1.ListBookResponse]
	createBook        *connect.Client[v1.CreateBookRequest, v1.CreateBookResponse]
	updateBook        *connect.Client[v1.UpdateBookRequest, v1.UpdateBookResponse]
	deleteBook        *connect.Client[v1.DeleteBookRequest, v1.DeleteBookResponse]
	getAuthor         *connect.Client[v1.GetAuthorRequest, v1.GetAuthorResponse]
	listAuthor        *connect.Client[v1.ListAuthorRequest, v1.ListAuthorResponse]
	createAuthor      *connect.Client[v1.CreateAuthorRequest, v1.CreateAuthorResponse]
	updateAuthor      *connect.Client[v1.UpdateAuthorRequest, v1.UpdateAuthorResponse]
	deleteAuthor      *connect.Client[v1.DeleteAuthorRequest, v1.DeleteAuthorResponse]
	throwPanic        *connect.Client[v1.ThrowPanicRequest, v1.ThrowPanicResponse]
	throwServiceError *connect.Client[v1.ThrowServiceErrorRequest, v1.ThrowServiceErrorResponse]
}

// GetBook calls book.v1.BookService.GetBook.
func (c *bookServiceClient) GetBook(ctx context.Context, req *connect.Request[v1.GetBookRequest]) (*connect.Response[v1.GetBookResponse], error) {
	return c.getBook.CallUnary(ctx, req)
}

// ListBook calls book.v1.BookService.ListBook.
func (c *bookServiceClient) ListBook(ctx context.Context, req *connect.Request[v1.ListBookRequest]) (*connect.Response[v1.ListBookResponse], error) {
	return c.listBook.CallUnary(ctx, req)
}

// CreateBook calls book.v1.BookService.CreateBook.
func (c *bookServiceClient) CreateBook(ctx context.Context, req *connect.Request[v1.CreateBookRequest]) (*connect.Response[v1.CreateBookResponse], error) {
	return c.createBook.CallUnary(ctx, req)
}

// UpdateBook calls book.v1.BookService.UpdateBook.
func (c *bookServiceClient) UpdateBook(ctx context.Context, req *connect.Request[v1.UpdateBookRequest]) (*connect.Response[v1.UpdateBookResponse], error) {
	return c.updateBook.CallUnary(ctx, req)
}

// DeleteBook calls book.v1.BookService.DeleteBook.
func (c *bookServiceClient) DeleteBook(ctx context.Context, req *connect.Request[v1.DeleteBookRequest]) (*connect.Response[v1.DeleteBookResponse], error) {
	return c.deleteBook.CallUnary(ctx, req)
}

// GetAuthor calls book.v1.BookService.GetAuthor.
func (c *bookServiceClient) GetAuthor(ctx context.Context, req *connect.Request[v1.GetAuthorRequest]) (*connect.Response[v1.GetAuthorResponse], error) {
	return c.getAuthor.CallUnary(ctx, req)
}

// ListAuthor calls book.v1.BookService.ListAuthor.
func (c *bookServiceClient) ListAuthor(ctx context.Context, req *connect.Request[v1.ListAuthorRequest]) (*connect.Response[v1.ListAuthorResponse], error) {
	return c.listAuthor.CallUnary(ctx, req)
}

// CreateAuthor calls book.v1.BookService.CreateAuthor.
func (c *bookServiceClient) CreateAuthor(ctx context.Context, req *connect.Request[v1.CreateAuthorRequest]) (*connect.Response[v1.CreateAuthorResponse], error) {
	return c.createAuthor.CallUnary(ctx, req)
}

// UpdateAuthor calls book.v1.BookService.UpdateAuthor.
func (c *bookServiceClient) UpdateAuthor(ctx context.Context, req *connect.Request[v1.UpdateAuthorRequest]) (*connect.Response[v1.UpdateAuthorResponse], error) {
	return c.updateAuthor.CallUnary(ctx, req)
}

// DeleteAuthor calls book.v1.BookService.DeleteAuthor.
func (c *bookServiceClient) DeleteAuthor(ctx context.Context, req *connect.Request[v1.DeleteAuthorRequest]) (*connect.Response[v1.DeleteAuthorResponse], error) {
	return c.deleteAuthor.CallUnary(ctx, req)
}

// ThrowPanic calls book.v1.BookService.ThrowPanic.
func (c *bookServiceClient) ThrowPanic(ctx context.Context, req *connect.Request[v1.ThrowPanicRequest]) (*connect.Response[v1.ThrowPanicResponse], error) {
	return c.throwPanic.CallUnary(ctx, req)
}

// ThrowServiceError calls book.v1.BookService.ThrowServiceError.
func (c *bookServiceClient) ThrowServiceError(ctx context.Context, req *connect.Request[v1.ThrowServiceErrorRequest]) (*connect.Response[v1.ThrowServiceErrorResponse], error) {
	return c.throwServiceError.CallUnary(ctx, req)
}

// BookServiceHandler is an implementation of the book.v1.BookService service.
type BookServiceHandler interface {
	// Book rpcs
	GetBook(context.Context, *connect.Request[v1.GetBookRequest]) (*connect.Response[v1.GetBookResponse], error)
	ListBook(context.Context, *connect.Request[v1.ListBookRequest]) (*connect.Response[v1.ListBookResponse], error)
	CreateBook(context.Context, *connect.Request[v1.CreateBookRequest]) (*connect.Response[v1.CreateBookResponse], error)
	UpdateBook(context.Context, *connect.Request[v1.UpdateBookRequest]) (*connect.Response[v1.UpdateBookResponse], error)
	DeleteBook(context.Context, *connect.Request[v1.DeleteBookRequest]) (*connect.Response[v1.DeleteBookResponse], error)
	// Authro rpcs
	GetAuthor(context.Context, *connect.Request[v1.GetAuthorRequest]) (*connect.Response[v1.GetAuthorResponse], error)
	ListAuthor(context.Context, *connect.Request[v1.ListAuthorRequest]) (*connect.Response[v1.ListAuthorResponse], error)
	CreateAuthor(context.Context, *connect.Request[v1.CreateAuthorRequest]) (*connect.Response[v1.CreateAuthorResponse], error)
	UpdateAuthor(context.Context, *connect.Request[v1.UpdateAuthorRequest]) (*connect.Response[v1.UpdateAuthorResponse], error)
	DeleteAuthor(context.Context, *connect.Request[v1.DeleteAuthorRequest]) (*connect.Response[v1.DeleteAuthorResponse], error)
	// Exampl rpcs
	ThrowPanic(context.Context, *connect.Request[v1.ThrowPanicRequest]) (*connect.Response[v1.ThrowPanicResponse], error)
	ThrowServiceError(context.Context, *connect.Request[v1.ThrowServiceErrorRequest]) (*connect.Response[v1.ThrowServiceErrorResponse], error)
}

// NewBookServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewBookServiceHandler(svc BookServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	bookServiceGetBookHandler := connect.NewUnaryHandler(
		BookServiceGetBookProcedure,
		svc.GetBook,
		connect.WithSchema(bookServiceGetBookMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceListBookHandler := connect.NewUnaryHandler(
		BookServiceListBookProcedure,
		svc.ListBook,
		connect.WithSchema(bookServiceListBookMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceCreateBookHandler := connect.NewUnaryHandler(
		BookServiceCreateBookProcedure,
		svc.CreateBook,
		connect.WithSchema(bookServiceCreateBookMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceUpdateBookHandler := connect.NewUnaryHandler(
		BookServiceUpdateBookProcedure,
		svc.UpdateBook,
		connect.WithSchema(bookServiceUpdateBookMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceDeleteBookHandler := connect.NewUnaryHandler(
		BookServiceDeleteBookProcedure,
		svc.DeleteBook,
		connect.WithSchema(bookServiceDeleteBookMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceGetAuthorHandler := connect.NewUnaryHandler(
		BookServiceGetAuthorProcedure,
		svc.GetAuthor,
		connect.WithSchema(bookServiceGetAuthorMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceListAuthorHandler := connect.NewUnaryHandler(
		BookServiceListAuthorProcedure,
		svc.ListAuthor,
		connect.WithSchema(bookServiceListAuthorMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceCreateAuthorHandler := connect.NewUnaryHandler(
		BookServiceCreateAuthorProcedure,
		svc.CreateAuthor,
		connect.WithSchema(bookServiceCreateAuthorMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceUpdateAuthorHandler := connect.NewUnaryHandler(
		BookServiceUpdateAuthorProcedure,
		svc.UpdateAuthor,
		connect.WithSchema(bookServiceUpdateAuthorMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceDeleteAuthorHandler := connect.NewUnaryHandler(
		BookServiceDeleteAuthorProcedure,
		svc.DeleteAuthor,
		connect.WithSchema(bookServiceDeleteAuthorMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceThrowPanicHandler := connect.NewUnaryHandler(
		BookServiceThrowPanicProcedure,
		svc.ThrowPanic,
		connect.WithSchema(bookServiceThrowPanicMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	bookServiceThrowServiceErrorHandler := connect.NewUnaryHandler(
		BookServiceThrowServiceErrorProcedure,
		svc.ThrowServiceError,
		connect.WithSchema(bookServiceThrowServiceErrorMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/book.v1.BookService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case BookServiceGetBookProcedure:
			bookServiceGetBookHandler.ServeHTTP(w, r)
		case BookServiceListBookProcedure:
			bookServiceListBookHandler.ServeHTTP(w, r)
		case BookServiceCreateBookProcedure:
			bookServiceCreateBookHandler.ServeHTTP(w, r)
		case BookServiceUpdateBookProcedure:
			bookServiceUpdateBookHandler.ServeHTTP(w, r)
		case BookServiceDeleteBookProcedure:
			bookServiceDeleteBookHandler.ServeHTTP(w, r)
		case BookServiceGetAuthorProcedure:
			bookServiceGetAuthorHandler.ServeHTTP(w, r)
		case BookServiceListAuthorProcedure:
			bookServiceListAuthorHandler.ServeHTTP(w, r)
		case BookServiceCreateAuthorProcedure:
			bookServiceCreateAuthorHandler.ServeHTTP(w, r)
		case BookServiceUpdateAuthorProcedure:
			bookServiceUpdateAuthorHandler.ServeHTTP(w, r)
		case BookServiceDeleteAuthorProcedure:
			bookServiceDeleteAuthorHandler.ServeHTTP(w, r)
		case BookServiceThrowPanicProcedure:
			bookServiceThrowPanicHandler.ServeHTTP(w, r)
		case BookServiceThrowServiceErrorProcedure:
			bookServiceThrowServiceErrorHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedBookServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedBookServiceHandler struct{}

func (UnimplementedBookServiceHandler) GetBook(context.Context, *connect.Request[v1.GetBookRequest]) (*connect.Response[v1.GetBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.GetBook is not implemented"))
}

func (UnimplementedBookServiceHandler) ListBook(context.Context, *connect.Request[v1.ListBookRequest]) (*connect.Response[v1.ListBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.ListBook is not implemented"))
}

func (UnimplementedBookServiceHandler) CreateBook(context.Context, *connect.Request[v1.CreateBookRequest]) (*connect.Response[v1.CreateBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.CreateBook is not implemented"))
}

func (UnimplementedBookServiceHandler) UpdateBook(context.Context, *connect.Request[v1.UpdateBookRequest]) (*connect.Response[v1.UpdateBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.UpdateBook is not implemented"))
}

func (UnimplementedBookServiceHandler) DeleteBook(context.Context, *connect.Request[v1.DeleteBookRequest]) (*connect.Response[v1.DeleteBookResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.DeleteBook is not implemented"))
}

func (UnimplementedBookServiceHandler) GetAuthor(context.Context, *connect.Request[v1.GetAuthorRequest]) (*connect.Response[v1.GetAuthorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.GetAuthor is not implemented"))
}

func (UnimplementedBookServiceHandler) ListAuthor(context.Context, *connect.Request[v1.ListAuthorRequest]) (*connect.Response[v1.ListAuthorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.ListAuthor is not implemented"))
}

func (UnimplementedBookServiceHandler) CreateAuthor(context.Context, *connect.Request[v1.CreateAuthorRequest]) (*connect.Response[v1.CreateAuthorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.CreateAuthor is not implemented"))
}

func (UnimplementedBookServiceHandler) UpdateAuthor(context.Context, *connect.Request[v1.UpdateAuthorRequest]) (*connect.Response[v1.UpdateAuthorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.UpdateAuthor is not implemented"))
}

func (UnimplementedBookServiceHandler) DeleteAuthor(context.Context, *connect.Request[v1.DeleteAuthorRequest]) (*connect.Response[v1.DeleteAuthorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.DeleteAuthor is not implemented"))
}

func (UnimplementedBookServiceHandler) ThrowPanic(context.Context, *connect.Request[v1.ThrowPanicRequest]) (*connect.Response[v1.ThrowPanicResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.ThrowPanic is not implemented"))
}

func (UnimplementedBookServiceHandler) ThrowServiceError(context.Context, *connect.Request[v1.ThrowServiceErrorRequest]) (*connect.Response[v1.ThrowServiceErrorResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("book.v1.BookService.ThrowServiceError is not implemented"))
}
