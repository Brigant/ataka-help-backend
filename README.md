# Ataka-help-backend.
This is a simple rest api server for Attacka Help site.

# How to run!
Install docker [docker](https://docs.docker.com/engine/install/). 
Install docker-compose [docker-compose](https://pkg.go.dev/github.com/docker/compose/v2#section-readme).
Make copy of `.env_example`  with name `.env`. Specify the need configuration values.

## For development
For the local development it is usefull make the application run without rebuiling of the application container. It just runs only database container and runs the app code localy.
For that purpose the command exists:
```
make run 
```
For migration up: 
```
make migration-up
```
!!! For downward migration in the file change the argument `"-all"` to the number of desired downward migrations, for example `"1"`. Use command:
```
make migrate-down
```

## For production
### How to run at first time or rebuild application
Run command:
```
make dbrun
make build
```
After add defaul user to database. Use sql cript in the **seed** directory for help
### Simple run
It runs existed images:
```
make up
```

### How to stop
Run command:
```
make down
```
