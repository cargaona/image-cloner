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
5. It's also recommended creating a secret with the container registry credentials and attach it to the serviceAccount in order to give access to the resources to your new registry. You can find info about it [here.](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account)
6. You can deploy all the content on ./manifests manually or just run ```make deploy``` if you are logged to your cluster in the right namespaces.

### Notes
Since this is just an MVP, there are some things you need to have in mind.

- The controller needs the Cluster Admin to provision the container registry credentials for each namespace. Create the secret with the credentials and add it to the default serviceAccount it's a good way to start. If you don't want to do that change on every namespace, you can use some tools like [imagepullsecret-patcher](https://github.com/titansoft-pte-ltd/imagepullsecret-patcher) or [kubernetes-reflector](https://github.com/EmberStack/kubernetes-reflector)
- It's not compatible with [Schema v1 images](https://docs.docker.com/registry/spec/deprecated-schema-v1/).

