# go-video-service

### Configure service

Create `.env` file:

```
HTTP_HOST="localhost"
HTTP_PORT=8080
HTTP_READ_TIMEOUT="10s"
HTTP_WRITE_TIMEOUT="10s"
HTTP_MAX_HEADER_MBYTES=1

MINIO_ENDPOINT="localhost:9000"
MINIO_ACCESS_KEY_ID="access-key-id"
MINIO_SECRET_ACCESS_KEY="secret-access-key"
MINIO_ENABLE_SSL=false
```

or create `configs/deployment-config.yaml` file:

```yaml
http:
  host:  "localhost"
  port:  8080
  read_timeout:  10s
  write_timeout:  10s
  max_header_mbytes:  10

minio:
  endpoint:  "localhost:9000"
  access_key_id:  "access-key-id"
  secret_access_key:  "secret-access-key"
  enable_ssl:  false
```

### API documentation
----

##### Health check
----
* **URL**
    `/health-check`
* **Method**
    `GET`
* **Success Response**
    *Status:* 204 (No Content)
    *Headers:*
    ```http
    Content-Type: application/json
    ```

##### Upload video
----
* **URL**
    `/api/video`
* **Method**
    `POST`
* **Headers**
    ```http
    X-Original-Name: filename.extension
    ```
* **Data params**
    Required: `binary`
* **Success response**
    *Status:* 201 (Created)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "uuid": "00000000-0000-0000-0000-000000000000",
        "links": [
            {
                "method": "GET",
                "url": "/api/video/00000000-0000-0000-0000-000000000000"
            },
            {
                "method": "GET",
                "url": "/api/video/00000000-0000-0000-0000-000000000000/stream"
            },
            {
                "method": "DELETE",
                "url": "/api/video/00000000-0000-0000-0000-000000000000"
            }
       ]
    }
    ```
* **Error response**
  * *Status:* 400 (Bad Request)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 500 (Internal Server Error)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```

##### Find video
----
* **URL**
    `/api/video/:uuid`
* **Method**
    `GET`
* **URL params**
    Required: `uuid=[Version 4 UUID]`
* **Success response**
    *Status:* 200 (OK)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "uuid": "00000000-0000-0000-0000-000000000000",
        "original_name": "filename.extension",
        "size": 0,
        "content_type": "video/*",
        "created_at": "0001-01-01T00:00:00Z",
        "links": [
            {
                "method": "GET",
                "url": "/api/video/00000000-0000-0000-0000-000000000000/stream"
            },
            {
                "method": "DELETE",
                "url": "/api/video/00000000-0000-0000-0000-000000000000"
            }
       ]
    }
    ```
* **Error response**
  * *Status:* 400 (Bad Request)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 404 (Not Found)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 500 (Internal Server Error)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```

##### Stream video
----
* **URL**
    `/api/video/:uuid/stream`
* **Method**
    `GET`
* **Headers**
    ```http
    Range: bytes=0-
    ```
* **URL params**
    Required: `uuid=[Version 4 UUID]`
* **Success response**
  * *Status:* 206 (Partial Content)
    *Headers:*
    ```http
    Accept-Ranges: bytes
    Content-Type: application/octet-stream
    Content-Range: bytes 0-1048576/2097152
    Content-Length: 1048576
    ```
    *Content:* `bytes`
  THEN
  * *Status:* 200 (OK)
    *Headers:*
    ```http
    Accept-Ranges: bytes
    Content-Type: application/octet-stream
    Content-Range: bytes 2097152-2097152/2097152
    Content-Length: 0
    ```
    *Content:* `bytes`
* **Error response**
  * *Status:* 400 (Bad Request)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 404 (Not Found)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 500 (Internal Server Error)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```

##### Remove video
----
* **URL**
    `/api/video/:uuid`
* **Method**
    `DELETE`
* **URL params**
    Required: `uuid=[Version 4 UUID]`
* **Success response**
    *Status:* 204 (No Content)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
* **Error response**
  * *Status:* 400 (Bad Request)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 404 (Not Found)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
  OR
  * *Status:* 500 (Internal Server Error)
    *Headers:*
    ```http
    Content-Type: application/json
    ```
    *Content:*
    ```json
    {
        "message": ""
    }
    ```
