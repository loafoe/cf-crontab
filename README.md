# cf-crontab
Run a crontab like scheduler as a (16MB) microservice on Cloud foundry. It knows how to trigger `http`, `amqp`, `cartel` and `iron` events.
Uses bound services for credentials retrieval. 

The binary also doubles as a CF CLI plugin and exposes commands to manage your crontab (add, remove, save)

Finally, it uses the server components environment to persist your crontab so you don't even need a backing store, we are going for cheap! Your crontab entries will thus survive restarts and Cloud foundry updates.

# Requirements
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

# Deployment
Use the below template manifest and fill in above values, save as `manifest.yml`

```yaml
applications:
- name: cf-crontab
  disk_quota: 128M
  docker:
    image: loafoe/cf-crontab:0.0.8
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
# CF plugin
For now you should clone this repo and build the binary. Binary distribution to follow:

```shell
go build .
cf install-plugin -f cf-crontab
```

This will enable a few new commands:

```
~ > cf plugins
Listing installed plugins...

plugin              version   command name    command help
cf-crontab          0.0.1     add-cron        Add a cron job
cf-crontab          0.0.1     crontab         List all crontab entries
cf-crontab          0.0.1     remove-cron     Remove a cron job
cf-crontab          0.0.1     save-crontab    Save crontab table to the environment
```

# Tasks
You define tasks using the below JSON. A more detailed description and of all supported and planned tasks will follow soon. For now we have `http`:

```json
[
  {
    "schedule": "* */10 * * * *",
    "job": {
      "type": "http",
      "command": {
        "method": "GET",
        "url": "https://icanhazip.com/"
      }
    }
  }
]
```

# Scheduling tasks
Save the above json as `tasks.json` and then:

```
cf add-cron tasks.json
```

```
Adding 1 entries ...
┌───┬────────────────┬──────┬────────────────────────────┐
│ # │ SCHEDULE       │ TYPE │ DETAILS                    │
├───┼────────────────┼──────┼────────────────────────────┤
│ 1 │ 0 */10 * * * * │ http │ GET https://icanhazip.com/ │
└───┴────────────────┴──────┴────────────────────────────┘
OK
```

Now every 10 minutes we will call icanhazip for no good reason.

# Listing tasks

```
cf crontab
```

# Remove a cron task
You should use the ID in the first column of the crontab list to select the entry you wish to remove.
```
cf remove-cron 1
```

# Saving your crontab
Whenever you have made changes to your crontab it's important to save a copy so the server can re-seed its scheduler in the event of a restart or Cloud foundry upgrade.

```
cf save-crontab
```

After a few seconds the active crontab will be saved in the applications environment. You can save the manifest locally to make an off-site backup. Yeah, no database required so we keep the costs to a minimum.

# Contact / Getting help

- andy.lo-a-foe@philips.com

# License
License is MIT
