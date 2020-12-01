# discord-notification-resource

A Concourse resource for webhook notifications in Discord.  Written in Go.

## Example

```yaml
resource_types:
  - name: discord-notification-resource
    type: registry-image
    source:
      repository: logsquaredn/discord-notification-resource
      tag: latest

resources:
  - name: notify
    type: discord-notification-resource
    source:
      ...

jobs:
  - name: some-job
    plan:
      ...
      - put: notify
        params:
          ...
```

## Source Configuration

| Parameter    | Required | Description                                                  |
| ------------ | -------- | ------------------------------------------------------------ |
| `webhook_id` | yes      | the id of the webhook to post to _see below (1)_                             |
| `token`      | yes      | the token to use to authenticate when posting to the webhook _see below (2)_ |

> _(1)_ `webhook_id` _will be the second to last path parameter in the url copied from the_ `Copy Webhook URL` _button below_

> _(2)_ `token` _will be the last path parameter in the url copied from the_ `Copy Webhook URL` _button below_

![webhook-id](https://user-images.githubusercontent.com/60495614/100556635-a8d29b80-3271-11eb-8b46-798d5ccc8e4e.png)

## Behavior

### `check`

not implemented

### `in`

not implemented

### `out`

see [discordgo.WebhookParams](https://godoc.org/github.com/bwmarrin/discordgo#WebhookParams). Addtionally:

| Parameter | Required | Description                                                                         |
| --------- | -------- | ----------------------------------------------------------------------------------- |
| `wait`    | no       | whether or not to wait on the response and gather version and metadata info from it |
