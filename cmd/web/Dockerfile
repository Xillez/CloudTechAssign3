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
 	gopkg.in/mgo.v2/bson \
	github.com/Xillez/CloudTechAssign3/cmd/web


#Move to the application within the container

# Start the application
CMD ["web"]