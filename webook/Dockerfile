# 基础镜像
FROM ubuntu:20.04
# 把编译后的打包进（/app）这个目录
COPY webook /app/webook
# 给执行文件
RUN chmod 775 /app/webook
# 设定工作目录为 /app
WORKDIR /app

ENTRYPOINT ["/app/webook"]