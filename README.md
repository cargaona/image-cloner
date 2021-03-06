## Image-cloner-controller

Image-cloner-controller is a Kubernetes controller that copies your public images deployed on your cluster to a container-registry you choose.
The motivation to this controller is to protect your resources from the usage of public images that can be randomly deleted by someone else.
Also, this controller can prevent you from issues like [this one with Docker Pull Limit](https://github.com/docker/hub-feedback/issues/1741).

### Installation
There is a makefile to help with the development of this project and also will serve to install the controller in your cluster.
1. Copy the content of `manifests/` to a new `.manifests/`. 
2. Create a container-registry wherever you wish. (docker-hub / quay.io are the common choices.)
   1. Run docker login with the registry you created.
   2. Look for the credentials on `~/.docker/config.json` You will need it for the setup of the controller and also for the imagePullSecret.
   3. Copy the content of the json to `.manifests/imagePullRegistry/configk8s.json`
   4. Copy the content of the json to the "docker-config" configMap. (Take care of copying just the part of the registry you will use for the controller.)It will look something like this:
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

3. Complete the "image-cloner-configuration" configMap. More info on [Configuration](#Configuration)
4. It's also recommended creating a secret with the container registry credentials and attach it to the serviceAccount in order to give access to the resources to your new registry. You can find info about it [here.](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account)
5. You can deploy all the content on `./manifests` manually or just run ```make install``` if you are logged to your cluster in the right namespace.


### Configuration
There are some configurations you need, and you can do to  **image-cloner-controller**.

#### ConfigMap
These configurations live in the "image-cloner-configuration" configMap.

Details:

| Value | Description | Default | Optional |
| ----  | ---         | ---     | ---      |
|MAX_CONCURRENT_RECONCILES| Max number of concurrent reconcile loops. Since uploading images is a heavy operation it can consumes a lot of your node resources if you don't limit the number of simultaneous uploads. | 5| Yes |
|NAMESPACES_TO_IGNORE | There are some namespaces you will want to ignore and leave it with the existing images. Place the name of the namespaces you want to ignore here comma separated. | kube-proxy (hard default, within the code) | Yes |
|BACKUP_REGISTRY | The direction of your backup registry. Example: quay.io/cargaona |  - | No |

#### Validating webhook.
You can activate validating webhooks for Daemonsets and Deployments just by applying the manifests/webhooks.yaml.
The webhooks work as admission controllers for the updated resources the image-cloner change.

You may want to exclude some namespaces from this admission controller, such as kube-system and the one you are going to install the image-cloner. To accomplish this behaviour you need to label your namespaces just as described below.

```bash
kubectl label namespace <namespace-to-exclude> validate-backups.image-cloner.io=disable
```

### Notes
Since this is just an MVP, there are some things you need to have in mind.

- The controller needs the Cluster Admin to provision the container registry credentials for each namespace. Create the secret with the credentials and add it to the default serviceAccount it's a good way to start. If you don't want to do that change on every namespace, you can use some tools like [imagepullsecret-patcher](https://github.com/titansoft-pte-ltd/imagepullsecret-patcher) or [kubernetes-reflector](https://github.com/EmberStack/kubernetes-reflector)
- It's not compatible with [Schema v1 images](https://docs.docker.com/registry/spec/deprecated-schema-v1/).


### Next features
- Add support to schema v1 images using go-container-registry libraries instead of just using crane.
- Replace current certs for webhooks and start using cert-manager.
- Implement mutating webhook to back up images for newly created Deployments and Daemonsets. 
