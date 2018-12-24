FROM golang:latest

# Setup the work directory
WORKDIR /home

# Copy the files to WORKDIR
COPY main.go /home/main.go

# Install the used go packages
RUN ["go", "get", "git.darknebu.la/GalaxySimulator/structs"]
RUN ["go", "get", "github.com/gorilla/mux"]

# Start the webserver
ENTRYPOINT ["go", "run", "/home/main.go"]
