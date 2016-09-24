




The following error codes exist among the methods of the ARGO Web API:

Error                    | HTTP Code  | Description   
------------------------ | ---------- | ------ 
Bad request              | 400 | One or more checks may have failed. More details on the carried out checks can be found [here](validations.md)
Wrong start_time         | 400 | Use start_time url parameter in zulu format (like `2006-01-02T15:04:05Z`) to indicate the query start time
Wrong end_time           | 400 | Use end_time url parameter in zulu format (like `2006-01-02T15:04:05Z`) to indicate the query end time
Wrong exec_time          | 400 | Use exec_time url parameter in zulu format (like `2006-01-02T15:04:05Z`) to indicate the exact probe execution time
Wrong granularity        | 400 | The parameter value can be either `daily` or `monthly`
Unauthorized             | 401 | The client needs to provide a correct authentication token using the header `x-api-key` 
Forbidden                | 403 | Access to the resource is forbidden due to authorization policy enforced
Item not found           | 404 | Either the path is not found or no results are available for the given query
Content not acceptable   | 406 | The `Accept` header either was not provided or was provided but did not contain any valid content types. Acceptable content types are `application/xml` or `application/json`




