# KNoT Babeltower

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/9140aa8c06934071ad6e3cf3b1b148ff)](https://www.codacy.com/manual/joaoaneto/knot-babeltower?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=CESARBR/knot-babeltower&amp;utm_campaign=Badge_Grade)

## Contents

- [Basic installation and usage](#basic-installation-and-usage)
  - [Requirements](#requirements)
  - [Configuration](#configuration)
  - [Setup](#setup)
  - [Compiling and Running](#compiling-and-running)
- [Docker installation and usage](#docker-installation-and-usage)
  - [Requirements](#requirements)
  - [Building and Running](#building-and-running)
- [Verify service health](#verify-service-health)

## Basic installation and usage

### Requirements

*   Install glide dependency manager (<https://github.com/Masterminds/glide>)
*   Be sure the local packages binaries path is in the system's `PATH` environment variable:

```bash
$ export PATH=$PATH:<your_go_workspace>/bin
```

### Configuration

You can set the `ENV` environment variable to `development` and update the `internal/config/development.yaml` when necessary. On the other way, you can use environment variables to configure your installation. In case you are running the published Docker image, you'll need to stick with the environment variables.

The configuration parameters are the following (the environment variable name is in parenthesis):

*   `server`
    *   `port` (`SERVER_PORT`) **Number** Server port number. (Default: 80)

### Setup

```bash
make tools
make deps
```

### Compiling and running

```bash
make run
```

> You can use the `make watch` command to run the application on watching mode, allowing it to be restarted automatically when the code changes.

## Docker installation and usage

### Requirements

*   Install docker engine (<https://docs.docker.com/install/>)

### Building and Running

A container is specified at `docker/Dockerfile`. To use it, execute the following steps:

01. Build the image:

    ```bash
    docker build . -f docker/Dockerfile -t cesarbr/knot-babeltower
    ```

01. Create a file containing the configuration as environment variables.

01. Run the container:

    ```bash
    docker run --env-file knot-babeltower.env -ti cesarbr/knot-babeltower
    ```

## Verify service health

```bash
curl http://<hostname>:<port>/healthcheck
```
