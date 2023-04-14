# giveawaybot

It's a great thing to win but it's an even greater thing to help others win just as you already have.

# Usage

Because there's not much time to work on a frontend interface, this bot is currently only accessible via API endpoints via the browser.

### Authentication

After completing the twitter sign in step, a JSON containing your details and JWT token will be returned. This token is only valid for 2 hours, however, for every subsequent authenticated calls, a new token is required which is also valid for another 2 hours. This way, it becomes possible to stay logged in for as long as you keep using the app.

## Endpoints

- `GET / - Returns a 200 OK response if the bot is running`

- `GET /v1/health - Returns a 200 OK response if the bot is running. This can include additional information about the bot's health`

- `GET /v1/auth/twitter/begin - Returns a URL to redirect the user to in order to authenticate with Twitter`

- `GET /v1/auth/twitter/signon - This endpoint is to be called directly by Twitter (for now). This will register the user if not already registered or log them in if they are`

- `GET /v1/giveaway/simulate/start?token= - Starts a new giveaway if one is not already running. This will also send a tweet to Twitter on behalf of the signed in user`

- `GET /v1/giveaway/simulate/disrupt?token= - Disrupts the current giveaway if one is running. This will also delete the tweet on Twitter on behalf of the signed in user`

- `GET /v1/giveaway/simulate/end?token= - Ends the current giveaway if one is running. It sets the completed status and the completed_at timestamp.`

- `GET /v1/giveaway/tweet/replies?token= - Makes a recursive call to Twitter to get all replies to the giveaway tweet. This is used to determine the winner of the giveaway.It also sets the giveaway status to inactive.`

- `GET /v1/giveaway/tweet/report/{id} - Compiles a report of the giveaway and returns it as a JSON object`

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
