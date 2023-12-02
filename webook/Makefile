.PHONY: mock
mock:
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\web\jwt\types.go -package=jwtmocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\web\jwt\mocks\handler.mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\user.go -package=repomocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\mocks\user_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\code.go -package=repomocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\mocks\code_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\dao\user.go -package=daomocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\dao\mocks\user_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\article\article.go -package=artrepomocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\article\mocks\article_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\article\article_author.go -package=artrepomocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\article\mocks\article_author_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\article\article_reader.go -package=artrepomocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\article\mocks\article_reader_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\cache\user.go -package=cachemocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\cache\mocks\user_mock.go
	@mockgen -package=redismocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\repository\cache\redismocks\cmdable_mock.go github.com/redis/go-redis/v9 Cmdable
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\pkg\ratelimit\types.go -package=ratelimitmocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\pkg\ratelimit\mocks\types_mock.go
	@mockgen -source=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\service\sms\types.go -package=smsmocks -destination=C:\Users\Administrator\Desktop\JGEBCamp\webook\internal\service\sms\mocks\types_mock.go
	@go mod tidy