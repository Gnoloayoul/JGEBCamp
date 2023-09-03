package memory

import (
	"context"
	"fmt"
)

type Service struct {

}

type NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}