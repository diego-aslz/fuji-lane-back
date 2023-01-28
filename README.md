# Fuji Lane - Backend
[![CircleCI](https://circleci.com/gh/nerde/fuji-lane-back.svg?style=svg&circle-token=7e16acfd129328a136688c39818f3be9f4204693)](https://circleci.com/gh/nerde/fuji-lane-back)

This app is responsible for responding to calls from Fuji Lane frontend.

## Code Organization

This app's code is organized in the following fashion:

- `flactions` contains the actions that can be invoked by API clients. All HTTP routes handled by this app must match
  an action.
- `flconfig` contains configuration files and configuration code.
- `fldiagnostics` contains code useful for tracking system operation so it can be logged and used for debugging.
- `flentities` is the lowest level package of the app and contains the persistent objects and basic database access code.
- `flservices` contains abstractions for third-parties used by this project, like S3 and Facebook.
- `flweb` contains code related to booting up the web server, middleware and routing requests to actions.
- `migrations` contains database migration files.
- `test` contains test code.

## Configuration

For test environment, the syste will read configuration from `flconfig/test.yml`. To run the app in development mode,
you'll probably need a `flconfig/development.yml` file, which is not version controlled for security reasons. For
basic usage, just `cp flconfig/test.yml flconfig/development.yml` and **fixing the database name** should be enough.
If you need credentials you'll have to ask a co-worker for them.

## Usage

- Install dependencies: `make dependencies`.
- Migrating the database: `make migrate`. For test database, use `STAGE=test make migrate`.
- Run all the tests: `make feature`.
- Run the web server: `APP_ROOT=. REDIS_URL=http://localhost:6379 make server`.

## Facebook Authentication

To authenticate Facebook tokens we need to use an access token that is generated with the following call:

```
curl -X GET "https://graph.facebook.com/oauth/access_token?client_id=FACEBOOK_APP_ID&client_secret=FACEBOOK_APP_SECRET&grant_type=client_credentials"
```

The result is stored in `FACEBOOK_CLIENT_TOKEN`. More details can be found
[here](https://developers.facebook.com/docs/facebook-login/access-tokens/#debug).

`FACEBOOK_APP_ID` and `FACEBOOK_APP_SECRET` are available in
[this page](https://developers.facebook.com/apps/299325330622870/settings/basic/).
