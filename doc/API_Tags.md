# Tags API

## 1. List tags
```
GET /api/tag?q=prefixOfName&limit=10
```

```json
[
  {
    "id": "1",
    "name": "tag1",
    "color": "#ffeedd"
  }
]
```

## 2. Add tag

```
POST /api/tag

{
    "name": "tag1",
    "color": "#ffeedd"
}
```

response (success):
```json5
{
    "id": "ID",                        // unique ID as assigned by storage engine
}
```
response (error):
```json
{
    "error": "error description"
}
```