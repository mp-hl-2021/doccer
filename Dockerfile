FROM golang:1.16.2-alpine3.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o doccer-server main.go


FROM golang:1.16.2-alpine3.13
COPY --from=builder /build/doccer-server .
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# executable
ENTRYPOINT [ "./doccer-server" ]