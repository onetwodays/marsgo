curl -X POST --location "http://127.0.0.1:38888/api/v1/textsecret/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\":\"otcexchange\",
        \"sign\":\"SIG_K1_Jz4KTq5v3dhcRKYehza6NRF5SxZaEzdPiVohRhX5SqoDCjmjf6hh3vyqfFHzEUagWZHC3L6G6SaJvKdH3UWEVwJRLWc6jL\"

        }"


-- eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjEzMzAwMDMsImlhdCI6MTYyMTI0MzYwMywibmFtZSI6Im90Y2V4Y2hhbmdlIn0.6N5Ni9o7Qjht1FqY1urd7Su7HDpZrEBfpk-jn5cdT8M



curl -X GET --location "http://192.168.0.177:38888/api/v1/textsecret/keys/otcexchange/1" \
    -H "Content-Type: application/json" \
    -H "x-user-name: otcexchange" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjEzMzAwMDMsImlhdCI6MTYyMTI0MzYwMywibmFtZSI6Im90Y2V4Y2hhbmdlIn0.6N5Ni9o7Qjht1FqY1urd7Su7HDpZrEBfpk-jn5cdT8M"




curl -X PUT --location "http://127.0.0.1:38888/api/v1/textsecret/keys" \
    -H "Content-Type: application/json" \
    -H "x-user-name: otcexchange" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjEzMjQ1MDAsImlhdCI6MTYyMTIzODEwMCwibmFtZSI6Im90Y2V4Y2hhbmdlIn0.wGlrSPtsAkChPSMLF6j1hi3VmzWM49GMXYztzTtB_2U" \
    -d "{
          \"identityKey\":\"identityKey\",
          \"signedPreKey\":{
              \"signature\": \"signature\",
              \"prekey\": {
                \"keyId\": 1,
                \"publickey\": \"54\"
              }
            },

          \"prekeys\":[
            {
              \"keyId\": 0,
              \"publickey\": \"pk1\"
            },
            {
              \"keyId\": 1,
              \"publickey\": \"pk2\"
            }
          ]

        }"



curl -X PUT --location "http://127.0.0.1:38888/api/v1/textsecret/messages/otcexchange" \
    -H "Content-Type: application/json" \
    -H "x-user-name: otcexchange" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjI4ODEyMTMsImlhdCI6MTYyMjc5NDgxMywibmFtZSI6Im90Y2V4Y2hhbmdlIn0.1MykydYaLnPc56O5y9Nf8Q8XOwwkQCpUNTQ3khA1rdk" \
    -d "{
        \"destination\":\"otcexchange\",
        \"online\":false,
        \"timestamp\": 1,
        \"messages\": [
           {
            \"type\": 1,
            \"destination\": \"otcexchange\",
            \"destinationDeviceId\": 1,
            \"destinationRegistrationId\": 1,
            \"body\": \"11111\",
            \"content\": \"22222\",
            \"relay\": \"wwww.baidu.com\"
          }

        ]

        }"



curl -X GET --location "http://127.0.0.1:38888/api/v1/textsecret/messages" \
    -H "Content-Type: application/json" \
    -H "x-user-name: otcexchange" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjI4ODEyMTMsImlhdCI6MTYyMjc5NDgxMywibmFtZSI6Im90Y2V4Y2hhbmdlIn0.1MykydYaLnPc56O5y9Nf8Q8XOwwkQCpUNTQ3khA1rdk"