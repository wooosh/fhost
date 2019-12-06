# fhost
A small and simple filehost written in go.

# Installation
```
go get github.com/wooosh/fhost
mkdir files # Stores all files uploaded
fhost # Must be root to serve on port 80
```

# Uploading files
```
curl -F'file=filename' myfilehost.com
```

# Configuration
fhost supports configuration via the following environment variables:

| Variable | Purpose                                                                                                                                                         |
|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| AUTH     | Sets the password. Use curl's -u flag when uploading: ``` curl -F'file=@foo.txt' -u ':mypassword' filehost.com # A colon must be prepended to your password ``` |
| PORT     | The port the filehost will accept uploads and receive files on                                                                                                  |
| WEBPATH  | The root of the filehost used to give an address back like `mysite.com/filehost/1IDeUw`                                                                         |

> NOTE: the authentication system used is HTTP basic auth which is NOT secure
