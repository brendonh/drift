DATA=$(cat <<EOF
{
    "inputs": {
        "bucket": "ShipLocation",
        "index": "Coords_bin",
        "key": "gqFYAKFZAA=="
    },
    "query": [
        {
            "map": {
                "language": "javascript",
                "source": "function(o) { return [o.values[0].data]; }"
            }
        }
    ]
}
EOF
)
echo -n $DATA | curl -v http://localhost:8098/mapred -H 'Content-Type: application/json' --data-binary @-
echo