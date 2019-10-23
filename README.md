# URL Watcher

Inspired by [Kubernetes's Efficient Detection of Changes](https://kubernetes.io/docs/reference/using-api/api-concepts/#efficient-detection-of-changes)

Assuming you have a list of URLs:
- Regularly fetch the URLs
- Bump the resourceVersion if it changes
- Store each version and a diff between each one
- Occasionally GC old versions, assuming you only support 24 hours of history or 10 versions (like revisionHistoryLimit)

## Example

Assume we have an API that returns a JSON array like:

```
[
{
  "name": "burger"
},
{
  "name": "pizza"
}
]
```

If someone fetches it via our API,  we can send them the latest version we have, and a resource version.

e.g.

```
{
  "version": 3,
  "items": [
    {
      "name": "burger"
    },
    {
      "name": "pizza"
    }
  ]
}
```

Then if someone adds a new food, like "ramen", an existing client would get the delta since we'd know they have version 3 already, but a new client would get the full list.

Existing client:

```
{
  "delta": true,
  "version": 4,
  "items": [
    {
      "name": "ramen"
    }
  ]
}
```

New client:

```
{
  "delta": true,
  "version": 4,
  "items": [
    {
      "name": "burger"
    },
    {
      "name": "pizza"
    },
    {
      "name": "ramen"
    }
  ]
}
```
