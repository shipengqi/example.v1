## READY

### Database

[MySQL container doc](https://hub.docker.com/_/mysql)

Start a mysql server instance:

```bash
docker run --name blog-mysql -p 33061:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
```

Connect to MySQL from the MySQL command line client:

```bash
docker run -it --rm mysql mysql -hshccdfrh75vm8.hpeswlab.net -uroot -p123456 -P 33061
```
