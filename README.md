# Benchdrill

Benchdrill is a benchmarking tool based on [Machinery (revision fdcbe0f of 31/5/2017)](https://github.com/RichardKnop/machinery/tree/fdcbe0ff6b8592b8ebc65f76fd63af3f6f90c3a7). It allows to load charges thanks to Machinery’s workers, which are distributed with Docker Compose in a Docker Swarm. Workers are then ready to execute tasks to run tools as [Filebench](https://github.com/filebench/filebench).

## Installation
First you need Go 1.8 and Docker v17.06. Then, after cloning this repository with Git, run the following commands in a terminal in the root directory of Benchdrill:

``` shell
$ docker build -t hub.rnd.alterway.fr/wolphin-project/benchdrill:master .
$ ./stack/benchdrill-network
$ ./stack/benchdrill-deploy
```

## Usage
Benchdrill commands are run with `stack/benchdrill-cli`. Currently there are 2 commands with Benchdrill: `send_cmd_args` and `send_cmd_file`. The first one send a command to a worker, possibly with arguments, the second one allows to send the content of a local file given as an argument to a worker, which will save the content on a file then use it.

Two benchmark tools are currently supported: [Sysbench](https://github.com/akopytov/sysbench) and [Filebench](https://github.com/filebench/filebench). Below two commands you can run, given as examples.

``` shell
$ ./stack/benchdrill-cli send_cmd_args "sysbench --time=5 cpu run"
```

This command will send the command quoted to a single worker, which will run a built-in CPU test of Sysbench for 5 seconds.

``` shell
$ ./stack/benchdrill-cli --times 3 send_cmd_file "filebench -f" < readfiles.f
```

This command will send the command quoted 3 times as 3 separate and identical tasks executed simultaneously (thanks to the ``--times`` option); it will run the test written in `readfiles.f` (provided in the repository as an example).

## Architecture

![Architecture schema of Benchdrill](architecture_schema.png)

Both Redis and Worker services are in a Docker Swarm. A Swarm is a cluster in which each node is a Docker engine. A node can be a manager or a worker. In Swarm mode concepts are not about containers but about services. A service is the a task a manager or a worker has to execute. It can be several identical containers distributed across nodes. In our case, the Redis service is composed of a unique container, but the Worker service has several identical containers and can be scaled up or down with a `docker stack` command. More information about Docker Swarm [here](https://docs.docker.com/engine/swarm/).

## License
Mozilla Public License 2.0
