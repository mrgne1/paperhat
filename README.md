# Paperhat
Paperhat is an API, database and optional website for sharing secrets with others.

A person can send a secret value to Paperhat.
Paperhat will encrypt that secret and store it in a database.
Once the secret is stored succesfully Paperhat will return a link to recover the secret value.

The link contains a unique id for the stored secret and the key used to encrypt the secret.
Paperhat doesn't store the encryption key.
This means that only a person with the correct link can recover a secret.

Once recovered, the secret will be deleted from the database.
The secret will be deleted even if an error occurs during recovery.

Secrets are always created with a duration, the period of time for which the secret will available for recovery.
If the secret is not recovered before that period of time has elapsed, the secret will be deleted.

# Requirements
Paperhat was developed with Golang 1.23.1 on Ubuntu
Paperhat uses Sqlite3 for its database and requires the ability to create a new database file where it is deployed

# Running Paperhat
1. Install Golang 1.23.1
2. Clone Paperhat repo
3. Run the command: `go build && go run .`
  - Optionally configure with a `.env` file (see Configuration section)

Website will be found at `{hostname}:{port}/v1/`
Example: `localhost:2060/v1/`

API will be found at `{hostname}:port/api/`
Example: `localhost:2060/api/`
API heartbeat example: `localhost:2060/api/heartbeat`

# Configuration
Paperhat uses a `.env` file for configurations.
Example Configuration:
```
# Port for the Website and the API
PORT=2060

# Operation mode (standalone or api)
# api will deploy on the api and database for saving secrets
#   it is intended to work with a custom website
# stand alone will deploy both the api and a website for creating secrets 
MODE=standalone

# Website path
# Path to the website files. Defaults to ./site/v1 if not provided
# SITEPATH=./site/v1

# Database path
# Path to the Sqlite3 database file. Defaults to ./secrets.db
# DBPATH=./secrets.db
```

# API
`GET /api/heartbeat`
Returns a 200 message with 'OK' in the body


`POST api/secrets?duration={time in seconds}`
Creates a new secret.
The entire body of the message is the secret.
If duration is omitted then the secret will last 60 seconds.

Returns a JSON object of the form: 
```json
{
    "id": "fc4051f1-6325-494d-912b-c0c963cc6dce",
    "key": "93135b1e7bf1c811c162"
}
```

`GET /api/secrets/{id}/{key}`
Returns the plaintext secret in the body of the request
