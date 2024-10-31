# Paperhat
Simple way to share secrets with another person.

User can create a secret.
The secret is encrypted and stored in a database.
User's are given an id for the secret and the key used to encrypt the secret.
Since only the user has the encryption key, only the user can decrypt the secret.
Secret can be recovered by passing the id and the key to the api.

Once accessed, the secret will be deleted.
Secrets will be deleted even if an error occurs and the secret is not returned.
The secret will also be deleted if the secret expires.

It's possible to run Paperhat as a stand-alone application with API, database and website or as a backend only instance with only the API and database.

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
