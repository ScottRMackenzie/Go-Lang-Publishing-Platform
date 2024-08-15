# Go-Lang-Publishing-Platform

## Distributed architecture differences
The api and frontend will run on different ports.
Impacts of this:
- for any protected route you need to use `credentials: "include"` so the cookies are send with the request


## API

### Valid JSON Requests

#### Search

`{
    "search_query": "19",
    "sort_by": {
        "field": "title",
        "order": "ASC"
    },
    "pagination": {
        "start_index": 0,
        "max_results": -1
    },
    "filters": {
        "case_sensitive": {
            "author": true
        },
        "values": {
            "author": "%Orwe%"
        }
    }
}`