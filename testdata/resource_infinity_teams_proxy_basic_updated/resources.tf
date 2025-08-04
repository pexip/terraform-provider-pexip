resource "pexip_infinity_teams_proxy" "teams-proxy-test" {
  name                    = "test-teams-proxy"
  description             = "Test Teams Proxy"
  address                 = "test-teams-proxy.dev.pexip.network"
  port                    = 8081
  azure_tenant            = "test-azure-tenant"
  eventhub_id             = "test-eventhub-id"
  min_number_of_instances = 1
  notifications_queue     = "test-notifications-queue"
}