# Fuji Lane - Backend

This app is responsible for responding to calls from Fuji Lane frontend.

## Facebook Authentication

To authenticate Facebook tokens we need to use an access token that is generated with the following call:

```
curl -X GET "https://graph.facebook.com/oauth/access_token?client_id=FACEBOOK_APP_ID&client_secret=FACEBOOK_APP_SECRET&grant_type=client_credentials"
```

The result is stored in `FACEBOOK_CLIENT_TOKEN`. More details can be found
[here](https://developers.facebook.com/docs/facebook-login/access-tokens/#debug).

`FACEBOOK_APP_ID` and `FACEBOOK_APP_SECRET` are available in
[this page](https://developers.facebook.com/apps/299325330622870/settings/basic/).
