## Test Curl Commands

### Auth

#### Successful Logins

```
curl -X POST -d username=sebvet -d password=astonmartin http://localhost:8000/login -v

Status 200:				{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNlYnZldCIsIm5hbWUiOiIiLCJleHAiOjE2MTU1NTgxMzZ9.dDhv7JpV1HmexRgQMFSH9YJH47nkckgFRWJLSIobdco"}
```

```
curl -X POST -d username=babydriver -d password=edgarwright http://localhost:8000/login -v

Status 200:
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJhYnlkcml2ZXIiLCJuYW1lIjoiIiwiZXhwIjoxNjE1NTU4NzA2fQ.i427TdElCviP8lr08EbqP5_6ocosabSttNIhvEkIa_s"}
```

#### Unsuccessful Login

```
Unsuccessful Login:
curl -X POST -d username=notauser -d password=notapassword http://localhost:8000/login -v

Status 401:
{"error": "Incorrect credentials provided"}
```

### Roster

#### Join Roster Successfully

```
curl -X POST -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNlYnZldCIsIm5hbWUiOiIiLCJleHAiOjE2MTU1NTg2MTN9.6XKKiAC0aL3Dny3OOaD64HL8OU9V34xceaOFjmiR-dU\",\"rate\":5}" http://localhost:8001/roster -v

Status 200:
{"username":"sebvet","name":"Sebastian Vettel","rate":5}
```

```
curl -X POST -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJhYnlkcml2ZXIiLCJuYW1lIjoiIiwiZXhwIjoxNjE1NTU4NzA2fQ.i427TdElCviP8lr08EbqP5_6ocosabSttNIhvEkIa_s\",\"rate\":15}" http://localhost:8001/roster -v

Status 200:
{"username":"babydriver","name":"Ansel Elgort","rate":15}
```

#### Join roster unsuccessfully

```
curl -X POST -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJhYnlkcml2ZXIiLCJuYW1lIjoiIiwiZXhwIjoxNjE1NTU4NzA2fQ.i427TdElCviP8lr08EbqP5_6ocosabSBADSIGNATURE\",\"rate\":15}" http://localhost:8001/roster -v

Status 401:
{"error": "Invalid JWT token"}
```

```
curl -X POST -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNlYnZldCIsIm5hbWUiOiIiLCJleHAiOjE2MTU1NTg2MTN9.6XKKiAC0aL3Dny3OOaD64HL8OU9V34xceaOFjmiR-dU\",\"rate\":5}" http://localhost:8001/roster -v

Status 400:
{"error": "User is already in roster"}
```

#### Change Rate

```
curl -X PUT -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNlYnZldCIsIm5hbWUiOiIiLCJleHAiOjE2MTU1NTg2MTN9.6XKKiAC0aL3Dny3OOaD64HL8OU9V34xceaOFjmiR-dU\",\"rate\":7}" http://localhost:8001/roster -v

Status 200:
{"username":"sebvet","name":"Sebastian Vettel","rate":7}
```

#### Leave Roster

```
curl -X DELETE -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNlYnZldCIsIm5hbWUiOiIiLCJleHAiOjE2MTU1NTg2MTN9.6XKKiAC0aL3Dny3OOaD64HL8OU9V34xceaOFjmiR-dU\"}" http://localhost:8001/roster -v

Status 200
```

```
curl -X DELETE -H "Content-Type: application/json" --data "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNlYnZldCIsIm5hbWUiOiIiLCJleHAiOjE2MTU1NTg2MTN9.6XKKiAC0aL3Dny3OOaD64HL8OU9V34xceaOFjmiR-dU\"}" http://localhost:8001/roster -v

Status 400:
{"error": "User is not in roster"}
```

#### Get Drivers in Roster

```
curl -X GET http://localhost:8001/roster -v

Status 200:
[{"username":"babydriver","name":"Ansel Elgort","rate":15},{"username":"sebvet","name":"Sebastian Vettel","rate":5}]
```

### Directions

```
curl -X GET http://localhost:8002/directions/Exeter/Crediton

Status 200:
{"TotalDistance":14007,"ARoadDistance":13403}
```

### Journey

```
curl -X GET http://localhost:8003/journey/Exeter/Crediton

Status 200:
{"start_point":"Exeter","end_point":"Crediton","total_distance":14007,"a_road_distance":13403,"best_driver":{"username":"sebvet","name":"Sebastian Vettel","rate":5},"cost":280} 	
```

