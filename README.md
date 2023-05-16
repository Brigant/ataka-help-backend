# Ataka-help-backend.
This is a simple rest api server for Attacka Help site.

# How to run!
Install docker [docker](https://docs.docker.com/engine/install/). 
Install docker-compose [docker-compose](https://pkg.go.dev/github.com/docker/compose/v2#section-readme).
Make copy of `.env_example`  with name `.env`. Specify the need configuration values.

## For development
For the development it is usefull make the application runs  without rebuiling of the application container. It just run only database container and runs localy the app code.
For that purpose  the command exists:
```
make dev-run 
```
For migration up: 
```
make migration
```
!!! For downward migration in the file change the argument `"-all"` to the number of desired downward migrations, for example `"1"`. Use command:
```
make migrate-down
```

## For production
### How to run at first time or rebuild application
Run command:
```
make build
```
### Simple run
It runs existed containers:
```
make run
```

### How to stop
Run command:
```
make stop
```
