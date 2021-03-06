---
title: Application YAML
weight: 100
---

The app configuration is a file called `qbec.yaml` and needs to be at the root of the directory tree.

```yaml
apiVersion: qbec.io/v1alpha1 # only supported version currently
kind: App # must always be "App"
metadata:
  name: my-app # app name. Allows multiple qbec apps to deploy different objects to the same namespace without GC collisions
spec:
  componentsDir: components    # directory where component files can be found. Not recursive. default: components
  paramsFile: params.libsonnet # file to load for `param list` and `param diff` commands. Not otherwise used.
  postProcessor: pp.jsonnet    # post processor file for injecting common metadata

  # additional library paths when executing jsonnet, no support currently for `http` URLs.
  libPaths:
  - additional
  - local
  - library
  - paths

  # list of components to exclude by default
  excludes:
  - default
  - excluded
  - components
 
  # when dealing with apps that deploy to multiple namespaces, use object lists at cluster scope for GC purposes.
  # by default, this will use namespaced queries for each namespace.
  clusterScopedLists: true

  # declaration of late-bound variable definitions that can be passed in on the command line using the --vm:* options. 
  vars:
    # external variables are accessed as std.extVar('var-name')
    external:
      - name: imageTag # the name of the external variable passed in using --vm:ext-str and related options
        default: 'latest' # the default value to use if this variable is not specified on the command line. Can be an arbitrary object.
        secret: false # when true qbec will not print the plain text value in any debug message

    # for top-level variables, your component's main "object" is a function that accepts a value, typically 
    # initialized with a default in the code.
    topLevel:
      - name: mySecret # the name of the top-level variable (i.e. the name of the function parameter of your component's function)
        components: [ 'service2' ] # the components that require this variable. Must be specified.
        secret: true

  # if the following attribute is set to true and the --app-tag argument is set on the command line, qbec will automatically
  # change the default namespace for the environment in question by suffixing it with <hyphen><tag-value> (e.g. 'myns-tag')
  namespaceTagSuffix: true

  # an arbitrary object to define baseline properties
  baseProperties:
    foo: base

  # map of environment names to environment objects. An environment is the combination of a server URL and default
  # namespace. The default namespace is used for objects that do not have an explicit namespace set.
  environments:
    # for enviromnents such as minikube that do not always have a stable server URL you can use a context name instead.
    minikube:
      context: minikube # named context, prefer server URLs for non-local clusters.
      defaultNamespace: my-ns # the namespace to use when namespaced object does not define it.
      includes: # components to include, subset of global exclusion list
      - components
      - to
      - include
      excludes: # additional components to exclude
      - more
      - exclusions

    dev:
      server: https://dev-server # server URL
      properties: # arbitrary properties can be attached to environments
        foo: bar

  # additional environments can be loaded from files. Files are loaded in the order specified.
  # It is explicitly allowed for a later file to replace an inline environment or one loaded from an earlier file.
  # The file path is relative to the directory where qbec.yaml resides. http(s) URLs and glob patterns are also supported
  envFiles:
  - more-envs.yaml
  - https://my.server/envs.yaml
  - envs/*.yaml

  # if the following attribute is set to true, qbec will add component names also as labels to Kubernetes objects. 
  addComponentLabel: true
```

### Environment files

Environments can be defined in external files that are then loaded and merged into the main environments object.

```yaml
apiVersion: qbec.io/v1alpha1
kind: EnvironmentMap
spec:
  # the environments key is exactly the same as the environments key in qbec.yaml
  environments:
    prod:
      server: https://prod-server
      includes:
      - service2
      properties:
        foo: bar
```

### Notes

* The list of components is loaded from the `componentsDir` directory.
* Once the list is loaded, all exclusion and inclusion lists are checked to ensure that they refer to valid components.
* The global exclusion list allows you to introduce a new component gradually by only including it in a dev environment.
