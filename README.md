# serve-api

A minimal webserver for local development.
Plus a database api (sqlite only for now).

## How to use

`serve-api` will just run a webserver on 0.0.0.0:3000, which serves the local directory.
`serve-api --verbose` does the same, but will also print the path every requests URI.
`serve-api --port <portnumber>` changes the port.
`serve-api --ipaddress <bind-ip>` changes the ipaddress where `serve-api` will bind to.
`serve-api --directory <path>` changes the directory which is served.
`serve-api --help` will display the shortcuts for these flags as well.

API requests should happen in this manner:
```
/api?table=...&columns=id,name&offset=0&limit=10
```
They'll return the results as json.

## How to release

1. Add information about new version to `CHANGELOG.md` & commit.
2. Push latest changes with `git push`.
3. List existing tags with `git tag`.
   ```
   v0.0.1
   v0.0.2
   ```
4. Select the next available tag and apply it with `git tag v0.0.3`.
5. Push tag with `git push origin v0.0.3`.

> If the precondition is not met or the build fails, the action will delete the (remote) tag itself.
> If you messed something up, you can delete local tags with `git tag -d v0.0.3`.
> And you can delete remote tags with `git push --delete origin v0.0.3`.
