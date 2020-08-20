provider "device42" {
  address = "https://10.10.10.10"
  username = "admin"
  password = "adm!nd42"
}

data "device42_subnet" "parent01" {
  name = "parent_name_01"
}

data "device42_subnet" "parent02" {
  name = "parent_name_02"
}

resource "device42_subnet" "child01" {
  name = "child_name_01"
  parent_subnet_id = data.device42_subnet.parent01.subnet_id
  mask_bits = 28
}

resource "device42_subnet" "child02" {
  name = "child_name_02"
  parent_subnet_id = data.device42_subnet.parent02.subnet_id
  mask_bits = 29
}

