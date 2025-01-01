# syntax=docker/dockerfile:1

FROM golang:1.22.10-alpine3.20

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
  apk update && \
  apk add --no-cache upx make build-base bash git redis

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./
RUN go mod download

COPY . .

COPY db/schema.sql /docker-entrypoint-initdb.d/
COPY db/migrations/*.sql /docker-entrypoint-initdb.d/migrations/

RUN apk add --no-cache postgresql-client

RUN go build -o /main

#Environment VARS for deployment
ENV APT_ENV=dev 

ENV DB_DRIVER="postgres"
ENV DB_HOST="localhost"
ENV DB_NAME="apt_registry_development"
ENV DB_USER="dev_user"
ENV DB_PASSWORD="password"
ENV DB_PORT=5432
ENV DB_USE_SSL=false

ENV COOKIE_HASH_KEY='y0b6|UBJQ(N$KB)jAJYL-aj=:q?;yK64^TPch0=|1XNnv{X@QrL#?80u$1]LcBF'
ENV COOKIE_BLOCK_KEY='4Qdnm4acxfAILGEFQ3jUj0PoLbMWbyMm'

ENV COOKIE_DOMAIN="localhost"
ENV SESSION_MAX_AGE=43200
ENV SESSION_COOKIE_NAME="aptrust_session"
ENV FLASH_COOKIE_NAME="aptrust_flash"
ENV PREFS_COOKIE_NAME="aptrust_prefs"

ENV HTTPS_COOKIES=false 

ENV NSQ_URL="http://localhost:4151"

ENV AWS_ACCESS_KEY_ID=<yourkeyid>
ENV AWS_SECRET_ACCESS_KEY=<yoursecretkey>
ENV AWS_REGION=us-east-1 

ENV ENABLE_TWO_FACTOR_SMS=true

ENV ENABLE_TWO_FACTOR_AUTHY=true 
ENV AUTHY_API_KEY=<yourkey> 

ENV OTP_EXPIRATION="15m" 

ENV EMAIL_ENABLED=false
ENV EMAIL_FROM_ADDRESS="help@aptrust.org" 

ENV REDIS_DEFAULT_DB=0 
ENV REDIS_PASSWORD=""
ENV REDIS_URL="localhost:6379"

ENV LOG_FILE="STDOUT"
ENV LOG_LEVEL=0
ENV LOG_CALLER=false
ENV LOG_TO_CONSOLE=true
ENV LOG_SQL=false
ENV AWS_SES_PWD=password
ENV AWS_SES_USER=system@user.org

EXPOSE 8080

CMD [ "/main" ]
