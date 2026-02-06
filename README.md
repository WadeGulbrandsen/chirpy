![Chirpy Logo](assets/logo.png)

# Chirpy

Repo for boot.dev [learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) course.

Chirpy is a basic X (formally Twitter) style API where users can post Chirps.

## Features
- API can manage users and chirps
- Serves the a directory as a static web site along side the API
- Users can use the (fictional) payment processor, **Polka**, to upgrade to Chirpy Red
- Chirpy Red gives users nothing besides status

## Prerequisites

### Database
You will need PostgreSQL server running and with an empty database for Chirpy

### Environment Variables

| Variable     | Required | Description                                                     |
| :----------- | :------- | :-------------------------------------------------------------- |
| `DB_URL`     | Yes      | The connection string for the database.                         |
| `JWT_SECRET` | Yes      | A string used to sign JWT access tokens.                        |
| `POLKA_KEY`  | Yes      | The API key used the Polka payment system's webhook.            |
| `APP_PATH`   | No       | Path to your static content. Defaults to the current directory. |
| `APP_PREFIX` | No       | The prefix for requests to static content. Defaults to `/app/`. |
| `PORT`       | No       | Port for Chirpy to listen on. Defaults to `8080`.               |
| `PLATFORM`   | No       | The environment Chirpy is running in, e.g. `prod`, `dev`, etc.  |

TIP: You can generate random strings to use for the `JWT_SECRET` and `POLKA_KEY` with `openssl rand`:
```shell
# Generate a JWT_SECRET
$ openssl rand -base64 48
rYF7fX2vss+uAOIUQSZfBNltfp4QUtVjN64ITZvMFjxGp6w9cebwD7e6BfABzI7k

# Generate a POLKA_KEY
$ openssl rand -hex 16
9e2202972e5faaa93dcadaf43b56edaa
```

You can set environment variables in 2 ways:
1. Via environment variables:
   ```shell
   export DB_URL=DB_URL
   export JWT_SECRET=JWT_SECRET
   export POLKA_KEY=POLKA_KEY
   export APP_PATH=APP_PATH
   export APP_PREFIX=APP_PREFIX
   export PORT=PORT
   export PLATFORM=PLATFORM
   ```
2. Via a `.env` file in the current directory with the corresponding variables.
   ```conf
   DB_URL="postgres://dbuser:dbpassword@dbhost:5432/chirpy?sslmode=disable"
   JWT_SECRET="uE0VySfVwuBGxWMxlCDT7z6XE/apuHYgsV5HaOrdgDXzM3PMg1h1Vnb3+iIPnsiA"
   POLKA_KEY="114e0ca83e5bcfee2f7f8b93119fffc8"
   APP_PATH="/path/to/static/content"
   APP_PREFIX="/myapp/"
   PORT=8888
   PLATFORM="dev"
   ```

## Installation
```shell
go install github.com/WadeGulbrandsen/chirpy@latest
```
This will install the `chirpy` binary to your `$HOME/go/bin` directory.

## Running Chripy
Just run the binary. You can add to `$HOME/go/bin` to your path.
```shell
# If you didn't set your path:
$ ~/go/bin/chirpy

# If you have $HOME/go/bin in your path:
$ chirpy
```

## Metrics
Chirpy keeps track of the number of hits to the static content since the server has been started.
You can view the hit count at `/admin/metrics`.

## API Authentication

