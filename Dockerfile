FROM alpine:3.6
RUN echo $'http://mirrors.aliyun.com/alpine/v3.6/main\n\
http://mirrors.aliyun.com/alpine/v3.6/community' > /etc/apk/repositories
RUN apk add --update ca-certificates
RUN update-ca-certificates
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
RUN mkdir /app
COPY main /app
CMD ["/app/main"]
