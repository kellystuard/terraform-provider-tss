terraform {
  required_providers {
    tss = {
      versions = ["0.2"]
      source = "hashicorp.com/edu/tss"
    }
  }
}

provider tss {
}

data tss_secret_field test_username {
  number = 1
  slug   = "username"
}

data tss_secret_field test_password {
  number = 1
  slug   = "password"
}

output test_password {
  value = "${data.tss_secret_field.test_username.value} : ${data.tss_secret_field.test_password.value}"
}
