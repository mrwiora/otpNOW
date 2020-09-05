This tutorial has Debian 10 as the only preprequirement. Most probably any freeradius installation will work...

Ensure that the required packages are installed:
```shell
apt install freeradius freeradius-rest -y
```

Stop the freeradius server as we will change the configuration and perform a test run in debug mode:
```shell
systemctl stop freeradius
```

Enable the rest  module for freeradius:
```shell
ln -s /etc/freeradius/3.0/mods-available/rest /etc/freeradius/3.0/mods-enabled/rest
chown freerad.freerad /etc/freeradius/3.0/mods-enabled/rest
```

Set the default auth type to Rest in `/etc/freeradius/3.0/users`- e.g. in first line of the file:
```nginx
DEFAULT        Auth-Type := rest
```

Make the authentication type "rest" available in `/etc/freeradius/3.0/sites-enabled/default`:
search for the authenticate section and add:
```
authenticate {
        Auth-Type rest {
                rest
        }
        ...
}
```

Customize the rest plugin to ask the otpNOW Rest Service in `/etc/freeradius/3.0/mods-enabled/rest`:
Note: Take care to NOT have an ending / on the connect_uri...
```
        connect_uri = "http://54.74.185.110:8080"
        ...
        authenticate {
                uri = "${..connect_uri}/totp?username=%{User-Name}&passcode=%{User-Password}"
                method = 'get'
                tls = ${..tls}
        }
```

Run in debug mode to see if it is basically working - in case something goes wrong you will find details here
```shell
/usr/sbin/freeradius -X
```
The server should come up with:
```shell
...
Listening on auth address * port 1812 bound to server default
Listening on acct address * port 1813 bound to server default
Listening on auth address :: port 1812 bound to server default
Listening on acct address :: port 1813 bound to server default
Listening on auth address 127.0.0.1 port 18120 bound to server inner-tunnel
Listening on proxy address * port 34846
Listening on proxy address :: port 46578
Ready to process requests
```

Use the radclient to verify that a not valid token is correctly rejected:
```shell
echo "User-Name=test, User-Password=123456" | radclient -sx 127.0.0.1 auth testing123
```
In the debug output you should find something like:
```shell
Ready to process requests
(0) Received Access-Request Id 207 from 127.0.0.1:42377 to 127.0.0.1:1812 length 44
(0)   User-Name = "test"
(0)   User-Password = "123456"
...
(0) Sent Access-Reject Id 207 from 127.0.0.1:1812 to 127.0.0.1:42377 length 20
Waking up in 3.9 seconds.
(0) Cleaning up request packet ID 207 with timestamp +6
Ready to process requests
```
On the Client side you should get a:
```shell
$ echo "User-Name=test, User-Password=123456" | radclient -sx 127.0.0.1 auth testing123
Sent Access-Request Id 207 from 0.0.0.0:42377 to 127.0.0.1:1812 length 44
        User-Name = "test"
        User-Password = "123456"
        Cleartext-Password = "123456"
Received Access-Reject Id 207 from 127.0.0.1:1812 to 127.0.0.1:42377 length 20
(0) -: Expected Access-Accept got Access-Reject
Packet summary:
        Accepted      : 0
        Rejected      : 1
        Lost          : 0
        Passed filter : 0
        Failed filter : 1
```

Do the same with a freshly generated token:
```shell
echo "User-Name=test, User-Password=092451" | radclient -sx 127.0.0.1 auth testing123
```
In the debug output you should find something like:
```shell
Ready to process requests
(2) Received Access-Request Id 168 from 127.0.0.1:50144 to 127.0.0.1:1812 length 44
(2)   User-Name = "test"
(2)   User-Password = "092451"
...
(2) Sent Access-Accept Id 168 from 127.0.0.1:1812 to 127.0.0.1:50144 length 0
(2) Finished request
Waking up in 4.9 seconds.
(2) Cleaning up request packet ID 168 with timestamp +18
Ready to process requests
```
On the Client side you should get a:
```shell
$ echo "User-Name=test, User-Password=092451" | radclient -sx 127.0.0.1 auth testing123
Sent Access-Request Id 168 from 0.0.0.0:50144 to 127.0.0.1:1812 length 44
        User-Name = "test"
        User-Password = "092451"
        Cleartext-Password = "092451"
Received Access-Accept Id 168 from 127.0.0.1:1812 to 127.0.0.1:50144 length 20
Packet summary:
        Accepted      : 1
        Rejected      : 0
        Lost          : 0
        Passed filter : 1
        Failed filter : 0
```


In case your tests were successful, just stop the debug instance (Ctrl+C) and restart the service:
```shell
systemctl start freeradius
```

The most basic logs will be now available in `/var/log/freeradius/radius.log`

Set up your own otpNOW Server and enjoy :)