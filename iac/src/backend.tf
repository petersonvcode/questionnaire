terraform {
  # Backend configuration is stored separately at ../config/<env>.conf
  backend "s3" {
    use_lockfile = true
  }
}
