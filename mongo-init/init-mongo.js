db = db.getSiblingDB('maindatabase');

db.createCollection('records');
db.records.createIndex({ createdAt: 1}); //indexing for eye-watering speed
// dummy data
db.records.insertMany([
    {
    "key": "TAKwGc6Jr4i8Z487",
    "createdAt": ISODate("2017-01-28T01:22:14.398Z"),
    "count": [500, 400, 450, 550, 300, 150, 350]
    },
    {
    "key": "NAeQ8eX7e5TEg70H",
    "createdAt": ISODate("2017-01-27T08:19:14.135Z"),
    "count": [540, 400, 450, 550, 300, 160, 350]
    },
    {
    "key": "cCddT2RPqWmUI4Nf",
    "createdAt": ISODate("2017-01-27T13:22:10.421Z"),
    "count": [120, 400, 450, 660, 500, 770, 250]
    }
    ])