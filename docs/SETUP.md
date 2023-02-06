## Requirements

There are several components required for the application to work. As a disclaimer, I developed this on Debian Bullseye so that's the only OS I can guarantee it will work on but it will probably work on your machine as I tried to not include anything too platform-specific


### Front-End

The front end is AlpineJS and HTMX with minimal additional Javascript. Download the minified JS bundles and place them in the public assets folder with the commands below (These are the latest versions as of writing this documentation)

#### **HTMX**
`wget https://unpkg.com/htmx.org@1.8.4 && mv htmx-org@1.8.4 src/front-end/public/bin/html-1.8.4-min.js`

#### **AlpineJS**
`wget https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js && mv cdn.min.js src/front-end/public/bin/alpine-3.10.5-min.js`

___


### Back-End

#### **Podman**

For running the containers. See https://podman.io/getting-started/installation for instructions, but your distro probably has a `podman` package already so just install that if it's available

If you have Docker already that will also work, just edit the start scripts to fix the image name

#### **PostgreSQL**

This is only needed if you want to connect to the Postgres container from your host. If you're fine with running `psql` from inside `pg`, you can skip this part.

Install instructions are at https://www.postgresql.org/download/. Note that at least on my system this starts a daemon that inteferes with the container we're going to start in the steps below, `systemctl stop postgresql` does the trick

#### **Go**

Version 1.18 or higher is required. Download from https://go.dev/dl/


## Running Locally

1. Add `127.0.0.1       marketplace.test` to your hosts file (probably `/etc/hosts`)

2. Start the database: `./src/db/start.sh`

3. Start the server: `cd src/server && go run .`. Credentials are by default read from `config/config.json` and should match the ones used in step 3

4. Start Nginx: `./src/nginx/start.sh`

5. If all went well, the application is now accessible at `http://marketplace.test`. Your browser will say it's insecure because I haven't set up HTTPS yet. Ignore the warning, it's fine

