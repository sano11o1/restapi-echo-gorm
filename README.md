DBに入る
```
docker compose exec db bash
mysql -uroot -ppass
use go_mysql8_development 
```

マイグレーション実行
```
migrate -database "mysql://webuser:webpass@tcp(localhost:3306)/go_mysql8_development" -source "file://database/migrations/" up
```
