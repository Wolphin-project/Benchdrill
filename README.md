Commands to run Beedrill:

``` shell
$ make
$ docker build -t beedrill .
$ docker build -t one-mega-nginx -f Dockerfile_nginx .
$ docker stack deploy --compose-file compose.yml lulz
$ ./bin/beedrill
```

You need to add a directory named “nginx-files/” which have to contain a file named “one_mega_file”.