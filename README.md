# Device42 Terraform Provider

## Disclaimer
It is my last first attempt actually to start coding. Please do not use this provider in production; it is currently very unstable.

## How to build
```
git clone https://github.com/Deriul/terraform-provider-device42.git
cd terraform-provider-device42 && go build && mv main ~/.terraform.d/plugins/linux_amd64/terraform-provider-device42
```

## How to use

[Device42 API reference](https://api.device42.com/#/)

```
provider "device42" {
  address = "https://10.10.10.10"
  username = "admin"
  password = "adm!nd42"
}
```

This plugin serves the one only purpose: it creates child subnetworks under parent subnets in Device42. 
To use it, you must have at least one parent network in Device42 pre-created. The network must be in the VRF group, not Assigned and not Allocated.
Additionally, during subnet creation, it sets the last available IP address in the child network as a gateway.

### Data
```
data "device42_subnet" "parent01" {
  name = "parent_name_01"
}
```

You might want to get a parent network ID by name. It is self-explanatory.

### Resources
```
resource "device42_subnet" "child01" {
  name = "child_name_01"
  parent_subnet_id = data.device42_subnet.parent01.subnet_id
  mask_bits = 28
}
```
**name** - (Required) Although it is not required by the API, it should be set (fight entrophy!).<br>
**parent_subnet_id** - (Required) String. Self-explainatory.<br>
**mask_bits** - (Required) Int. Required network size.

Parameter name could be changed later, size and parent are unchangeable. If you try to change it, the provider will create a new network (don't do it).

## Known issues

- You can not create more than one child from the same parent at a time. If you need more than one, create one at a time. This behaviour should be easy to fix once I have more free time.
- Changing the existing resource's parent_subnet_id or mask_bits should raise an error. I will do it later.
- I should add the child name uniqueness check.