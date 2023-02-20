# requirements
    * Docker
    * Docker Compose
    * Go >= 1.19
# setup

`git clone https://github.com/itzurabhi/companies-micro.git && cd companies-micro`

create a `.env` file with necessary changes from `.dev.env` or use the same file for `--env-file` option.

`docker-compose -f docker-compose.dev.yaml --env-file .env up`
