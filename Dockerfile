FROM golang:1.16 as builder

WORKDIR /app

# Cache might be used here.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest.
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -v -o cmd/game cmd/main.go

# Use a Docker multi-stage build.
# Multistage builds in Docker are a powerful feature that allows you to create
# lean and secure images by separating the build environment from the runtime
# environment. You can use multiple FROM statements in a Dockerfile, where
# each FROM instruction can use a different base image, and can be considered a
# separate stage of the build. This approach enables you to include tools and
# dependencies needed for building the application in the initial stages without
# adding unnecessary size and potential security vulnerabilities to the final
# image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a non-root user and group
RUN addgroup -S game && adduser -S executor -G game

# Set the working directory to /home/executor (a directory the non-root user has access to)
WORKDIR /home/executor

# Copy the binary and static from the builder stage to the workdir.
COPY --from=builder /app/cmd/game ./app/cmd/game
COPY --from=builder /app/configs ./app/configs
COPY --from=builder /app/web ./app/web

# Change the ownership of the binary to the non-root user
RUN chown executor:game app

# Use the non-root user to run the application
USER executor

EXPOSE 4000

WORKDIR /home/executor/app

# Dubug with: CMD ["sleep","3600"]
CMD ["./cmd/game"]
