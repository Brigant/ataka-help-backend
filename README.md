# ataka-help-backend.
This is a sample Attacka Help api server

# How to run.
Install docker [docker](https://docs.docker.com/engine/install/).
Install docker-compose [docker-compose](https://pkg.go.dev/github.com/docker/compose/v2#section-readme)
Make copy of .env_example  with name .env. Specify the need configuration values

## For development
For the development it is usefull make the application runs  without rebuiling of the application container.It just run only database container and runs localy the app code
For that purpose  the command exists
```
make dev-run 
```
## For production
### How to run at first time or rebuild application
Run command
```
make build
```
### Simple run
It runs existed containers
```
make run
```

### How to stop
Run command
```
make stop
```
