# OpenNet Worker

Install Wrangler:

npm install -g wrangler

Login:

wrangler login

Create a token:

openssl rand -hex 32

Set the token:

wrangler secret put OPENNET_TOKEN

Deploy:

wrangler deploy
