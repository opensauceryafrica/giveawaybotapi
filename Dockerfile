# syntax = docker/dockerfile:1.2

FROM golang:1.18-alpine
ENV CGO_ENABLED=0

# # define build args and environment variables
ARG PORT
ENV PORT $PORT

# mount env file
RUN --mount=type=secret,id=_env,dst=/etc/secrets/.env cat /etc/secrets/.env

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o giveawaybot .

EXPOSE $PORT

CMD [ "./giveawaybot" ]