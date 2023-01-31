# kubecli

This project I have just started as a draft to manage daily activities with kubernetes. Its a very initial draft and might be having errors. Need to create unit test cases as well.

### It depends on kubectl binary

Subcommands 

**swap** :- Generally we add multiple context as connected with different kubernetes clusters. This subcommand can help in changing the context

**authdelete** :- To delete the context run this sub command

**scale**:- scale up and down all the objects together in a namespace
  - Pass filter if need to select specfic set of objects
  - Pass count to specify the count of replicas

**logs**:- read logs of specfic container. It will ask for the selection of pod and container

**selector**:- Its difficult to get the traverse path of nodes in deployment. Run this subcommand and pass `attrib` as most inner node name and it will return the whole path.

**restart**:- select deployment to restart it.

**view**:- return thevalue of any node in the deployment

**backup**:- It will backup the node labels present in all of the deployments within namespace. Backup will be saved in file named as namespace. Necessary to run before remove 

**remove**:- This will remove the required node from all deployments in the namespace. Use filter and it will match only selected deployments. Use attrib to pass the name of node

**apply**:- It will reply the node labels saved in backup file. 
  - Pass file name with --file flag. 
  - If file name is not passed then pass the path of node in --query flag. For eg. --query spec.templates.spec.nodeSelector.nodetype
  - Pass value of node. For eg. --value GPU
  
**delete**:- Delete a pod in a namespace. Run this subcommand and it will ask for the pod. Just specify the number like this
Select Pod: 1


## How to use
copy this binary in **/usr/bin** by the name `kubectl-kubecli`
Then use it as

`kubectl kubecli [subcommands]`

we can also change the name of binary as per interest
