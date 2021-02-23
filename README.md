# Agit PR alaram docker action

This action that notifies agit when pr is open.

## Inputs

### `url`

**Required** Agit webhook URL. Default `""`.

### `event`

**Required** Github PR Event api url. Default `""`.

## Outputs

### `response`

Response when sending agit webhook.

## Example usage
```
uses: comento/pr-alarm-to-agit-action@master
with:
    url: ${{ secrets.agit_webhook }}
    event: ${{ github.event.pull_request.url }}
```