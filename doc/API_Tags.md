# Tags API

## 1. List tags
```
GET /api/tag?q=prefixOfName&page=2&pageSize=10
```
Returns sorted alphabetically

```json
{
  "page": [
    {
      "id": "1",
      "name": "tag1",
      "color": "#ffeedd"
    }
  ],
  "total": 23
}
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

## 3. Edit tag
(Must pass all values)
```
PUT /api/tag

{
    "id": "ID",
    "name": "newName",
    "color": "#aabbcc"
}
```
