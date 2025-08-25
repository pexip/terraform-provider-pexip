Your job is to port terraform-provider-pexip away from using github.com/hashicorp/terraform-plugin-sdk/v2 in favor of using github.com/hashicorp/terraform-plugin-framework and maintain the repository.

You have access to the current terraform-provider-pexip repository as well as the ../go-infinity-sdk repository.

Make a commit to the current branch after every single file edit to a local branch called plugin-library-porting.

Use the terraform-provider-pexip/.agent/ directory as a scratchpad for your work. Store long term plans and todo lists there.

The original project was mostly tested by manually running the code. When porting, you will need to write end to end and unit tests for the project. But make sure to spend most of your time on the actual porting, not on the testing. A good heuristic is to spend 80% of your time on the actual porting, and 20% on the
testing.
