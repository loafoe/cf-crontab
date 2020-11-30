# cf-standalone
Run a crontab like scheduler as a (16MB) microservice on Cloud foundry. It knows how to trigger `http`, `amqp`, `cartel` and `iron` events.
Uses bound services for credentials retrieval. 

# Requirements

Pick one and replace where you see `HOST_HERE`

# Deployment
Use the below template manifest and fill in above values, save as `manifest.yml`

```yaml
applications:
- name: cf-crontab
  disk_quota: 128M
  docker:
    image: loafoe/cf-standalone:0.0.8
  env:
    CF_CRONTAB_SECRET: PASSWORD_HERE
  instances: 1
  memory: 16M
  routes:
  - route: HOST_HERE.cloud.pcftest.com
```

Save the above template with your values and then deploy:

```shell script
cf push -f manifest.yml
```
# Config
For now you can only inject the config via the ENV:


# Tasks
You define tasks using the below JSON. A more detailed description and of all supported and planned tasks will follow soon. For now we have `http`:

```json
[
    {
      "schedule": "*/5 * * * * *",
      "job": {
        "type": "http",
        "command": {
          "headers": {
              "Authorization": "Basic cG9sbGVyLXdlbGxjZW50aXZlOnRlc3QtcG9sbGVyCg=="
          },  
          "body": "{ \"countToProcess\": 0, \"serverName\": \"string\"}",
          "method": "POST",
          "url": "https://dm-sftp-poller.eu-west.philips-healthsuite.com/poll"
        }
      },
      "entryID": 1
    }
]
```

# Contact / Getting help

- andy.lo-a-foe@philips.com

# License
License is MIT
