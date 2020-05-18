FROM golang:1.14.2-alpine3.11

# Need git for dep
RUN apk add --no-cache git

# Need npm for react
# RUN apk add --update npm

# Copy current dir (outside docker) to the proper directory (inside docker)
COPY . src/github.com/etzelm/blog-in-golang/

# Change current directory
WORKDIR src/github.com/etzelm/blog-in-golang/

# Build react app except it doesn't produce the same results as local build
# Probably has something to do with npm version, for now I'll build locally
# RUN npm i
# RUN npm run build

# Change current directory
# WORKDIR ../

# Build go server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-in-golang .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/blog-in-golang .
COPY --from=0 /go/src/github.com/etzelm/blog-in-golang/realtor/build/* /realtor/build/
COPY . .
CMD ["./blog-in-golang"]

EXPOSE 80
EXPOSE 443