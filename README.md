### SENSOR - API

> **Note**: Before running the application locally you need an `config.yaml` file

All commands and their descriptions are described in `Makefile`

## local start

You can launch all aplication using one single command `make compose_up`

### routes:

- `/api/v1/group/:groupName/transparency/average` - [method GET] -current average transparency inside the group.
  Example: `http://localhost:5000/api/v1/group/alpha/transparency/average`
- `/api/v1/group/:groupName/temperature/average` - [method GET] - current average temperature inside the group
  Example: `http://localhost:5000/api/v1/group/alpha/temperature/average`
- `/api/v1/group/:groupName/species` - [method GET] -  full list of species (with counts) currently detected inside the group
  Example: `http://localhost:5000/api/v1/group/alpha/species`
- `/api/v1/group/:groupName/species/:top` - [method GET] - list of top N species (with counts) currently detected inside the group  support ?from=<fromDateTime>&till=<untillDateTime> parameters
  Example: `http://localhost:5000/api/v1/group/alpha/species/3?from=1689278400&till=1689364800`
- `/api/v1/region/temperature/min?xMin=xMin&xMax=xMax&yMin=yMin&yMax=yMax&zMin=zMin&zMax=zMax`:[method GET] current minimum temperature inside the region
  Example: `http://localhost:5000/api/v1/region/temperature/min?xMin=-8.213864897523635&xMax=7.868109888194829&yMin=-0.6530181503282156&yMax=4.494854709411525&zMin=-4.4049550107467885&zMax=-2.693363601487414`
- `/api/v1/region/temperature/max?xMin=xMin&xMax=xMax&yMin=yMin&yMax=yMax&zMin=zMin&zMax=zMax`: [method GET] current maximum temperature inside the region
  Example: `http://localhost:5000/api/v1/region/temperature/max?xMin=-8.213864897523635&xMax=7.868109888194829&yMin=-0.6530181503282156&yMax=4.494854709411525&zMin=-4.4049550107467885&zMax=-2.693363601487414`
- `/api/v1/sensor/:codename/temperature/average` : [method GET] average temperature detected by a particular sensor between the specified date/time pairs (UNIX timestamps)
  Example: `http://localhost:5000/api/v1/sensor/alpha5/temperature/average?from=1689278400&till=1689599444`

Swagger documentation can see on `http://localhost:5000/swagger/`

### Tests:

for running integration test use command: `make int_test`
