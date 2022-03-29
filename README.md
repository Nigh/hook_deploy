# hook_deploy

auto deploy from WebHook

`WebHook -> Pull -> Build -> Run`

## env

`Golang 1.18`

## config

Edit `config.json`

```json
{
	"port": 13737,
	"apps":[
		{
			"type": "gitee",
			"git_name": "YUI",
			"git_url" : "",
			"branch": "master",
			"project_dir" : "~/apps/YUI",
			"build": "go build -o ./AppBin",
			"deploy": "mv ./AppBin ./deploy/application"
		}
	]
}
```
