# cf-crontab
Run a crontab like scheduler as a (16MB) microservice on Cloud foundry. It knows how to trigger `http`, `amqp`, `cartel` and `iron` events.
Uses bound services for credentials retrieval. 

The binary also doubles as a CF CLI plugin and exposes commands to manage your crontab (add, delete, backup)

Finally, it uses the server components environment to persist your crontab so you don't even need a backing store, we are going for cheap! Your crontab entries will thus survive restarts and Cloud foundry updates.

# requirements
1. A clever password to protect the crontab:

```shell script
pwgen -s 32
```
pick one and replace where you see `PASSWORD_HERE`

2. A unique router name, so the CF plugin can talk to it.
```shell script
pwgen -A 20
```

Pick one and replace where you see `HOST_HERE`

# deployment
Use the below template manifest and fill in above values, save as `manifest.yml`
```yaml
applications:
- name: cf-crontab
  disk_quota: 128M
  docker:
    image: loafoe/cf-crontab:0.0.2
  env:
    CF_CRONTAB_SECRET: PASSWORD_HERE
  instances: 1
  memory: 16M
  routes:
  - route: HOST_HERE.cloud.pcftest.com
```

```shell script
cf push -f manifest.yml
```
# CF CLI plugin
TODO

# Contact / Getting help

- andy.lo-a-foe@philips.com

# License
License is MIT
