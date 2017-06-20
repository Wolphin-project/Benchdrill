Commands to set up Beedrill:

``` shell
$ make
$ docker build -t hub.rnd.alterway.fr/wolphin-project/beedrill:master .
$ ./stack/beedrill-network
$ ./stack/beedrill-deploy
```

Example commands to use Beedrill:

``` shell
$ ./stack/beedrill-cli send_cmd_args "sysbench --help"
$ ./stack/beedrill-cli --times=3 send_cmd_file "filebench -f" < readfiles.f
```