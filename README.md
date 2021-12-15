Compose Switch 
--------------

Compose Switch is a replacement to the Compose V1 `docker-compose` (python) executable. It translates the command line into Compose V2 `docker compose` then run the latter. 

## installation

We provide an script for automated installation:

```console
$ curl -fL https://raw.githubusercontent.com/docker/compose-switch/master/install_on_linux.sh | sh
```

### Manual installation

1. download compose-switch binary for your architecture
```console
$ curl -fL https://github.com/docker/compose-switch/releases/latest/download/docker-compose-linux-amd64 -o /usr/local/bin/compose-switch
```
2. make compose-switch executable
```console
$ chmod +x /usr/local/bin/compose-switch
```
3. rename `docker-compose` binary if you already have it installed as `/usr/local/bin/docker-compose` 
```console
$ mv /usr/local/bin/docker-compose /usr/local/bin/docker-compose-v1
```
4. define an "alternatives" group for `docker-compose` command: 
```console
$ update-alternatives --install /usr/local/bin/docker-compose docker-compose <PATH_TO_DOCKER_COMPOSE_V1> 1
$ update-alternatives --install /usr/local/bin/docker-compose docker-compose /usr/local/bin/compose-switch 99
```

## check installation

```console
$ update-alternatives --display docker-compose
docker-compose - auto mode
  link best version is /usr/local/bin/compose-switch
  link currently points to /usr/local/bin/compose-switch
  link docker-compose is /usr/local/bin/docker-compose
/usr/bin/docker-compose - priority 1
/usr/local/bin/compose-switch - priority 99
```

## select Compose implementation to run by `docker-compose`

```console
$ update-alternatives --config docker-compose
There are 2 choices for the alternative docker-compose (providing /usr/local/bin/docker-compose).

  Selection    Path                           Priority   Status
------------------------------------------------------------
* 0            /usr/local/bin/compose-switch   99        auto mode
  1            /usr/bin/docker-compose         1         manual mode
  2            /usr/local/bin/compose-switch   99        manual mode

Press <enter> to keep the current choice[*], or type selection number: 
```
