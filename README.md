## Image-cloner-controller

Image-cloner-controller is a Kubernetes controller that copies your public images deployed on your cluster to a container-registry you choose.
The motivation to this controller is to protect your resources from the risk of not owned images that can be randomly deleted by someone else.

### Installation
There is a makefile to help with the development of this project and also will serve to install the controller in your cluster.
1. Create a container-registry wherever you wish. (docker-hub / quay.io are the common choices.)
2. Run docker login with the registry you created.
3. Look for the credentials on ~/.docker/config.json. You will need it for the setup of the controller and also for the imagePullSecret.
4. Copy the content of the json to the "docker-config" configMap. (Take care of copying just the part of the registry you will use for the controller.)
   It will look something like this:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
name: docker-configuration
data:
config.json: |
        {
          "auths": {
            "quay.io": {
              "auth": "yourkeys"
            }
          }
        }
```
5. Complete the "image-cloner-configuration" configMap. More info on [Configuration](#Configuration)
5. It's also recommended creating a secret with the container registry credentials and attach it to the serviceAccount in order to give access to the resources to your new registry. You can find info about it [here.](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account)
7. You can deploy all the content on ./manifests manually or just run ```make deploy``` if you are logged to your cluster in the right namespaces.


### Configuration
There are some configurations you need, and you can do to  ##image-cloner-controller##.
These configurations live in the "image-cloner-configuration" configMap.
Details:

| Value | Description | Default | Optional |
| ----  | ---         | ---     | ---      |
|MAX_CONCURRENT_RECONCILES| Max number of concurrent reconcile loops. Since uploading images is a heavy operation it can consumes a lot of your node resources if you don't limit the number of simultaneous uploads. | 5|
|NAMESPACES_TO_IGNORE | There are some namespaces you will want to ignore and leave it with the existing images. Place the name of the namespaces you want to ignore here comma separated. | kube-proxy (hard default, within the code) |
|BACKUP_REGISTRY | The direction of your backup registry. Example: quay.io/cargaona |  - | No |

### Notes
Since this is just an MVP, there are some things you need to have in mind.

- The controller needs the Cluster Admin to provision the container registry credentials for each namespace. Create the secret with the credentials and add it to the default serviceAccount it's a good way to start. If you don't want to do that change on every namespace, you can use some tools like [imagepullsecret-patcher](https://github.com/titansoft-pte-ltd/imagepullsecret-patcher) or [kubernetes-reflector](https://github.com/EmberStack/kubernetes-reflector)
- It's not compatible with [Schema v1 images](https://docs.docker.com/registry/spec/deprecated-schema-v1/).

### Next features
- Implement mutating webhook to change the images of new incoming pods.
- Implement validating webhook to make sure the pods we are deploying have images from our backup registry.
- Activate/Deactivate webhooks feature. It could be and argument or an environment variable passed through confipMap.
- Add support to schema v1 images using go-container-registry libraries instead of just using crane.