# giveawaybot

It's a great thing to win but it's an even greater thing to help others win just as you already have.

# Usage

Because there's not much time to work on a frontend interface, this bot is currently only accessible via API endpoints via the browser. The endpoints are:

- `GET / - Returns a 200 OK response if the bot is running`

- `GET /v1/health - Returns a 200 OK response if the bot is running. This can include additional information about the bot's health`

- `GET /v1/auth/twitter/begin - Returns a URL to redirect the user to in order to authenticate with Twitter`

- `GET /v1/auth/twitter/signon - The callback URL for Twitter authentication. This will register the user if not already registered or log them in if they are`

- `GET /v1/giveaway/simulate/start - Starts a new giveaway if one is not already running. This will also send a tweet to Twitter on behalf of the signed in user`

- `GET /v1/giveaway/simulate/disrupt - Disrupts the current giveaway if one is running. This will also delete the tweet on Twitter on behalf of the signed in user`

- `GET /v1/giveaway/simulate/end - Ends the current giveaway if one is running. It sets the completed status and the completed_at timestamp.`

- `GET /v1/giveaway/tweet/replies - Makes a recursive call to Twitter to get all replies to the giveaway tweet. This is used to determine the winner of the giveaway.It also sets the giveaway status to inactive.`

- `GET /v1/giveaway/tweet/report - Compiles a report of the giveaway and returns it as a JSON object`

# Running

First, clone the repo and install the dependencies:

```bash
git clone https://github.com/opensaucerers/giveawaybot.git
cd giveawaybot
go mod tidy
```

Next, create a .env file and copy the contents of .env.example into it. Then, fill in the values for the environment variables.

```bash
touch .env
cp .env.example .env
```

Now, start the app with the make command:

```bash
make run
```
