# Workshop Exercise

This section is intended to help you set up and manage the local Chainlink CRE (Chainlink Runtime Environment) workshop 
environment. It provides a script that automates various tasks related to the workshop, such as preparing repositories, 
starting the local CRE, deploying workflows, etc.

The main goal is to become familiar with Workflow development, and while doing so, you will learn how 
to add a custom user log to a V2 workflow.

This contains a script (`script.sh`) which provides a set of commands to manage the local Chainlink CRE 
(Chainlink Runtime Environment) workshop environment, including repository preparation, environment management, and workflow deployment.

# Prerequisites
You can see all the prerequisites [here](https://docs.google.com/document/d/1HtVLv2ipx2jvU15WYOijQ-R-5BIZrTdAaumlquQVZ48).

Make sure you copy this script to the root of your workspace, in which you have cloned both the capabilites and chainlink repositories:
```bash
cp ./script.sh ~/workspace/cre-workshop/

cre-workshop/
├── capabilities/
├── chainlink/
├── script.sh (you can rename this to whatever you want(e.g `workshop.sh`))
````

Give the script executable permissions:
```bash
chmod +x script.sh
```

You are now good to go.

# Step-by-Step Guide

Stop any running Docker containers to avoid potential port collisions:
```bash
./script.sh stop_docker
```

Checkout the `dx-1407-workshop-cre` branch in the `capabilities` repository and compile the `cron` capability binary:
```bash
./script.sh prepare_capabilities_repo
```

Checkout the branch `dx-1407-workshop-cre` in chainlink repository (if not already checked out), and start the local CRE environment:
```bash
./script.sh prepare_core_repo && ./script.sh start_local_cre
```

Open your editor/IDE and navigate to the `chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/main.go` file 
and add a log inside the `onTrigger()` function. Something like:
```go
wcx.Logger.Info(“<YOUR MESSAGE GOES HERE”>) // note: use a string literal
```

Compile and upload the workflow:
```bash
./script.sh deploy_example_workflow
```

Open the `cre` topic view in your [Local Red Panda](http://localhost:8080/topics/cre?p=-1&s=500&o=-3#messages).

Create a JS filter to display your custom log (this will be enabled by default after creation):
```javascript
var msg = value.logLines[0].message
return msg.includes("<YOUR LOG MSG HERE>")
```

Set `START OFFSET` as `Newest - 500` and you should see rows with your message.

## Usage

```bash
./script.sh <command>
```

## Available Commands

- `stop_docker`  
  Stops all running Docker containers on your system (to avoid potential port collisions).

- `prepare_capabilities_repo`  
  Prepares the capabilities repository for the workshop by checking out the correct branch and building the cron capability.

- `prepare_core_repo`  
  Prepares the core repository for the workshop by checking out the correct branch.

- `start_local_cre`  
  Starts the local CRE environment and the Beholder service.

- `run_example_workflow`  
  Starts the local CRE environment and runs the example workflow with the Beholder.

- `deploy_example_workflow`  
  Deploys the example workflow to the local CRE environment.

- `shutdown_local_cre`  
  Shuts down the local CRE environment.

## Example

```bash
./script.sh prepare_core_repo
./script.sh prepare_capabilities_repo
./script.sh start_local_cre
./script.sh deploy_example_workflow
```

## Notes

- Each command is independent and can be run as needed.
- The script will exit with an error message if any step fails.
- Additional arguments after the command will be forwarded to the executed command.

