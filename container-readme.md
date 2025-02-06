
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
