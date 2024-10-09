# Geographic Data REST API Service
A RESTful API service built in Go (Golang) with PostgreSQL for handling geographic data—specifically points and contours. The service supports CRUD operations and advanced spatial queries, including intersection calculations and point-in-polygon queries.

Please refer to the architecture doc for the comprehensive explanation of the project structure [ARCHITECTURE.md](./ARCHITECTURE.md) 

## Prerequisites

- **Go**: Version 1.16 or higher.
- **Docker** and **Docker Compose**: For containerization.
- **PostgreSQL**: With PostGIS extension.

The project uses `make` to make your life easier. If you're not familiar with Makefiles you can take a look at [this quickstart guide](https://makefiletutorial.com).
Whenever you need help regarding the available actions, just use the following command.

```bash
make help
```

## Setup

To get your setup up and running the only thing you have to do is

```bash
make all
```

This will initialize a git repo, download the dependencies in the latest versions and install all needed tools.
If needed code generation will be triggered in this target as well.

## Test & lint

Run linting

```bash
make lint
```

Run unit tests

```bash
make test-unit
```

Run repository tests

The tests for repository is separated with actual connection to database, please make sure you have the database running `make start-components`

```bash
make test-repository
```

## Local Development

### Run dependency

```bash
make start-components
```

### Stop dependency

```bash
make stop-components 
```

### Build docker image

```bash
make docker-build
```

### Run the server localy

Copy the `.env.sample` to `.env` and put the appropriate value for your setup, then run:

```bash
make server
```

### Sample API 

#### Create Points

Request

```bash
curl --location 'localhost:8080/points' \
--header 'Content-Type: application/json' \
--data '{
    "data": {
        "type": "Point",
        "coordinates": [
            1.0,
            2.0
        ]
    }
}'
```

Response

```json
{
    "id": 22,
    "data": {
        "type": "Point",
        "coordinates": [
            1,
            2
        ]
    }
}
```

#### Get Points

Request

```bash
curl --location --request GET 'localhost:8080/points'
```

Response

```json
{
    "count": 7,
    "next": "http://localhost:8080/points?page=1",
    "previous": null,
    "results": [
        {
            "id": 28,
            "data": {
                "type": "Point",
                "coordinates": [
                    17,
                    17
                ]
            }
        }
    ]
}
```

#### Get Contours

Request

```bash
curl --location 'localhost:8080/contours'
```

Response

```json
{
    "count": 10,
    "next": "http://localhost:8080/contours?page=1",
    "previous": null,
    "results": [
        {
            "id": 11,
            "data": {
                "type": "Polygon",
                "coordinates": [
                    [
                        [
                            1,
                            1
                        ],
                        [
                            3,
                            1
                        ],
                        [
                            3,
                            3
                        ],
                        [
                            1,
                            3
                        ],
                        [
                            1,
                            1
                        ]
                    ]
                ]
            }
        }
    ]
}
```




#### Create Contours

Request

```bash
curl --location 'localhost:8080/contours' \
--header 'Content-Type: application/json' \
--data '{
    "data": {
        "type": "Polygon",
        "coordinates": [
            [
                [
                    30.0,
                    10.0
                ],
                [
                    40.0,
                    40.0
                ],
                [
                    20.0,
                    40.0
                ],
                [
                    10.0,
                    20.0
                ],
                [
                    30.0,
                    10.0
                ]
            ]
        ]
    }
}'
```

Response

```json
{
    "id": 43,
    "data": {
        "type": "Polygon",
        "coordinates": [
            [
                [
                    30,
                    10
                ],
                [
                    40,
                    40
                ],
                [
                    20,
                    40
                ],
                [
                    10,
                    20
                ],
                [
                    30,
                    10
                ]
            ]
        ]
    }
}
```

#### Update Contours

Request

```bash
curl --location --request PUT 'localhost:8080/contours/43' \
--header 'Content-Type: application/json' \
--data '{
    "data": {
        "type": "Polygon",
        "coordinates": [
            [
                [
                    30.0,
                    10.0
                ],
                [
                    40.0,
                    40.0
                ],
                [
                    20.1,
                    40.0
                ],
                [
                    10.0,
                    20.1
                ],
                [
                    30.0,
                    10.0
                ]
            ]
        ]
    }
}'
```

Response

```json
{
    "id": 43,
    "data": {
        "type": "Polygon",
        "coordinates": [
            [
                [
                    30,
                    10
                ],
                [
                    40,
                    40
                ],
                [
                    20.1,
                    40
                ],
                [
                    10,
                    20.1
                ],
                [
                    30,
                    10
                ]
            ]
        ]
    }
}
```

#### Get Contours By ID

Request

```bash
curl --location 'localhost:8080/contours/9'
```

Response

```json
{
    "id": 9,
    "data": {
        "type": "Polygon",
        "coordinates": [
            [
                [
                    25,
                    25
                ],
                [
                    45,
                    5
                ],
                [
                    25,
                    5
                ],
                [
                    25,
                    25
                ]
            ]
        ]
    }
}
```

#### Delete Contours

Request

```bash
curl --location --request DELETE 'localhost:8080/contours/43'
```

Response

```
204 No Content
```

#### Get Points By Contour

Request

```bash
curl --location 'localhost:8080/points?contour=2'
```

Response

```json
{
    "count": 1,
    "next": "http://localhost:8080/points?page=1",
    "previous": null,
    "results": [
        {
            "id": 28,
            "data": {
                "type": "Point",
                "coordinates": [
                    17,
                    17
                ]
            }
        }
    ]
}
```

#### Get Contours Intersections Area

Request

```bash
curl --location 'localhost:8080/intersections?contour_1=1&contour_2=2'
```

Response

```json
[
    {
        "data": {
            "type": "Polygon",
            "coordinates": [
                [
                    [
                        15,
                        20
                    ],
                    [
                        20,
                        20
                    ],
                    [
                        20,
                        15
                    ],
                    [
                        15,
                        15
                    ],
                    [
                        15,
                        20
                    ]
                ]
            ]
        }
    }
]
```