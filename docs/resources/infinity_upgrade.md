# infinity_upgrade

Triggers a system upgrade on the Infinity service. This is an action resource that initiates an upgrade process. Note: This resource only supports creation and reading - upgrades cannot be updated or undone once initiated. The resource represents the upgrade trigger action, not the upgrade state itself.

## Example Usage

```hcl
resource "pexip_infinity_upgrade" "example" {
  package = "pexip-infinity-29.3.1"
}
```

## Argument Reference

The following arguments are supported:

* `package` - (Optional) Specific upgrade package to use. If not specified, the system will use the default upgrade package.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier for this upgrade trigger.
* `timestamp` - Timestamp when the upgrade was triggered.

## Important Notes

- This is an **action resource** - it triggers an upgrade when created
- Upgrades **cannot be undone** once initiated
- The resource does not track the actual upgrade progress or status
- Updates and deletions are not supported for this resource
- Each apply will trigger a new upgrade if the resource is recreated
- Use this resource with caution in production environments

## Upgrade Process

When this resource is created:
1. It triggers an upgrade process on the Pexip Infinity system
2. The upgrade runs asynchronously in the background
3. The resource stores a timestamp of when the upgrade was triggered
4. Monitor the system separately to track upgrade progress and completion

## Warning

⚠️ **Use with extreme caution**: This resource will immediately trigger a system upgrade which may cause:
- Service interruption during the upgrade process
- Potential system instability if the upgrade fails
- Irreversible changes to the system configuration

Always test upgrades in a development environment before applying to production.