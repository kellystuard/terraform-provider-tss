# Contributing

Thanks for considering contributing to the Thycotic Secret Server Terraform Provider!

## :inbox_tray: Opening issues

If you find a bug, please feel free to [open an issue](https://github.com/kellystuard/terraform-provider-tss/issues).

If you taking the time to mention a problem, even a seemingly minor one, it is greatly appreciated, and a totally valid contribution to this project. Thank you!

## :beetle: Fixing bugs

We love pull requests. Here’s a quick guide:

1. [Fork this repository](https://github.com/kellystuard/terraform-provider-tss/fork) and then clone it locally:

  ```bash
  git clone https://github.com/[your-name-here]/terraform-provider-tss
  ```

2. Create a topic branch for your changes:

  ```bash
  git checkout -b fix-for-that-thing
  ```
3. Commit a failing test for the bug:

  ```bash
  git commit -am "Adds a failing test to demonstrate that thing"
  ```

4. Commit a fix that makes the test pass:

  ```bash
  git commit -am "Adds a fix for that thing!"
  ```

5. Run the tests:

  ```bash
  make test
  ```

6. If everything looks good, push to your fork:

  ```bash
  git push origin fix-for-that-thing
  ```

7. [Submit a pull request.](https://help.github.com/articles/creating-a-pull-request)

8. Acceptance tests will automatically run on the pull request and let you know if anything got broken.

## :love_letter: Adding new features

Thinking of adding a new feature? [Open an issue](https://github.com/kellystuard/terraform-provider-tss/issues).
