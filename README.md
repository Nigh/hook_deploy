# hook_deploy

auto deploy from WebHook

`WebHook -> Pull -> Build -> Run`

## config

Edit `config.json`

```json
{
  "port": 13737,
  "git": "https://github.com/Nigh/hook_deploy",
  "remote": "origin",
  "branch": "master",
  "build": "go build",
  "deploy": "./application"
}
```
