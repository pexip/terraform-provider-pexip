/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

# TODO: add nodes to this test case
resource "pexip_infinity_ssh_authorized_key" "tf-test-ssh-key" {
  key     = "AAAAB3NzaC1yc2EAAAADAQABAAACAQDGWIGMczXIMRNassH/IuFPSoyryEyn3uhUqn1s3tSSDOV0b3xogwejZJZKZfUo+oFoYKLbeD70CuZCSIHOx5uZmTYk04vN8r4fX0nzEfHYSCty5ZSvPXevdxyZD+CLnTEtYxbBq4k3xIsmprRKWz70MoVXQqM9jZpR5sOc1LarW24HJhM22iVVghrDX6tsI13Kvld3QRg6Y+jh6rZnH8k3EBwqP+BndSp4ECUM+XA5OEFN4ylZSlk/VS6V9XcVnERFbA3m+qkIhx/K8dc5XmGDGO1Aayn78z2lBtdUul4YdQnUYczu6hpJa2Swasatip0CL6o3vJX344MwkU3MMzJ+ynPdOMOLqQjFgX1gNboWa5udNNdzKdLmRYd3//Fwx9ZE6lPlPrApb6C1VZNgqvFl7yz0F0eSVOJZ7iEL6WzYybbtPbrWi0kO5bYpB/muP2jficXwCqaVxG9Qj/at6ALGPAgkZWbLh0MZFlH0fQzQYxnq2aLRe0KPdgoWXOW1gU7fycR/0j28yBYX5XAI1DMvwB+6vONuEo27Ty6etwHHJWYpVzmzwoElBcqfeBRxtdAgB4Rbq+SX3kNE4J5bsWxY0D6UkUuZ0xdRAgjcRwWxcJsTwIMKSyjUWzoihtIaANQE2sX/6LuR8xI3tI7ckSpY9QZzci6W/o6PuOWeP/Njkw=="
  keytype = "ssh-rsa"
  comment = "tf-test SSH Key"
  #nodes   = ["/api/admin/configuration/v1/worker_vm/1/"]
}
