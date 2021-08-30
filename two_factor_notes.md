# Two-Factor Authentication Notes

## Task List

- [ ] Allow user to enable/disable 2FA
- [ ] Allow user to choose SMS or Authy as primary 2FA method
- [ ] Verify user phone number for SMS
- [ ] Register Authy account for new users
- [x] Choose 2FA method page (post-login)
- [x] Generate and send OTP via SMS
- [x] Verify/Reject SMS code
- [ ] Send Authy OneTouch
- [ ] Accept/Reject OneTouch response
- [ ] Generate Backup Codes
- [ ] Verify/Reject Backup Code
- [x] Restrict access of users awaiting second factor

## Login Workflows

1. Show user sign-in

2. Validate email and password
    * Return to 1 if sign-in fails

3. If account is not 2FA
    * Set user session cookie with full access
    * Redirect to dashboard
    * End

    - Else -

    * Set user session cookie with restricted access & go to #4

4. Show screen to choose second factor method

5. If user chooses Authy:
    * Send Authy one-touch and wait for reply
    * If confirmed, redirect to dashboard
    * Else sign out & return to sign-in screen

6. If user chooses SMS:
    * Generate SMS & store encrypted in DB w/expiration time
    * Send SMS code
    * Go to 8

7. If user chooses backup code
    * Show page with backup code entry box

8. Show page with code entry box
	* Include hidden field with auth type: SMS or backup code
    * If SMS, page should include a link to generate another code (timeout, not received, etc)

9. Verify code
    * If invalid, re-display page with error message.
        * On fourth invalid attempt, lock account for 5 minutes.
    * If valid, clear restricted access flag & redirect to dashboard.


## Restrictions

Restrict users awaiting two-factor token verification to the following pages:

| URL                    | Description |
| ---------------------- | -------------------------------------------------- |
| /users/2fa\_backup     | Enter backup code                                  |
| /users/2fa\_choose     | Choose second factor method                        |
| /users/2fa\_sms        | Send SMS code, show enter auth token form          |
| /users/2fa\_push       | Push Authy one-touch and wait for response         |
| /users/2fa\_resend     | Resend SMS code                                    |
| /users/2fa\_verify     | Verify SMS/Backup code                             |

These users are marked by user.AwaitingSecondFactor = true.

If visiting unauthorized page, redirect to Choose Second factor page. This is in middleware/authenticate.go

## Two-Factor Setup

* Add enable/disable two-factor auth button to My Account
* When switching to enable, choose two-factor method as described below
* When switching to disable, confirm 2fa is disabled, but don't change any other settings

### Choose Two Factor Auth Method

* SMS
  * If chosen, have user enter and confirm phone number
* Authy
  * If chosen, have user enter phone and then we'll register them with Authy
* None (2fa disabled)

### Authy Registration

It looks like Pharos uses Authy OneTouch, not time-based codes. So we shouldn't need to show the user a QR code on Authy setup. Just register the user with Authy's [RegisterUser](https://github.com/dcu/go-authy#creating-users) call.

### Phone Confirmation

Send an initial OTP via SMS and wait for the user to enter it.

### Generate Backup Codes

Generate and display five backup codes. These should be random strings of 10 characters: upper-case, lower-case and digits.

## Misc Notes

If we ever do have to generate QR codes for Authy, here's a [basic overview of the registration process](https://socketloop.com/tutorials/golang-how-to-implement-two-factor-authentication) and here's some [sample code for generating QR codes](https://www.socketloop.com/tutorials/golang-generate-qr-codes-for-google-authenticator-app-and-fix-cannot-interpret-qr-code-error).

This gist provides a simple example of [how to send a text message](https://gist.github.com/BizarroDavid/40f644de19a93039de5e67439de704b4). The main SNS library, in all its horror, is [documented here](https://docs.aws.amazon.com/sdk-for-go/api/service/sns/).

For logging purposes, see the documentation on [Publish](https://docs.aws.amazon.com/sdk-for-go/api/service/sns/#SNS.Publish) and [PublishOutput](https://docs.aws.amazon.com/sdk-for-go/api/service/sns/#PublishOutput). It would be nice if we could link the message ID of the publish output to the CloudTrail log entry that describes the message's disposition. that would simplify the process of tracing problematic texts.
