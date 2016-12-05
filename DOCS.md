Use this plugin for caching build artifacts to speed up your build times. This
plugin can create and restore caches of any folders.

## Config

The following parameters are used to configure the plugin:

* **url** - url of the artifactory server
* **username** - authenticate with this user against artifactory server
* **password** - authenticate with this password against artifactory server
* **path** - path on the artifactory server
* **mount** - one or an array of folders to cache
* **rebuild** - boolean flag to trigger a rebuild
* **restore** - boolean flag to trigger a restore

The following secret values can be set to configure the plugin.

* **ARTIFACTORY_CACHE_URL** - corresponds to **url**
* **ARTIFACTORY_CACHE_USERNAME** - corresponds to **username**
* **ARTIFACTORY_CACHE_PASSWORD** - corresponds to **password**

It is highly recommended to put the **ARTIFACTORY_CACHE_USERNAME** and
**ARTIFACTORY_CACHE_PASSWORD** into a secret so it is not exposed to users.
This can be done using the drone-cli.

```bash
drone secret add --image=jmccann/drone-artifactory-cache:latest \
    octocat/hello-world ARTIFACTORY_CACHE_USERNAME octocat

drone secret add --image=jmccann/drone-artifactory-cache:latest \
    octocat/hello-world ARTIFACTORY_CACHE_PASSWORD pa55word
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/secrets/) for additional
information on secrets

## Example

The following is a sample configuration in your `.drone.yml` file:

```yaml
pipeline:
  artifactory_cache_restore:
    image: jmccann/drone-artifactory-cache:latest
    url: https://company.com/artifactory
    path: repo-key/path/archive.tar
    restore: true
    mount:
      - node_modules

  build:
    image: node:latest
    commands:
      - npm install

  artifactory_cache_rebuild:
    image: jmccann/drone-artifactory-cache:latest
    url: https://company.com/artifactory
    path: repo-key/path/archive.tar
    rebuild: true
    mount:
      - node_modules
```
