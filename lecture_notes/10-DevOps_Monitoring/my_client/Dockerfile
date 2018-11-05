FROM golang:1.11

# Install any needed packages
# RUN go get gopkg.in/mgo.v2

# Set the working directory to /app
WORKDIR /src

# Copy the current directory contents into the container at /app
COPY client_sim.go /src/client_sim.go

# set the hostname of the dbserver to the name given in the docker-compose file
ENV SERVER="192.168.33.10"
ENV PORT="8080"

# Run app.py when the container launches
RUN go build /src/client_sim.go
ENTRYPOINT ./client_sim