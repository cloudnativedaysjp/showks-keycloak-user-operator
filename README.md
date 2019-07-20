# showks-keycloak-user-operator

## install

```
$ cat ./config/secret.env
KEYCLOAK_BASE_PATH=http://example.com:8080/
KEYCLOAK_USERNAME=xxxxxxxxxx
KEYCLOAK_PASSWORD=xxxxxxxxxx
KEYCLOAK_REALM=xxxxxxxxxx
```

```
$ make deploy
```

## Usage

```yaml
apiVersion: showks.cloudnativedays.jp/v1beta1
kind: KeyCloakUser
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: keycloakuser-sample
spec:
  username: alice
  passwordSecretName: XXXXXX
  realm: master
```

## Development

KeyCloakのAPIを実行するため、環境変数`KEYCLOAK_BASE_PATH`と`KEYCLOAK_USERNAME`、`KEYCLOAK_PASSWORD`、`KEYCLOAK_REALM`をセットします。

その後、以下のようにコントローラーを手元で実行します。

```
$ kubectl apply -f config/crds
$ make run
```


aaa
