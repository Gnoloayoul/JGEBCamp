package ioc

import (
	intrv1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intr/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGRPCInteractiveServiceClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr string
		Secure bool
		Threshold int32
	}

	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}

	// 配置连接安全部分
	var opts []grpc.DialOption
	if cfg.Secure {
		// 安全证书什么的
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 启动 grpc 部分
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}

	// 构建client
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewInteractiveServiceAdapter(svc)
	res := client.NewGreyScaleInteractiveServiceClient(remote, local)

	// 监听配置后续变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		var cfg Config
		err := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			// 输出日志
		}
		res.UpdateThreshold(cfg.Threshold)
	})

	return res
}