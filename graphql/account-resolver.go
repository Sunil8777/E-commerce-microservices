package main

type accountResolver struct {
	server *Server
}

func (a *accountResolver)  Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	
}