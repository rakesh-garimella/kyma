---
title: Create addons repository
type: Details
---

The repository in which you create your own addons must contain at least one `index.yaml` file and have a specific structure, depending on the type of server that exposes your addons.

## The index yaml file

Your remote addons repository can contain many addons defined in index files. Depending on your needs and preferences, you can create one or more index files to categorize your addons. In the index file, provide an entry for every single addon from your addons repository. The index file must have the following structure:
```
apiVersion: v1
entries:
  {addon_name}:
    - name: {addon_name}
      description: {addon_description}
      version: {addon_version}
```

See the example:
```yaml
apiVersion: v1
entries:
  redis:
    - name: redis
      description: Redis service
      version: 0.0.1
```

>**NOTE:** You must place your addons in the same directory where the `index.yaml` file is stored.

## Supported protocols

Expose your addons repository so that you can provide URLs in the [AddonsConfiguration](#custom-resource-addonsconfiguration) (AC) and [ClusterAddonsConfiguration](#custom-resource-clusteraddonsconfiguration) (CAC) custom resources. The Helm Broker supports exposing addons using the following protocols:

<div tabs>
  <details>
  <summary>
  HTTP/HTTPS
  </summary>

>**NOTE:** The HTTP protocol is supported only in `DevelopMode`. To learn more, read [this](#details-registration-rules-using-http-urls) document.

If you want to use an HTTP or HTTPS server, you must compress your addons to `.tgz` files. The repository structure looks as follows:
```
sample-addon-repository
  ├── {addon_x_name}-{addon_x_version}.tgz           # An addon compressed to a .tgz file
  ├── {addon_y_name}-{addon_y_version}.tgz        
  ├── ...                                      
  ├── index.yaml                                     # A file which defines available addons
  ├── index-2.yaml                              
  └── ...                                                    
```

See the example of the Kyma `addons` repository [here](https://github.com/kyma-project/addons/releases).

>**TIP:** If you contribute to the Kyma [`addons`](https://github.com/kyma-project/addons/tree/master/addons) repository, you do not have to compress your addons as the system does it automatically.

These are the allowed addon repository URLs provided in CAC or AC custom resources for HTTP or HTTPS servers:
```yaml
apiVersion: addons.kyma-project.io/v1alpha1
kind: ClusterAddonsConfiguration
metadata:
  name: addons-cfg-sample
spec:
  repositories:
    # HTTPS protocol
    - url: "https://github.com/kyma-project/addons/releases/download/latest/index.yaml"
    # HTTP protocol
    - url: "http://github.com/kyma-project/addons/releases/download/latest/index.yaml"
```

  </details>
  <details>
  <summary>
  Git
  </summary>

If you want to use Git, place your addons directly in addons directories. The repository structure looks as follows:
```
sample-addon-repository
  ├── {addon_x_name}-{addon_x_version}               # An addon directory
  ├── {addon_y_name}-{addon_y_version}        
  ├── ...                                      
  ├── index.yaml                                     # A file which defines available addons
  ├── index-2.yaml                              
  └── ...                                                    
```

See the example of the Kyma `addons` repository [here](https://github.com/kyma-project/addons/tree/master/addons).


You can specify a Git repository URL by adding a special `git::` prefix to the URL address. After this prefix, provide any valid Git URL with one of the protocols supported by Git. In the URL, you can specify a branch, commit, or tag version. You can also add the `depth` query parameter with a number that specifies the last revision you want to clone from the repository.

>**NOTE:** If you use `depth` together with `ref`, make sure that `depth` number is big enough to clone a proper reference. For example, if you have `depth=1` and `ref` set to a commit from the distant past, the URL will not work as you clone only the first commit from the `master` branch and there is no option to do the checkout.

These are the allowed addon repository URLs provided in CAC or AC custom resources for Git:
```yaml
apiVersion: addons.kyma-project.io/v1alpha1
kind: ClusterAddonsConfiguration
metadata:
  name: addons-cfg-sample
spec:
  repositories:
    # Git HTTPS protocol with a path to index.yaml
    - url: "git::https://github.com/kyma-project/addons.git//addons/index.yaml"
    # Git HTTPS protocol with a path to index.yaml of a specified version and a depth query parameter
    - url: "git::https://github.com/kyma-project/addons.git//addons/index.yaml?ref=1.2.0&depth=3"
    # github.com URL with no prefix. It is automatically interpreted as a Git repository source.
    - url: "github.com/kyma-project/addons//addons/index.yaml"
    # bitbucket.org URL with no prefix. It is automatically interpreted as a Git repository source.
    - url: "bitbucket.org/kyma-project/addons//addons/index.yaml"
```

>**NOTE:** For now, the SSH protocol is not supported.

  </details>
</div>
