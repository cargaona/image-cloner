kubectl create secret generic registrypullsecret --from-file=.dockerconfigjson=/home/char/go/src/github.com/cargaona/image-cloner/configk8s.json --type=kubernetes.io/dockerconfigjson

