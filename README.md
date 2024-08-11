# Go-Lang-Publishing-Platform

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