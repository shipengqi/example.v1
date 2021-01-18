## Requirements

### Database

[MySQL container doc](https://hub.docker.com/_/mysql)

Start a mysql server instance:

```bash
docker run --name blog-mysql -p 33061:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
```

Connect to MySQL from the MySQL command line client:

```bash
docker run -it --rm mysql mysql -hshccdfrh75vm8.hpeswlab.net -uroot -p123456 -P 33061
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
4 rows in set (0.41 sec)

mysql> create database blog;
Query OK, 1 row affected (1.48 sec)
mysql> use blog;
Database changed
mysql> show tables;
Empty set (0.01 sec)

mysql> CREATE TABLE `blog_tag` (
    ->   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    ->   `name` varchar(100) DEFAULT '' COMMENT '',
    ->   `created_at` int(10) unsigned DEFAULT '0' COMMENT '',
    ->   `created_by` varchar(100) DEFAULT '' COMMENT '',
    ->   `updated_at` int(10) unsigned DEFAULT '0' COMMENT '',
    ->   `modified_by` varchar(100) DEFAULT '' COMMENT '',
    ->   `deleted_at` int(10) unsigned DEFAULT '0' COMMENT '',
    ->   `state` tinyint(3) unsigned DEFAULT '1' COMMENT ' 01',
    ->   PRIMARY KEY (`id`)
    -> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';
Query OK, 0 rows affected, 6 warnings (4.10 sec)

mysql> CREATE TABLE `blog_article` (
    ->   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    ->   `tag_id` int(10) unsigned DEFAULT '0' COMMENT 'ID',
    ->   `title` varchar(100) DEFAULT '' COMMENT '',
    ->   `desc` varchar(255) DEFAULT '' COMMENT '',
    ->   `content` text COMMENT '',
    ->   `cover_image_url` varchar(255) DEFAULT '' COMMENT '',
    ->   `created_at` int(10) unsigned DEFAULT '0' COMMENT '',
    ->   `created_by` varchar(100) DEFAULT '' COMMENT '',
    ->   `updated_at` int(10) unsigned DEFAULT '0' COMMENT '',
    ->   `modified_by` varchar(255) DEFAULT '' COMMENT '',
    ->   `deleted_at` int(10) unsigned DEFAULT '0',
    ->   `state` tinyint(3) unsigned DEFAULT '1' COMMENT '',
    ->   PRIMARY KEY (`id`)
    -> ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='';
Query OK, 0 rows affected, 7 warnings (0.11 sec)

```

### Redis

Start a redis instance:
```bash
docker run --name blog-redis -p 6379:6379 -d redis
```

For more about Redis Persistence, see <http://redis.io/topics/persistence>.

connecting via redis-cli:
```bash
docker run -it --rm redis redis-cli -h 16.155.194.49 -p 6379
```

Additionally, If you want to use your own `redis.conf`:
You can create your own Dockerfile that adds a `redis.conf` from the context into `/data/`, like so.
```dockerfile
FROM redis
COPY redis.conf /usr/local/etc/redis/redis.conf
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
```

Alternatively, you can specify something along the same lines with `docker run` options.
```bash
docker run -v /myredis/conf:/usr/local/etc/redis --name blog-redis redis redis-server /usr/local/etc/redis/redis.conf
```
 
Where `/myredis/conf/ `is a local directory containing your `redis.conf` file. Using this method means that there 
is no need for you to have a `Dockerfile` for your redis container.

The mapped directory should be writable, as depending on the configuration and mode of operation, Redis may need 
to create additional configuration files or rewrite existing ones.

## Document

### Swagger

[swaggo/swag official doc](https://github.com/swaggo/swag/blob/master/README_zh-CN.md).

Download swag:
```bash
go get -u github.com/swaggo/swag/cmd/swag

$ swag -v
swag version v1.6.5
```

Install gin-swagger:
```bash
go get -u github.com/swaggo/gin-swagger@v1.2.0 
go get -u github.com/swaggo/files
go get -u github.com/alecthomas/template
```

Run `swag init` in the project's root folder which contains the `main.go` file. 
This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).

Generate: 
```bash
docs/
├── docs.go
└── swagger
    ├── swagger.json
    └── swagger.yaml
```

Run your app, and browse to `http://localhost:8080/swagger/index.html`. 

## Signal

| 命令      | 信号    | 含义                                                                                                    |
| --------- | ------- | ------------------------------------------------------------------------------------------------------- |
| ctrl + c  | SIGINT  | 强制进程结束                                                                                            |
| ctrl + z  | SIGTSTP | 任务中断，进程挂起                                                                                      |
| ctrl + \  | SIGQUIT | 进程结束 和 `dump core`                                                                                 |
| ctrl + d  |         | EOF                                                                                                     |
|           | SIGHUP  | 终止收到该信号的进程。若程序中没有捕捉该信号，当收到该信号时，进程就会退出（常用于 重启、重新加载进程） |

### All signals

```
$ kill -l
 1) SIGHUP   2) SIGINT   3) SIGQUIT  4) SIGILL   5) SIGTRAP
 6) SIGABRT  7) SIGBUS   8) SIGFPE   9) SIGKILL 10) SIGUSR1
11) SIGSEGV 12) SIGUSR2 13) SIGPIPE 14) SIGALRM 15) SIGTERM
16) SIGSTKFLT   17) SIGCHLD 18) SIGCONT 19) SIGSTOP 20) SIGTSTP
21) SIGTTIN 22) SIGTTOU 23) SIGURG  24) SIGXCPU 25) SIGXFSZ
26) SIGVTALRM   27) SIGPROF 28) SIGWINCH    29) SIGIO   30) SIGPWR
31) SIGSYS  34) SIGRTMIN    35) SIGRTMIN+1  36) SIGRTMIN+2  37) SIGRTMIN+3
38) SIGRTMIN+4  39) SIGRTMIN+5  40) SIGRTMIN+6  41) SIGRTMIN+7  42) SIGRTMIN+8
43) SIGRTMIN+9  44) SIGRTMIN+10 45) SIGRTMIN+11 46) SIGRTMIN+12 47) SIGRTMIN+13
48) SIGRTMIN+14 49) SIGRTMIN+15 50) SIGRTMAX-14 51) SIGRTMAX-13 52) SIGRTMAX-12
53) SIGRTMAX-11 54) SIGRTMAX-10 55) SIGRTMAX-9  56) SIGRTMAX-8  57) SIGRTMAX-7
58) SIGRTMAX-6  59) SIGRTMAX-5  60) SIGRTMAX-4  61) SIGRTMAX-3  62) SIGRTMAX-2
63) SIGRTMAX-1  64) SIGRTMAX
```
