# stuttgart-things/kaeffken

<div align="center">
  <p>
    <img src="https://github.com/stuttgart-things/docs/blob/main/hugo/sthings-coffee.png" alt="sthings" width="350" />
  </p>
  <p>
    <strong>[ˈkɛfkən]</strong>- gitops apps & cluster management cli

  </p>
</div>

## COMMANDS

<details><summary><b>CREATE</b></summary>

```bash
kaeffken create --survey=true --questions "tests/questions1.yaml,tests/questions2.yaml"
```

</details>

<details><summary><b>DECRYPT</b></summary>

## DECRYPT FILE (SOPS)

### STDOUT

```bash
export SOPS_AGE_KEY=AGE-SECRET-KEY-1T22K05UTRDU..
kaeffken decrypt \
--source encrypted.yaml
```

### TO FILE

```bash
kaeffken decrypt \
--source encrypted.yaml \
--key AGE-SECRET-KEY-1T22K05UTRDU.. \
--output file \
--destination /tmp/decrypted.yaml
```

</details>


<details><summary><b>APPS</b></summary>

```bash
kaeffken \
--output file \
--clusterPath=clusters/labul/test1 \
--apps tests/apps.yaml
--pr true
```

```bash
kaeffken \
--output stdout \
--apps /home/sthings/projects/stuttgart-things/kaeffken/apps/flux/apps.yaml \
--appDefaults /home/sthings/projects/stuttgart-things/kaeffken/apps/flux/app-defaults.yaml \
--defaults /home/sthings/projects/stuttgart-things/kaeffken/apps/flux/flux-defaults.yaml
```

</details>

<details><summary><b>ENCRYPT FILE</b></summary>

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

```bash
kaeffken encrypt \
--source tests/secret.yaml \
--output file \
--destination ~/projects/sops/ \
--name config \
--age age1g438...
```

```bash
kaeffken encrypt \
--source tests/secret.yaml \
--output file \
--pr true \
--destination /tmp \
--clusterPath=clusters/labul/test1
```

</details>

<details><summary><b>RENDER (BUILTIN) TEMPLATE AND ENCRYPT FILE</b></summary>

```bash
kaeffken encrypt \
--template k8s \
--values "password=mysecretvalue, username=admin" \
--output stdout
```

</details>

## INSTALL

<details><summary><b>INSTALL (LINUX)</b></summary>

</details>

## DEV

<details><summary><b>CREATE BRANCH</b></summary>

```bash
task branch
```

</details>

:bulb: :computer: :floppy_disk: Add Features, fixes, documentation ...

<details><summary><b>LINT, TEST, BUILD, RUN</b></summary>

```bash
task run
```

</details>

<details><summary><b>CREATE/MERGE PULL REQUEST</b></summary>

```bash
task pr
```

</details>

<details><summary><b>RELEASE VERSIONED ARTIFACTS</b></summary>

```bash
task release
```

</details>

<details><summary><b>ENV FILE</b></summary>

```bash
cat <<EOF > .env
SOPS_AGE_KEY=AGE-SECRET-KEY-1T2...
EOF
```

</details>

<details><summary><b>ALL TASKS</b></summary>

```bash
task: Available tasks for this project:
* branch:              Create branch from main
* build:               Install
* build-ko:            Build KO Image
* commit:              Commit + push code into branch
* delete-branch:       Delete branch from origin
* lint:                Lint code
* pr:                  Create pull request into main
* release:             Release
* run:                 Run
* test:                Test code
* tests:               Built cli tests
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
