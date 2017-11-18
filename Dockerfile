# Dockerfile for the go application

# Using the official golang repository
From golang:1.9.2-stretch

# Working directory
# TODO switch to somthing that makes sense (Require knowledge about image)

# Update from repository commented while testing
RUN apt-get update && apt-get install -y

# Go getting  govendor and mgo
RUN go get -u github.com/kardianos/govendor \ 
	gopkg.in/mgo.v2 \
 	gopkg.in/mgo.v2/bson

#Make directory for application
RUN mkdir -p /go/src/github.com/Xillez/CloudTechAssign3

WORKDIR /go/src/github.com/Xillez/CloudTechAssign3

#Adding the application
ADD  ./CloudTechAssign3 .

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Echo hello to make sure it is alive
CMD echo "hello"