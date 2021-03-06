# build
FROM golang:1.17-alpine as build
ENV GO111MODULE=on

WORKDIR /go/src/app

RUN apk --no-cache add make ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY backend/kagu_miru_api backend/kagu_miru_api
COPY backend/internal backend/internal
COPY backend/pkg backend/pkg

RUN CGO_ENABLED=0 go build -o bin/server -ldflags "-w -s" ./backend/kagu_miru_api

# exec
FROM scratch
COPY --from=build /go/src/app/bin/server ./server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT ["./server"]
