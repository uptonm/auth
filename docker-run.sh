docker run --network="host" \
 -e HOST=127.0.0.1 \
 -e PORT=8080 \
 -e ENV=prod \
 -e AUTH0_CLIENT_ID=nI9NUFy5rPjlMQYuTgMdLPhedFzRxRBo \
 -e AUTH0_CLIENT_SECRET=ngx58r78Hd9Z-hjNKAlC1URsjK9FQDVcQEL6TmMhIZdopY5umRGf2ygEGI4dMMWB \
 -e AUTH0_CALLBACK_URL=http://localhost:8080/api/v1/auth/auth0/callback \
 -e AUTH0_DOMAIN=uptonm.us.auth0.com uptonm/uptonm.io