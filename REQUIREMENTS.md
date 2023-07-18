### The brief

The task is to implement a backend API for a set of underwater sensors. Unfortunately, we don’t have the sensors themselves physically installed (yet), so you will also need to implement some simulation (more on that later).

### The domain

Each sensor has following properties:

* codename
* 3D-coordinates (x, y, z)
* data output rate (in seconds)

Sensors are organised into the groups; each group has a name (we will use the greek letter names for these: alpha, beta, gamma, etc.). Sensor’s codename consists of the group name and the sensor’s index (integer number, monotonically increasing within a group), for example: gamma 3.

Sensors provide the following data:

* temperature (Celsius, floating point number)
* transparency (%, integer number)
* the list of fish species detected near the sensor since the previous measurement; for each fish specie sensor returns the name (string) & the count (integer), for example:
  * Atlantic Cod: 12
  * Sailfish: 4

### The task

The task is to implement the backend service which will continuously generate the (fake) data, described in the[ Domain](https://www.notion.so/NestJS-test-task-1f431b598ec24ddca70b1388c2e1d009?pvs=21) section; and will provide an API to access this data. The service should be implemented using Go language (you are free to choose the web framework for the task). The data should be stored in PostgreSQL database. It is expected that the database for running the project locally will be described in a Docker Compose file — so the project could be started locally via the npm / docker commands only. It is also expected that API will have a Swagger spec. End-to-end tests would be a big plus.

It might be useful to split the data generation process into three phases:

1. One-time “kickoff” phase: generate the sensor groups and the sensors in each group.
2. Regularly repeated phase: generate the data for sensors. Fish specie names could be taken from[ this page](https://oceana.org/ocean-fishes/). For each sensor the data should change once per its data output rate. The data should be randomized
3. Regularly repeated phase: aggregate statistics (to see what statistics to gather and how often please read the API description below).

The API should provide the following endpoints:

* /group/<groupName>/transparency/average: current average transparency inside the group
* /group/<groupName>/temperature/average: current average temperature inside the group
* /group/<groupName>/species: full list of species (with counts) currently detected inside the group
* /group/<groupName>/species/top/<N>: list of top N species (with counts) currently detected inside the group
* /region/temperature/min?xMin=<xMin>&xMax=<xMax>&yMin=<yMin>&yMax=<yMax>&zMin=<zMin>&zMax=<zMax>: current minimum temperature inside the region\*\*\*\*;\*\*\*\* region here and below is an area represented by the range of coordinates
* /region/temperature/max?xMin=<xMin>&xMax=<xMax>&yMin=<yMin>&yMax=<yMax>&zMin=<zMin>&zMax=<zMax>: current maximum temperature inside the region
* /sensor/<codeName>/temperature/average?from=<fromDateTime>&till=<untillDateTime>: average temperature detected by a particular sensor between the specified date/time pairs (UNIX timestamps)

The /group/<groupName>/species/top/<N> endpoint also should (optionally) support ?from=<fromDateTime>&till=<untillDateTime> parameters.

The results of the the following endpoints should be cached in Redis with 10s TTL:

* /group/<groupName>/transparency/average
* /group/<groupName>/temperature/average

Please note that Redis connection should be configurable through environment variables. While running the project locally, there should be a convenient way to bring Redis instance up via Docker Compose.

**
