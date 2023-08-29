package tencent

import (
	"github.com/ecodeclub/ekit"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type ServiceV1 struct {
	client *sms.Client
	appId *string
	signName *string
}

func NewServiceV1(c *sms.Client, appId string, signName string) *Service {
	return &Service{
		client: c,
		appId: toPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *ServiceV1) Send(ctx context.Context, tplId string, args map[string]string, numbers ...string) error {
	req := sms.
}















