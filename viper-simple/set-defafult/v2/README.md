# set default
这里是第二种用法的示例
就效果而言，比方法1用途大
## 特性1：
可以在要用的时候才这样设定
## 特性2：  
- 如何同时还用了 viper.ReadInConfig 读取配置文件，且配置文件里也有相同字段，那么在这里的写入会覆盖掉配置文件中相同位置的内容
````
 例如
 
 配置文件
 -----------
    db.mysql:
      dns: ""
 -----------
 
 而在程序中
 -----------
    type Config struct {
        Name string `yaml:name`
    }
    cfg := Config{
        Name: "123456789",
    }
 -----------

那么在 viper.UnmarshalKey 之后，能读到cfg.Name，值为 123456789
````
- 如何同时还用了 viper.ReadInConfig 读取配置文件，且配置文件里没有相同字段，那么能写入并且能读取到
````
 例如
 
 配置文件
 -----------
     db.mysql:
       dns: "XXXXX"
 -----------
 
 在程序中
 -----------
	type Config struct {
		Name string `yaml:name`
	}
	cfg := Config{
		Name: "123456789",
	} 
 -----------
 那么在 viper.UnmarshalKey 之后，能读到cfg.Name，值为 123456789
````