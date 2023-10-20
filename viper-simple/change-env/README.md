# viper 不同环境加载不同配置文件
## 其实就是对 viper 传入参数
## 示例
假设在工作环境的 ./config 里有 config.yaml 与 config1.yaml
````
# config.yaml
Table:
  name: "123456"
  
# config1.yaml
Table:
  name: "asdbef"
````
想通过输入参数来指定 viper 读取那份配置文件，代码可以这么写
````
cfile := pflag.String("config", "config/config.yaml", "配置文件路径“)
pflag.Parse()
viper.SetConfigFile(*cfile)
err := viper.ReadInConfig()
if err != nil {
    panic(err)
}
````
这里 `pflag.String` 一句，  
设定了参数名为 `config` ，默认值（就是没有输入的话就读谁）为 `config/config.yaml`  
当有 `config` 的输入，就是由 `pflag.Parse()` 这句来解析

## 实际效果
### 在 Goland 中
#### 默认
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\IDE_none.png)
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\IDE_none_res.png)
#### config.yaml
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\IDE_1.png)
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\IDE_1_res.png)
#### config1.yaml
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\IDE_2.png)
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\IDE_2_res.png)
### 在命令行中
![](C:\Users\Administrator\Desktop\JGEBCamp\viper-simple\change-env\pic\cmd.png)