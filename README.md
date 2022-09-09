## Prerequisites
- Go (only for developement)
- Docker
- Docker Compose

## Configuration
- Create `.env`:
```
APP_VERSION=1.1.1

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
start (develoment)  : ./run.sh
start (production)  : ./run.sh prod
stop                : ./run.sh stop
clean database      : ./run.sh clean
```

## How to

### Create a new Game
1. Go to `phpMyAdmin`->`leaderboard`->`games`
2. Click `Insert`
3. Set `id` (lowercase, no spaces), `name` and `data` (json)
4. Click `Go`

### Create a new Leaderboard
1. Go to `phpMyAdmin`->`leaderboard`->`leaderboards`
2. Click `Insert`
3. Select `gid`
4. Set `name` (whatever you want. doesn't really matter)
5. Set `value1_order`, `value2_order` and `value3_order` based on what you need: <br />

```
 0 :    unused
 1 :    higher is better (ex: score, num levels,...)
-1 :    lower is beter (ex: time)
```

The leaderboard will be sorted later based on `value3_order` (if not zero).

The rows having equal `value3`, will be sorted based on `value2_order` (if not zero)

The rows having equal `value2`, will be sorted based on `value1_order` (if not zero)

6. Set `value1_limit`, `value2_limit`, `value3_limit`

Any player who exceeds these limits will get banned and also their IP address get banned for a week. (0 means no limit)

### Get Game
**GET** `https://your.domain/api/awesome-game`

First result is from database and then for the next 20 minutes results are from cache.

### Get Leaderboard
**GET** `https://your.domain/api/leaderboards/put_leaderboard_uid_here?offset=0&count=50`

To get data from cache set `offset` to 0 and `count` to 100 (or not set them)

Cached results are 0 ~ 20 minutes old but they are faster.

### Get Leaderboard by Player UID
**GET** `https://your.domain/api/leaderboards/put_leaderboard_uid_here/put_player_uid_here?offset=1&count=3`

### Set new leaderboard data
**POST** `https://your.domain/api/leaderboards/put_leaderboard_uid_here/new`

**BODY**
```
{
    "pid": "put_player_uid_here",
    "name": "John",
    "value1": 400,
    "value2": 3,
    "value3": 1
}
```

Not setting `pid` (player uid) means it's a new player.
The `pid` will be returned in the response.

If `pid` is set, `name` can be set or not. If `name` is set then it gets updated in database if it's changed.

Values only get updated if they are better than what already exists in database. For example if a player has finished the game in 300 seconds before now only values less than 300 are accepted. 

If there are limits set for values, any player who exceed those limits gets banned as noted before.

## Postman
https://documenter.getpostman.com/view/13117984/VUqpscfn

## License
MIT