otpNOW [![Build Status](https://travis-ci.org/mwiora/otpNOW.svg)](https://travis-ci.org/mwiora/otpNOW) [![Code Climate](https://codeclimate.com/github/mwiora/otpNOW/badges/gpa.svg)](https://codeclimate.com/github/mwiora/otpNOW)
=========

requirements and getting started
---------------

* compile yourself (requirement: install go (minimum 1.15.x) as described here https://golang.org/doc/install)
```
go get github.com/pquerna/otp
go get github.com/mwiora/otpNOW/
cd $GOPATH/src/github.com/mwiora/otpNOW/
go build
./otpNOW
```

sample output of current version (debug off)
---------------

```
$ ./otpNOW
Issuer:       Example.com
Account Name: alice@example.com
Secret:       KY3DBQYFXZ4MYW6NSWEGPWTYD7BD4QPY
Writing PNG to qr-code.png....

Please add your TOTP to your OTP Application now!
```

sample interaction with the api
---------------

- Just calling the Server - the offered QR can be scanned by any authenticator app - e.g. Google Authenticator app
```
http://ip:8080/
```
![qr](_resources_readme/qr.png)

- next we could ask the server if a token is valid, which is certainly not
```
http://ip:8080/token?passcode=123456
```
![qr](_resources_readme/check_notok.png)

- next we could ask the server if a generated token is valid
```
http://ip:8080/token?passcode=611201
```
![qr](_resources_readme/check_ok.png)

- doing this by a commandline client - with some more information offering - it looks like:

![qr](_resources_readme/curl.png)

checklist
---------------
basics
- [x] provide reference key to check if it's working
- [x] provide totp validation rest interface
- [ ] provide hotp validation rest interface
- [ ] base user key management storage

integrations
- [x] freeradius guidance

nice2have
- [ ] implement test driven development
