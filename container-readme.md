
Build image:
`docker buildx build -t wip-1 .`

Run image as named container:
`docker container run --interactive --tty --name container-name-1 --hostname hostname-1 wip-1`

Start container again (required for persistance):
`docker container start --interactive --attach container-name-1`

Build container and run as bash:
`docker compose run --build go-dev`

Cleanup container after exit:
`docker compose run --rm go-dev`

Build and remove (make sure to rebuild when changing code):
`docker compose run --build --rm go-dev`

When using volumes to persist source code, make sure you grant permissions to the folder on the host.

https://stackoverflow.com/questions/47197493/docker-mounting-volume-permission-denied

The question title does not reflect the real problem in my opinion.

mkdir /srv/redis/redisTest
mkdir: cannot create directory ‘/srv/redis/redisTest’: Permission denied
This problem occurs very likely because when you run:

docker run -d -v /srv/redis:/data --name myredis redis
the directory /srv/redis ownership changes to root. You can check that by

ls -lah /srv/redis
This is normal consequence of mounting external directory to docker. To regain access you have to run

sudo chown -R $USER /srv/redis

# Permissions

If you get a permission error inside the container when trying to run docker--You probably haven't granted permissions to the docker socket on the host machine. i.e.
```
permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get "http://%2Fvar%2Frun%2Fdocker.sock/v1.24/containers/json?all=1&filters=%7B%22label%22%3A%7B%22com.docker.compose.config-hash%22%3Atrue%2C%22com.docker.compose.project%3Dmy-app%22%3Atrue%7D%7D": dial unix /var/run/docker.sock: connect: permission denied
```

Change ownership:

`sudo chown -R 1000:1000 /var/run/docker.sock`

inside your container, check

ls -lh

# Multiple user accounts

When adding a second user account I ran into several issues with docker. 

1. Add second user to `docker-users` group. Windows key -> Computer Management -> Local Users and Groups -> docker-users -> Add ...

2. Set docker container mode to linux to work in WSL2. Docker icon in system tray -> right click -> Windows/linux mode 
