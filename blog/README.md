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

## Document

### Swagger

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
