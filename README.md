# Setup

## Generate keys

On MacOS:

```bash
ssh-keygen -t rsa -b 4096 -f id_rsa -m PEM
openssl rsa -in id_rsa -pubout -outform pem > id_rsa.pem
```

## Setup Config

Create config/default.json, should contain the fields listed in the `fox.go/Config` struct
Paths to the keys are relative to the root of the project

# Running

From scratch:

```bash
docker-compose up -d
make migrate
make run
```
