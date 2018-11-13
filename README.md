# csv-server

### Building the docker container
1. Clone this repository under your $GOPATH
2. Run `docker build . -t ishankhare07/csv-server` from inside the cloned repo.
3. Run `docker run -d -p 80:8080 ishankhare07/csv-server`
> the port mapping 80->8080 is done because the postman collection queries at port 80


### Assumptions made
The CSV file attached in the email was actually an excel file, and would have required the use of a much more complex library such as [github.com/tealeg/xlsx](https://github.com/tealeg/xlsx).  
Hence I've exported the excel file into csv and have used that as the input for the api.
The file is included in the repository as well.
