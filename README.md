# csv-server

### Building the docker container
1. Clone this repository
2. Run `docker build . -t ishankhare07/csv-server` from inside the cloned repo.
3. Run `docker run -d -p 80:8080 ishankhare07/csv-server`
> the port mapping 80->8080 is done because the postman collection queries at port 80
