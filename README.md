##  How It Works
**_docker-compose up --build_ command will be enough to run the project.**
**you have to wait about 10 second for the database to be ready.**

```bash
    docker-compose up --build
```

**or you can run this api _make run_ command.**

---

**This API is written using hexagonal architecture and consists of 5 endpoint designed to fulfill the following
requirements of the case. This repo included unit tests and integration tests. Below you can see this endpoint and the
specific topics of the data they provide:**

#### CreateLocation _(it creates a location)_
This endpoint creates a location.

**REQUEST**
```bash 
  curl --location 'http://localhost:96/location' \
    --header 'Content-Type: application/json' \
    --data '{
        "name": "test1",
        "latitude": 124.1222,
        "longitude": 134.1222,
        "marker_color": "FFFAFF"
    }'
```
**200 - response**
```json
{
  "id":"67d6ba9821e5359a8b2ebb26"
}
```
**400 - response**
```json
{
  "error":"Invalid request"
}
```

#### GetLocation _(it returns a location using id)_
This endpoint returns a location using id.

**REQUEST**
```bash 
  curl --location 'http://localhost:96/location?id=67d56bb1634ce74585317d40'
```
**200 - response**
```json
{
  "id": "67d6bd8821e5359a8b2ebb27",
  "name": "test1",
  "latitude": 124.1222,
  "longitude": 134.1222,
  "marker_color": "FFFAFF"
}
```
**500 - response**
```json
{
  "error": "mongo: no documents in result"
}
```

#### GetLocations _(it returns locations)_
This endpoint returns locations. You can page and limit options to get locations.

**REQUEST**
```bash 
  curl --location 'http://localhost:96/locations?page=1&limit=3'
```
**200 - response**
```json
{
  "locations":[
    {"id":"67d6ba8c21e5359a8b2ebb25","name":"test3","latitude":124.12,"longitude":134.12,"marker_color":"FFFAFF"},
    {"id":"67d6ba9821e5359a8b2ebb26","name":"test1","latitude":124.1222,"longitude":134.1222,"marker_color":"FFFAFF"},
    {"id":"67d6bd8821e5359a8b2ebb27","name":"test1","latitude":124.1222,"longitude":134.1222,"marker_color":"FFFAFF"}
  ]
}
```
**500 - response**
```json
{
  "error": "mongo: no documents in result"
}
```

#### UpdateLocations _(it can update locations)_
This endpoint updates one or more locations using a json body which is an array.

**REQUEST**
```bash 
  curl --location --request PATCH 'http://localhost:96/locations' \
    --header 'Content-Type: application/json' \
    --data '{
        "locations": [
            {
                "id": "67d6ba8c21e5359a8b2ebb25",
                "name": "test",
                "latitude": 1224.21332,
                "longitude": 13344.13243532,
                "marker_color": "FFFBFF"
            },
            {
                "id": "67d562e3d955d225ca4d9918",
                "name": "test",
                "latitude": 1224.1332,
                "longitude": 13344.13243532,
                "marker_color": "FFFBFF"
            }
        ]
    }
    '
```
**206 - response**
```json
{
  "updated_ids":["67d6ba8c21e5359a8b2ebb25"],
  "failed_ids":["67d562e3d955d225ca4d9918"],
  "updated_count":1
}
```

#### GetRoutes _(it returns a list of routes that are sorted by distance)_
This endpoint returns a list of routes that are sorted by distance.

**REQUEST**
```bash 
  curl --location 'http://localhost:96/routes?latitude=123.123&longitude=123.123'
```
**200 - response**
```json
{
  "routes":[
    {"id":"67d6ba9821e5359a8b2ebb26","name":"test1","distance":685.5056322379919,"marker_color":"FFFAFF"},
    {"id":"67d6bd8821e5359a8b2ebb27","name":"test1","distance":685.5056322379919,"marker_color":"FFFAFF"},
    {"id":"67d6ba8c21e5359a8b2ebb25","name":"test","distance":7242.741611883656,"marker_color":"FFFBFF"}
  ]
}
```
**400 - response**
```json
{
  "error":"Key: 'GetRoutesRequest.Longitude' Error:Field validation for 'Longitude' failed on the 'required' tag"
}
```