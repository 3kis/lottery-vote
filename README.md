投票系统

# 查询票数接口
```cmd
curl --location --request POST '127.0.0.1:8080/graphql' \
--header 'Content-Type: application/json' \
--data-raw '{
  "query": " { getVoteCount(username: \"alice\") { username, votes } }"
}'
```

# 投票接口
```cmd
curl -X POST 'http://localhost:8080/graphql' \
-H 'Content-Type: application/json' \
-d '{
"query": "mutation { vote(ticket: \"sample-ticket\", usernames: [\"alice\", \"bob\"]) { success } }"
}'
```

# 获取ticket接口
```cmd
curl --location --request POST '127.0.0.1:8080/graphql' \
--header 'Content-Type: application/json' \
--data-raw '{
"query": "{ getTicket }"
}'
```
