---
name: Bug report
about: Create a report to help us improve. Bugs are handled on a best-effort basis
  and help is always appreciated.
title: ''
labels: bug
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**Terraform Version**
Run `terraform -v` to show the version. If you are not running the latest version of Terraform, please upgrade because your issue may have already been fixed.

**Affected Resource(s)**
Please list the resources as a list, for example:
- data.tss_secret_field
If this issue appears to affect multiple resources, it may be an issue with Terraform's core, so please mention this.

**Terraform Configuration Files**
Please include a file that has only the minimal providers/data/resources necessary to understand the issue and no more. Including unrelated configuration could drastically increase the amount of time and effort needed to understand the issue.
```hcl
# Copy-paste your Terraform configurations here.
```

**Debug Output**
If available, please provide a link to a GitHub Gist containing the complete debug output: https://www.terraform.io/docs/internals/debugging.html. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

**Panic Output**
If Terraform produced a panic, please provide a link to a GitHub Gist containing the output of the `crash.log`.

**Expected Behavior**
What should have happened?

**Actual Behavior**
What actually happened?

**Steps to Reproduce**
Steps required to reproduce the issue:
1. `terraform init`
2. `terraform apply`

**Important Factoids**
Are there anything atypical about your accounts that we should know? For example: Running in EC2 Classic? Custom version of OpenStack? Tight ACLs?

**Additional context**
Add any other context about the problem here.
