terraform {
  backend "s3" {
    endpoints = {
      s3 = "{{ .s3Endpoint }}"
    }
}
