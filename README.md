# Go-Lang-Publishing-Platform

## Distributed architecture differences
The api and frontend will run on different ports.
Impacts of this:
- for any protected route you need to use `credentials: "include"` so the cookies are send with the request


## API

### Valid JSON Requests

#### Search

`{
    "search_query": "1984",
    "sort_by": {
        "field": "title",
        "order": "ASC"
    },
    "pagination": {
        "start_index": 0,
        "max_results": -1
    },
    "filters": {
        "exact_match": {
            "case_sensitive": {
                "genre": true
            },
            "values": {
                "genre": "Dystopian",
                "language_code": "en"
            }
        },
        "partial_match": {
            "case_sensitive": {
                "author": false
            },
            "values": {
                "author": "George"
            }
        }
    }
}`