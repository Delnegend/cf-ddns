# CF-DDNS

## Installation
- Clone the repository
- Copy `docker-compose.example.yml` to `docker-compose.yml`
- Edit `docker-compose.yml` and set the environment variables
- Run `docker-compose up -d --build`

## Update
```bash
docker down && git stash && git pull --rebase && git stash apply && docker up -d --build
```