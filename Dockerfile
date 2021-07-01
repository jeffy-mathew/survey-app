FROM golang:1.16-alpine as builder
RUN cd ..
RUN mkdir survey-platform
WORKDIR survey-platform
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor  -ldflags "-X 'survey-platform/internal/app.ApiVersion=1.0.0'"  -o survey-platform ./cmd/main.go

FROM alpine
RUN mkdir survey-platform
WORKDIR app
COPY --from=builder /go/survey-platform/survey-platform /app/
ENTRYPOINT ["/app/survey-platform"]
