## batch_proc
The application requires GO version `1.14` or higher to run.

### Running it

#### Download release
Download a pre-compiled binary from the [releases](https://github.com/rolandvarga/batch_proc/releases/) section of this repository.

The following systems are currently supported:
* Linux
* Mac

#### via Docker
This approach requires you to have Docker [installed](https://docs.docker.com/get-docker/) and running.

```bash
# use Make to build the container
make docker

# alternatively you COULD run the build command yourself
docker build -t rolandvarga/batch_proc .

# the build process will end showing you a tag for the container:
# Successfully built ed6c6fb2582c
# Successfully tagged rolandvarga/batch_proc:HEAD

# use the tag created by the build command and enter the container:
docker run -it --entrypoint /bin/sh rolandvarga/batch_proc:HEAD

# inside the container
cd /go/src/github.com/rolandvarga/batch_proc
./batch_proc data.json
```

#### build locally
```bash
make all
./batch_proc data.json
```

The processing step of the program prints out each batch as soon as a sequence is 'complete'. For the current data set it is advised to redirect all output to a file for ease of inspection e.g.:
```bash
./batch_proc data.json > out
```
