# stuttgart-things/kaeffken
gitops cluster management cli

<details><summary><b>DEPLOY</b></summary>

```bash
kaeffken --name [CLUSTERNAME]
#kaeffken --name michigan --env labul
```

</details>

<details><summary><b>ENCRYPT</b></summary>

```bash
cat <<EOF >> tests/secret.yaml
kind: Secret
apiVersion: v1
metadata:
  name: secret
data:
  password: wHat6ver
EOF
```

```bash
kaeffken encrypt \
--source tests/secret.yaml \
--output stdout
```

</details>

## AUTHOR

```bash
Patrick Hermann, stuttgart-things 12/2023
```

## License

Licensed under the Apache License, Version 2.0 (the "License").

You may obtain a copy of the License at [apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0).

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an _"AS IS"_ basis, without WARRANTIES or conditions of any kind, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
