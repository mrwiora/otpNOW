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
