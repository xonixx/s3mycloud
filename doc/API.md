# API (MVP 0.0.1)

## 1. Upload
   
### I. Upload metadata
```
POST /api/file/upload

{
    "name": "file_name.ext",           # required
    "size": 123,                       # required, size in bytes
    "tags": ["tag1", "tag2", "tag3"]   # optional
}
```
response (success):
```
{
    "id": "ID"                         # unique ID as assigned by storage engine
    "uploadUrl"                        # S3 URL for PUT
}
```
response (error):
```
{
    "error": "error description"
}
```

### II. Upload file
```
PUT https://s3-url

file bytes...
```

## 2. List files
```
GET /api/file?name=partOfName&tags=tag1,tag2,tag3&page=2&pageSize=10&sort=name,desc
```
                                           
| param    | optional | default       | description                           |
|----------|----------|---------------|---------------------------------------|
| name     | yes      |               | part of name to filter by             |
| tags     | yes      |               | list of tags to filter by (AND-logic) |
| page     | yes      | 0             | 0-based                               |
| pageSize | yes      | 10            |                                       |
| sort     | yes      | uploaded,desc |                                       |

Sort fields:
- name
- size
- uploaded

response (success):
```
{
    "page": [
        { 
            "id", 
            "name", 
            "size", 
            "tags", 
            "url", 
            "uploaded"    # timestamp or YYYY-MM-DD hh:ss TODO 
        },
        ...
    ],
    total: 27             # total records
}
```

## 3. Delete file
```
DELETE /api/file/{ID}
```
response (success):
```
{
    "success": true
}
```
response (error):
```
{
    "error": "file not found"
}
```

## 4. Assign tags
```
POST /api/file/{ID}/tags

["tag1", "tag2", "tag3"]
```
returns status 200 and body
```
{"success": true}
```

## 5. Remove tags

```
DELETE /api/file/{ID}/tags

["tag1", "tag3"]
```

returns status 200 if all OK and body
```
{"success": true}
```

returns status 400 if one of the tags is not present on the file and body
```
{
    "success": false,
    "error": "tag not found on file"
}
```
