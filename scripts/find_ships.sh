DATA=$(cat <<EOF
{
    "inputs": {
        "bucket": "Ship",
        "index": "Owner_bin",
        "key": "brendonh"
    },
    "query": [
        {
            "map": {
                "language": "javascript",
                "source": "function(o) { return [[o.key, o.values[0].data]]; }"
            }
        }
    ]
}
EOF
)
echo -n $DATA | curl -v http://localhost:8098/mapred -H 'Content-Type: application/json' --data-binary @-
echo