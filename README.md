# Httpmon
This application reads a HTTP request stream from either a log file or 
stdin. Then displays it in a terminal user interface made with termdash.

*Note: click quickly and repeatedly on UI tabs, sometimes the event is not
processed because the UI loop is not perfectly synced with the backend
loop.*

It can display the following metrics:
- requests per host: top hits among http callers
    - callers are sorted from top to bottom
    - shows top 5 request repartition with sorted horizontal bars
    - be wary of `cardinality` here, the more hostnames, the more memory it takes
    - filtering input fields would be a good feature to add
- requests per second: it shows a plot of total requests per `--period`.
    - it is possible to zoom in on data points using the mouse
    - on ordinate is the average request per second, the abscissa being
    the id of the data point.
    - the right pane shows the last data point average value, per section
    so that it is possible to have a look at this counter while streaming in data
    (displaying various series look messy)
- the status dashboard shows the different sections found in the stream, ordered
    by status code.
    - the right pane shows the global repartition of http requests in the form
    of horizontal bars, orderered by status codes
- finally, the alert dashboard shows a list of enabled alerts with their current
    state: `inactive`, `pending`, `active`.
    - the right pane shows an alert activity log - when an alert is active,
    its name, state and the date of the event are written their.
    - when the same alert recovers, this event is also logged there.

## Important
It is assumed values in the input CSV file are sorted by time in increasing
order. If logs are not sorted by date some features may fail (such as alerting).

## Alerting
Alerts have 3 states (like in prometheus):
- `alert.Inactive`: the alerting rule associated to the alert is false
- `alert.Pending`: the alerting rule is true but for a shorter time than
    `--alert-duration`. This mechanism reduces flappy-alert occurences.
- `alert.Active`: the alerting rule has been true for at least `--alert-period`
    time.

## Log Ingestion and general application flow
Log ingestion is controlled by an ingestor object - it tails a file or stream using
a bounded-size buffer. Read lines are then sent over a channel for a safe asynchronous
read from the parser. That way logs can come and go without locking the app. Lecture is
paused if the buffer is full.

Then registered parsers (only csv at the moment), are run on ingested logs. They generate
a `Trace` structure when parsed successfully. That structure is format-agnostic which
makes it possible for the aggregator to process data from various type of streams. Raw
logs are parsed line by line.

Parsed logs, called `Trace` in the application are then passed to the `MetricsCollector`
component. It generates metrics, aggregating received traces by different criteria.

Generated metrics updates (one by as per the line by line above) are communicated to the
`AlertManager`, evaluating registered alerting rules about metrics. If a metrics is evaluated
for the first time a presentation alert is sent to the view through a channel. It is a presentation
message to let the view know what alerts are going to be displayed without coupling too tightly the two
compenents. Then only changes in alert states are sent over to the view through the channel (i.e
inactive -> pending, pending -> active or active -> inactive).

From other goroutines, the view monitors the emitted alerts and metrics to display them.
Metrics and alert feeds are first set to the view by the `App` object. This object makes
the glue work to bind the backend and frontend.

## Build and run the app locally
The application is coded in `go`. To build this locally you need to install go
`1.21` at least. Do not worry it is also possible to run it with docker.
``` sh
go build -o httpmon main.go
```

Run from a file:
``` sh
./httpmon --file ./sample_csv.txt
```
To quit the app either press `escape` or `ctrl+c`.

Run from stdin:
``` sh
cat sample_csv.txt | ./httpmon --stdin
```

## Show flags and default values
``` sh
./httpmon --help
```

``` sh
Usage of ./httpmon:
  -alert-duration duration
        if requests/s > --threshold for --alert-duration then the alert is active (go duration format) (default 1m0s)
  -alert-threshold uint
        requests/s threshold over wich the alert becomes active (default 10)
  -debug
        wait a few seconds before starting
  -file string
        csv file to read http traces from
  -lines uint
        size of the line buffer when reading logs (default 100)
  -period duration
        log aggregation period used to generate metrics values (go duration format) (default 10s)
  -stdin
        read http logs from stdin, takes precendence over --file
```

## Build and run using docker
``` sh
docker build -t httpmon
docker run -it --rm httpmon
```

By default the docker image reads `sample_csv.txt`. If you want to
pass it another file, mount a volume as follows:
``` sh
docker run -it --rm -v local_full_path:container_full_path httpmon --file container_full_path
```

Finally, the app runs intentionally in an `alpine` image so that it is possible to use a
shell to debug it. Since it is a developper oriented application, I thought it did not
worth it to save extra-weight using `scratch`.

Funny note, when using docker, buttons become green.

## Run the tests
``` sh
go test ./...
```

## Improvements
- Add way more tests
- a CI/CD for gitlab/github
    - go binary: runs `gofmt -l`; go test; 
    docker image: docker build; docker run the tests (to check the image can run)
    - on tagging publishes new version on private artifactory manager for both 
    the go binary and the docker image. The tag is the version, it uses
    `semantic versioning`
- Sync UI with Backend loop so that event management is synchronised
    - I used the UI controller at first to do so, that worked but realised it is not
    thread safe... so crashes occured sometimes.
- Move code from `view.go` to other files by role. Maybe create another package to make
      simpler to read and understand (as well as decoupling elements).
- Filters for the ingestor so that cardinality can be reduced
- Implement more data format parsers (only CSV is supported for now)
- Provide a repartition view of requests per method to see if a section is
being more accessed in reading or writting (which could drive infrastructure
optimisation)
- Add recording rules for more complex metrics, especially the ones working
    on aggregation in order to save resources
- Display alert threshold on associated metric's dashboard
    - add toggle to hide or display the threshold
- a request selector for the request/s dashboard, it would list available
    requests for users to pick the series to display
- read configuration from a config file
- have configurable log aggregation rules to automatise metric generation


## Final note
On a distributed system, only the backend of this application would be of use, exposing
metrics in prometheus format on `/metrics` so that it can be scrapped by prometheus, displayed
and explored with grafana.

For such an exporter it is good practice to expose:
- `/config` to know how the exporter is configured, explore configs globally...
- `/health` and `/ready` so that it can be a good kubernetes citizen and fit well with its
    pod scheduling.

It can also be useful to expose `profiling data`. Though, profiling can be expensive so it
is better to control this by a feature flag. In go it can be done easily with `pprof` which
then exposes `/debug/pprof`. Using `/debug` for debug information focused on a specific
agent is also a good thing to have.

