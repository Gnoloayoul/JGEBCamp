# set default
这里是第一种用法的示例  
需要注意的地方：  
- 用该方法设置的 default ，必须与其他内容**无相关**
````
 例如
 db.mysql:
    dsn: xxxxxxxxxxxx
    name: 111111 # 这里用该方法写为 Default
    
那么在 viper.UnmarshalKey 时，是读不到该 Default 的 name
````