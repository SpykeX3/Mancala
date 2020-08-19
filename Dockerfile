FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY . .

FROM build AS build_go
RUN go mod download && go build -o /server /go/src/app/Server/

FROM build AS build_elm
RUN apk add --no-cache curl && cd elm && curl -L -o elm.gz https://github.com/elm/compiler/releases/download/0.19.1/binary-for-linux-64-bit.gz && gunzip elm.gz && chmod +x elm && ./elm make src/Main.elm

FROM alpine:3.7
COPY --from=build_go /server /
COPY --from=build_elm /go/src/app/elm/index.html /elm/

EXPOSE 1337

CMD ["/server"]