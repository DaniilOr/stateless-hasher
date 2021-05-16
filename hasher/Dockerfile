FROM golang:1.15-alpine AS build
ADD . /app
ENV CGO_ENABLED 0
WORKDIR /app
RUN go build -o hasher ./cmd

FROM scratch
COPY --from=build /app/hasher /hasher
CMD ["/hasher"]
EXPOSE 9999