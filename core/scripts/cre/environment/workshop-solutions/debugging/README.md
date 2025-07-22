# Workshop Debugging Exercise

This section is intended to help you understand how to use smoke tests and the local CRE env together.

We are leveraging the existing [workshop guide](../workflow/README.md) to implement, and debug a capability update that breaks the workflow.

We are introducing a broken capability update that you will need to debug and fix. We'll walk you through it!

This contains a script (`script.sh`) which provides a set of commands to help you with this.

# Prerequisites
You can see all the prerequisites [here](https://docs.google.com/document/d/1HtVLv2ipx2jvU15WYOijQ-R-5BIZrTdAaumlquQVZ48).

Make sure you copy this script to the root of your workspace, in which you have cloned both the capabilites and chainlink repositories:
```bash
cp ./script.sh ~/workspace/cre-workshop/

cre-workshop/
├── capabilities/
├── chainlink/
├── script.sh (you can rename this to whatever you want(e.g `debugging.sh`))
````

Just make sure the script name is not conflicting with the one you copied from the [workflow exercise](../workflow/README.md), 
in case you'd like to keep it there.

Give the script executable permissions:
```bash
chmod +x script.sh
```

You are now good to go.

# Step-by-Step Guide

Stay on your modified `dx-1407-workshop-cre` branch in `chainlink` repository.

Checkout the `dx-1343-broken-capability` branch of `capabilities` repository and build the update:
```bash
./script.sh checkout_broken_capability
```

Restart the local CRE env and the observability stack:
```bash
./script.sh start_observability_stack && ./script.sh restart_local_cre
```

Compile and upload the workflow:
```bash
./script.sh deploy_example_workflow ### the script referenced here is the one from the [workflow exercise](../workflow/README.md)
```

Open your [Local Red Panda](http://localhost:8080/topics/cre?p=-1&s=500&o=-3#messages) and try to find your log message (hint: it won't be there 😉).

Open the Grafana [workflow-engine](http://localhost:3000/d/ce589a98-b4be-4f80-bed1-bc62f3e4414a/workflow-engine?orgId=1&refresh=5s&from=now-5m&to=now) 
dashboard (previously started by `start_observability_stack`) and look at the Workflow Errors. Investigate errors with a similar message to `failed to register trigger`.

We have a skeleton smoke test the will check for successful execution. You need to fill in some details and run the test. It will fail until you fix it.

Let's save some resources and stop the local CRE environment:
```bash
./script.sh shutdown_local_cre
```

Edit the `Test_V2_Workflow_Workshop` test in the `chainlink/system-tests/tests/smoke/cre/workshop_test.go` file to run 
your workflow (line 119).

Run the test to confirm the failure:
```bash
./script.sh run_workflow_test
```

Observe the failure and check the dashboards to confirm the error.

The capability schedule is set in `cron/trigger/trigger.go:30` file in the `capabilities` repository.

Changes to the capability code will require yoiu to compile the capability binary again. You can do this with:
```bash
./script.sh compile_cron_capability
```

There is a second bug that you will stumble over that also has to be fixed...

We are checking [Chainlink Node Logs dashboard](http://localhost:3000/d/a7de535b-3e0f-4066-bed7-d505b6ec9ef1/cl-node-errors?orgId=1&refresh=5s) for clues.

To fix this issue, you need to open `cron/trigger/trigger.go` again and delete this incorrect line, which is trying to call a function on nil pointer:
```go
s.lggr.Infof("schedule error: %s", err.Error())
```

Let's compile the capability binary again:
```bash
./script.sh compile_cron_capability
```

And run the test again:
```bash
./script.sh run_workflow_test
```

## Usage

```bash
./script.sh <command>
```

## Available Commands

- `checkout_broken_capability`  
  Checks out a branch with a broken capability and builds it for testing purposes.

- `start_observability_stack`  
  Starts the observability stack using the `ctf obs u` command.

- `restart_local_cre`  
  Restarts the local CRE environment and the Beholder service.

- `shutdown_local_cre`  
  Shuts down the local CRE environment.

- `run_workflow_test`  
  Runs the workflow test (`Test_V2_Workflow_Workshop`).

- `compile_cron_capability`  
  Compiles the cron capability binary.

## Example

```bash
./script.sh checkout_broken_capability
./script.sh start_observability_stack
./script.sh restart_local_cre
./script.sh run_workflow_test
```

## Notes

- Each command is independent and can be run as needed.
- The script will exit with an error message if any step fails.
- Additional arguments after the command will be forwarded to the executed command.

