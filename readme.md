# Scratch 

## How to run the application

Clone the application with `git clone https://github.com/toniastro/scratchpay.git`


### Docker

```
cd scratchpay
docker build -t scratch-challenge .
docker run -p 9000:9000 -it scratch-challenge`
```

### Locally    

To run the application locally, please use the following command -

    ```shell
        go run ./
    ```
> Note: By default the port number its being run on is **9000**.

## Endpoints Description

### Get All Clinics

```
    URL - *http://localhost:9000/
    Method - GET
```

### Search (GET Request)

```JSON
    "URL" - *http://localhost:9000?name=Cleveland Clinic
    Method - GET
```

### Search (POST Request)

```JSON
    "URL" - *http://localhost:9000/
    Method - POST
    Body - (content-type = application/json)
    {
        "name":"Cleveland Clinic",
        "state": "California",
        "availability":{
            "to": "22:00"
        }
    }
```

## Test Driven Development Description

To run the unit test cases, please do the following -

1. `go run ./`
2. `go test `
