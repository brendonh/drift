# 998: gqFY0QPmoVnRA+c=
# 999: gqFY0QPnoVnRA+c=
# 0 0: 

DATA=$(cat <<EOF
{
    "inputs": {
        "bucket": "PoweredBody",
        "index": "Coords_bin",
        "key": "gqFY0QPmoVnRA+c=",
        "keep": "true"
    },
    "query": [
        {
            "map": {
                "language": "erlang",
                "module": "riak_kv_mapreduce",
                "function": "map_object_value"
            }
        }
    ]
}
EOF
)
echo -n $DATA | curl -v http://localhost:8098/mapred -H 'Content-Type: application/json' --data-binary @-
echo