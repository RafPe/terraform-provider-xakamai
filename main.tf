provider "xakamai" {
  edgerc  = "~/.edgerc"
  section = "dummy"
}

resource "xakamai_network_list" "my-server" {
  name  = "MY_TF_TEST_1"
  type  = "IP"
  items = ["1.2.3.3/32", "9.8.7.6/24"]
}
