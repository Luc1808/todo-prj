#!/usr/bin/env bash

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imx1Y2FzQGdtYWlsLmNvbSIsImV4cCI6MTc0NTQwMzg5MSwidXNlcl9pZCI6MX0.I6wmW4SmQ_LA9GksBjMcd8jSf2-5qVe2Tjvv9RQROA0"

# 1
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Buy milk, bread, and eggs","complete":false,"priority":"Medium","dueAt":"2025-05-01T10:00:00Z","category":"Health"}'

# 2
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Submit project report","description":"Finalize and submit the project report to manager","complete":false,"priority":"High","dueAt":"2025-04-25T17:00:00Z","category":"Finance"}'

# 3
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Workout session","description":"Attend gym for strength training","complete":false,"priority":"Low","dueAt":"2025-04-24T07:30:00Z","category":"Self Development"}'

# 4
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Call mom","description":"Check in with mom to see how she’s doing","complete":false,"priority":"Low","dueAt":"2025-04-24T20:00:00Z","category":"Social"}'

# 5
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Read book","description":"Read at least 50 pages of current book","complete":false,"priority":"Low","dueAt":"2025-04-26T21:00:00Z","category":"Self Development"}'

# 6
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Pay bills","description":"Pay electricity and internet bills","complete":false,"priority":"High","dueAt":"2025-04-24T12:00:00Z","category":"Finance"}'

# 7
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Plan team meeting","description":"Organize agenda for next team sync","complete":false,"priority":"Medium","dueAt":"2025-04-27T09:00:00Z","category":"Social"}'

# 8
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Dentist appointment","description":"Routine dental check-up","complete":false,"priority":"Medium","dueAt":"2025-05-02T11:00:00Z","category":"Health"}'

# 9
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Meditation","description":"Meditate for 15 minutes","complete":false,"priority":"Low","dueAt":"2025-04-24T06:30:00Z","category":"Self Development"}'

# 10
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Invest research","description":"Research potential stocks for investment","complete":false,"priority":"Medium","dueAt":"2025-04-28T14:00:00Z","category":"Finance"}'

# 11
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Friend’s birthday gift","description":"Buy a gift for Alex’s birthday","complete":false,"priority":"Low","dueAt":"2025-04-30T18:00:00Z","category":"Social"}'

# 12
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Yoga class","description":"Attend virtual yoga session","complete":false,"priority":"Low","dueAt":"2025-04-25T19:00:00Z","category":"Health"}'

# 13
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Update resume","description":"Add recent projects to resume","complete":false,"priority":"Medium","dueAt":"2025-04-29T16:00:00Z","category":"Self Development"}'

# 14
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Review budget","description":"Analyze monthly expenses and adjust budget","complete":false,"priority":"Medium","dueAt":"2025-04-30T12:00:00Z","category":"Finance"}'

# 15
curl -s -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Coffee with Sarah","description":"Catch up over coffee","complete":false,"priority":"Low","dueAt":"2025-04-26T10:30:00Z","category":"Social"}'

echo "Done sending 15 dummy tasks."

