package aliyun

import (
	"context"
	"errors"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/goccy/go-json"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
	"time"
)

type Service struct {
	client *sms.Client
}

func NewService(client *sms.Client) *Service {
	return &Service{
		client: client,
	}
}

// SendSms
// 单次
func (s *Service) SendSms(ctx context.Context, signName, tplCode string, phone []string) error {
	phoneLen := len(phone)

	for i := 0; i < phoneLen; i++ {
		phoneSignle := phone[i]

		// 1.生成验证码
		code := fmt.Sprintf("%06v",
			rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000000))

		// 强耦合式生成验证码
		bcode, err := json.Marshal(map[string]interface{}{
			"code": code,
		})
		if err != nil {
			return err
		}

		// 2. 初始化短信结构体
		smsRequest := &sms.SendSmsRequest{
			SignName:      ekit.ToPtr[string](signName),
			TemplateCode:  ekit.ToPtr[string](tplCode),
			PhoneNumbers:  ekit.ToPtr[string](phoneSignle),
			TemplateParam: ekit.ToPtr[string](string(bcode)),
		}

		// 3. 发送短信
		smsResponse, err := s.client.SendSms(smsRequest)
		if err != nil {
			return err
		}
		if *smsResponse.Body.Code == "OK" {
			fmt.Println(phoneSignle, string(bcode))
			fmt.Printf("发送手机号： %s 的短信成功， 验证码为【%s】\n", phoneSignle, code)
		}
		fmt.Println(error.New(*smsResponse.Body.Message))
	}
	return nil
}
