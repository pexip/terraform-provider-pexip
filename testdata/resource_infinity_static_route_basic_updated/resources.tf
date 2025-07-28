resource "pexip_infinity_static_route" "static_route-test" {
  name = "static_route-test"
  address = "10.0.0.0"  // Updated address
  prefix = 16  // Updated prefix
  gateway = "10.0.0.1"  // Updated value
}