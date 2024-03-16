package ioc

import (
	intrRepov1 "github.com/Gnoloayoul/JGEBCamp/webook/api/proto/gen/intrRepo/v1"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/grpc/client"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGRPCInteractiveRepositoryClient(repo repository.InteractiveRepository) intrRepov1.InteractiveRepositoryClient {
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

	var opts []grpc.DialOption
	if cfg.Secure {

	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}

	remote := intrRepov1.NewInteractiveRepositoryClient(cc)
	local := client.NewInteractiveRepositoryAdapter(repo)
	res := client.NewGreyscaleInteractiveRepositoryClient(remote, local)

	viper.OnConfigChange(func(in fsnotify.Event){
		var newCfg Config
		err = viper.UnmarshalKey("grpc.client.intr", &newCfg)
		if err != nil {
			panic(err)
		}
		res.UpdateThreshold(newCfg.Threshold)
	})

	return res
}