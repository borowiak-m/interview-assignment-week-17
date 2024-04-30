Create a RESTFUL API with two endpoints.

- One of them that fetches the data in the provided MongoDB collection and returns the results in the requested format.
- Second endpoint is to create (POST) and fetch(GET) data from an in-memory database.

No frameworks allowed.
# DELIVERABLES
public repo with clear instructions on configuration and running the application locally.

## requirements

mongodb+srv://challengeUser:WUMgIwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true
/GET from Mongo
1. request payload:
---JSON: 4 fields
startDate date
endDate date
minCount int
maxCount int

2. response payload:
---JSON: 3 fields
code: status code (0 for success)
msg: status description (success)
records: array of filtered items including fields 
- key
- createdAt
- totalCount = sum of count in the document/table

/GET /POST from/to array 
1. POST request payload:
---JSON: 2 fields
key string
value string
2. POST response payload:
echo of request
3. GET request payload:
key param in query param
sample: localhost/in-memory?key=active-tabs
4. GET response payload
---JSON 2 fields
key
value