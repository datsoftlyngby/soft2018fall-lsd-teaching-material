FROM golang:1.11

# Install any needed packages
RUN go get github.com/prometheus/client_golang/prometheus

# Set the working directory to /app
WORKDIR /src

# Copy the current directory contents into the container at /app
COPY server.go /src/server.go

# Make port 8080 available to the world outside this container
EXPOSE 8080

# set the hostname of the dbserver to the name given in the docker-compose file
ENV IP="0.0.0.0"
ENV PORT="8080"

# Run app when the container launches
RUN go build /src/server.go
ENTRYPOINT ./server