After a [user is created](#post-apiusers) they can [login](#post-apilogin).
If the login is successful they can get the `token` and `refresh_token` from the response.

The `token` is a JWT token and is valid for 1 hour. The `refresh_token` can be used to get
a new `token` by sending a [`POST /api/refresh`](#post-apirefresh) request. The `refresh_token` is valid for 60 days. After 60 days the user will need to
[login](#post-apilogin) again.

The `token` should be included as a Bearer token in the header for any requests that
require authentication.
```
Authorization: Bearer TOKEN
```

## API Endpoints

### `GET /api/healthz`
Returns response code `200` and the text `OK` if Chirpy is up.

### `GET /api/chirps`
Get a list of chirps. Can optionally be filtered for chirps from a specific user.

#### Query Parameters
- `author_id`: The ID of a user to filter the chirps by.
   If ommited all chirps will be returned.
- `sort`: Contols how the list of chirps is sorted.
   If `sort` is ommitied or has an invalid value it will default to `asc`
   - `asc`: oldest first
   - `desc`: newest first

#### Responses
- `200`: JSON representation of the list of [chirps](#chirp):
   ```json
   [
    CHIRP,
    CHIRP,
    ...
    CHIRP
   ]
   ```
- `400`: Invalid `author_id`
- `500`: Error retrieving chirps from the database

### `POST /api/chirps`

Creates a new chrip for the authenticated user with the provided `body`.

#### Input
JSON object with the following structure:
```json
{
    "body": STRING
}
```

#### Responses
- `201`: The chirp was successfully created. Returns the [chrip](#chirp)
- `401`: A valid access token was not provided
- `500`: Other errors such as malformed input

### `GET /api/chirps/{chirpID}`

Gets a chirp by its ID.

#### Responses
- `200`: Returns the [chrip](#chirp)
- `400`: The `chirpID` is not valid
- `404`: The chirp was not found

### `DELETE /api/chirps/{chirpID}`

Deletes the chirp with the provided `chirpID`.
Only the user that created the chirp can delete it.

#### Responses
- `204`: The chirp was successfully deleted
- `400`: The `chirpID` is not valid
- `401`: A valid access token was not provided
- `403`: The user doesn't have permission to delete the chirp
- `404`: The chirp was not found
- `500`: There was an error when deleting the chrip

### `POST /api/login`

Logs in a user with the given `email` and `password`. If successful it returns a
[UserWithToken](#userwithtoken) that contains the authentication tokens needed to
access other endpoints.

#### Input
A JSON object with the following structure:
```json
{
    "email":    STRING,
    "password": STRING
}
```

#### Responses
- `200`: Successfully logged in. Returns a [UserWithToken](#userwithtoken)
- `401`: Invalid email or password
- `500`: Error processing the request

### `POST /api/refresh`

Use the `refresh_token` to get a new access token.

The refresh token is specifed in the headers for the request:
```
Authorization: Bearer REFRESH_TOKEN
```

#### Responses
- `200`: The refresh was successful. Returns the new [Token](#token)
- `401`: The provided `refresh_token` wasn't valid

### `POST /api/revoke`

Revokes a `refresh_token` so it can no longer be used. A user will need to
[login](#post-apilogin) again to get a new `refresh_token`.

The `refresh_token` is specifed in the headers for the request:
```
Authorization: Bearer REFRESH_TOKEN
```

#### Responses
- `204`: The `refresh_token` was revoked
- `401`: The provided `refresh_token` wasn't valid

### `POST /api/users`

Create a new user with the given `email` and `password`.

#### Input
A JSON object formatted as:
```json
{
    "email":    STRING,
    "password": STRING
}
```

#### Responses
- `201`: The user was created. Returns the [User](#user)
- `400`: Invalid email or password
- `500`: Error creating the user

### `PUT /api/users`

Updates a user's `email` and `password`. Both are required.

#### Input
A JSON object formatted as:
```json
{
    "email":    STRING,
    "password": STRING
}
```

#### Responses
- `200`: The user was updated. Returns the [User](#user)
- `400`: Invalid email or password
- `401`: The request was not authenticated
- `500`: Error updating the user

### `POST /api/polka/webhooks`

This endpoint simulates a webhook for a fictional payment processor called **Polka**.

Requests to this endpoint need to include an API key in the headers for the request:
```
Authorization: ApiKey POLKA_API_KEY
```

#### Input
A JSON object with the following structure:
```json
{
    "event": STRING,
    "data": {
        "user_id": STRING
    }
}
```
- `event`: Currently supported event types (unknown event types will return `204`)
   - `"user.upgraded"`: Upgrades the user to **Chirpy Red**
- `user_id`: The user's UUID in string form

#### Responses
- `204`: The event was processed successfully
- `401`: The API Key is invalid
- `404`: The user couldn't be found
- `500`: There was an error processing the event

### `POST /admin/reset`
Deletes all users and chirps from the database and resets the hit count to `0`.
This only works if `PLATFORM` is set to `dev`.

#### Responses
- `200`: The database and hit counts have been reset
- `403`: `PLATFORM` is not set to `dev`
- `500`: There was an error when trying to clear the database

## API Objects

### Chirp
```json
{
    "id":         UUID,
    "created_at": TIMESTAMP,
    "updated_at": TIMESTAMP,
    "body":       STRING,
    "user_id":    UUID
}
```
- `id`: Unique identifier for the chirp
- `created_at`: Time the chirp was created
- `updated_at`: Time the chirp was last updated
- `body`: The text of the chirp
- `user_id`: Unique identifier for the user that created the chirp

### Token
```json
{
    "token": STRING
}
```
- `token`: A JWT access token

### User
```json
{
    "id":            UUID,
    "created_at":    TIMESTAMP,
    "updated_at":    TIMESTAMP,
    "email":         STRING,
    "is_chirpy_red": BOOLEAN
}
```
- `id`: Unique identifier for the user
- `created_at`: Time the user was created
- `updated_at`: Time the user was last updated
- `email`: Email address of the user
- `is_chirpy_red`: Indicates if the user has **Chirpy Red**

### UserWithToken
```json
{
    "id":            UUID,
    "created_at":    TIMESTAMP,
    "updated_at":    TIMESTAMP,
    "email":         STRING,
    "is_chirpy_red": BOOLEAN,
    "token":         STRING,
    "refresh_token": STRING
}
```
- `id`: Unique identifier for the user
- `created_at`: Time the user was created
- `updated_at`: Time the user was last updated
- `email`: Email address of the user
- `is_chirpy_red`: Indicates if the user has **Chirpy Red**
- `token`: The user's access token (valid for 1 hour)
- `refresh_token`: The user's refresh token (valid for 60 days)
