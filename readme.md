# Local set up for unit test

## Database

This project assume that you have installed docker engine locally or docker-machine - in case you are a windows or mac user.
If you don't, refer this: [install docker](https://docs.docker.com/engine/installation/)
If you are using windows or mac, confirm you ip address of the vm running docker engine. Replace the ip address in main.go (The address 192.168.99.100:28015).

Run command:
~~~~
docker run -d -p 8080:8080 -p 28015:28015 --name rethinkdb rethinkdb
~~~~

## Run the application

Run command(in the folder where this repository located)
~~~~
go run main.go
~~~~

## Test the api
### sign up new user:

url: http://localhost:8000/signup
method: post
payload sample(json):
~~~~
{
  "email": "test@abc.com",
  "name": "tester",
  "password" "password"
}
~~~~

### login
url: http://localhost:8000/login
method: post
payload sample(json):
~~~~
{
  "email": "test@abc.com",
  "password" "password"
}
~~~~

### check the database:
url: http://ip_address_of_your_vm_running_docker:8080

# Integrated version
