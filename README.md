# Go-Lang-Publishing-Platform

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