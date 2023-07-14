### SENSOR - API

> **Note**: Before running the application locally you need an `config.yaml` file

All commands and their descriptions are described in `Makefile`

## local start

You can launch all aplication using one single command `make compose_up`

### routes:

- `/api/v1/group/:groupName/transparency/average` - [method GET] -current average transparency inside the group
- `/api/v1/group/:groupName/temperature/average` - [method GET] - current average temperature inside the group
- `/api/v1/group/:groupName/species` - [method GET] -  full list of species (with counts) currently detected inside the group
- `/api/v1/group/:groupName/species/:top` - [method GET] - list of top N species (with counts) currently detected inside the group  support ?from=<fromDateTime>&till=<untillDateTime> parameters

Swagger documentation can see on `http://localhost:5000/swagger/`
