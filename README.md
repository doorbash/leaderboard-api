## Prerequisites
- Go
- Docker
- Docker Compose

## Configuration
- Put SSL key files in:
```
docker/nginx/cert/fullchain.pem
docker/nginx/cert/privkey.pem
```

- Create `.env`:
```
APP_VERSION=1.0.2

DATABASE_ROOT_PASSWORD="DATABASE_ROOT_PASSWORD"
DATABASE_USER="DATABASE_USER"
DATABASE_PASSWORD="DATABASE_PASSWORD"
DATABASE_NAME="DATABASE_NAME"

PMA_URI="https://your.domain/pma/"

IMAGE_MYSQL="mysql:8.0.30"
IMAGE_PMA="phpmyadmin/phpmyadmin:5.2.0"
IMAGE_NGINX="nginx:1.23.1-alpine"
IMAGE_REDIS="redis:7.0.4-alpine3.16"
```

## Run
```
./run.sh prod
```

## Postman
https://documenter.getpostman.com/view/13117984/VUqpscfn

## License
MIT