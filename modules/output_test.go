/*
Copyright © 2024 PATRICK HERMANN PATRICK.HERMANN@SVA.DE
*/
package modules

import (
	"testing"
)

var rawSecretManifest = `apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm`

func TestHandleOutput(t *testing.T) {

	HandleOutput("stdout", "", rawSecretManifest)
	HandleOutput("file", "/tmp/file.yaml", rawSecretManifest)

}
