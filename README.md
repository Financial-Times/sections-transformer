# sections-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/sections-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/sections-transformer/tree/master)

Retrieves Sections taxonomy from TME and transforms the sections to the internal UP json model.
The service exposes endpoints for getting all the sections and for getting section by uuid.

# Usage
`go get github.com/Financial-Times/sections-transformer`

```
Options:
  --tme-username=""                                           TME username used for http basic authentication ($TME_USERNAME)
  --tme-password=""                                           TME password used for http basic authentication ($TME_PASSWORD)
  --token=""                                                  Token to be used for accessing TME ($TOKEN)
  --base-url="http://localhost:8080/transformers/sections/"   Base url ($BASE_URL)
  --tme-base-url="https://tme.ft.com"                         TME base url ($TME_BASE_URL)
  --port=8080                                                 Port to listen on ($PORT)
  --maxRecords=10000                                          Maximum records to be queried to TME ($MAX_RECORDS)
  --slices=10                                                 Number of requests to be executed in parallel to TME ($SLICES)
```

With Docker:

`docker build -t coco/sections-transformer .`

`docker run -ti --env BASE_URL=<base url> --env TME_USERNAME=<foo> --env TME_PASSWORD=<bar> --env TOKEN=<foobar> --env PORT=8080 --env MAX_RECORDS=1000 --env SLICES=10 --env= TME_BASE_URL=foobar coco/sections-transformer`
