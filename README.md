#json-charlieegan3

A rake task to keep the live content on 
[charlieegan3.com](http://charlieegan3.com) up-to-date.

Fetches data from:

* Instagram
* Twitter
* last.fm
* GitHub
* Strava

The task is currently hosted on [hyper.sh](https://hyper.sh) and updates this 
[status file](https://s3.amazonaws.com/charlieegan3/status.json).

Deployment:

```bash
hyper run --size=s2 -d --restart=on-failure --env-file ~/Desktop/json-charlieegan3.env charlieegan3/json-charlieegan3:master
```
