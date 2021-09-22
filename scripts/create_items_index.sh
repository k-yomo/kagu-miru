curl -XPUT  -H "Content-Type: application/json" -d @./elasticsearch/mappings/items.json 'localhost:9200/items/'
curl -XPUT  -H "Content-Type: application/json" -d @./elasticsearch/mappings/items.query_suggestions.json 'localhost:9200/items.query_suggestions/'
