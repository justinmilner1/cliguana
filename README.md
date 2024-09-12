# CLIguana

A command line interface built for interaction with the Greptile API. 

Currently only supported for ubuntu.

API Reference: https://docs.greptile.com/api-reference/introduction


## Setup:
The package is currently hosted by launchpad, can be accessed via ppa.
```
sudo add-apt-repository ppa:cliguana
sudo apt-get update
sudo apt-get install cliguana
```

You must have github and greptile tokens set to env variables
```
export GREPTILE_AUTH_TOKEN=your_greptile_auth_token
export GITHUB_TOKEN=your_github_token

Add to .bashrc if not already there:
echo 'export GREPTILE_AUTH_TOKEN=your_greptile_auth_token' >> ~/.bashrc
echo 'export GITHUB_TOKEN=your_github_token' >> ~/.bashrc
```

## Usage:

### 1. Index repository
Index a repository with Greptile.

Arguments:
- postion1: path to repo. Default: current directory
- --monitor-progress. Default: true

```
cliguana index 
```

### 2. Remove index 
Remove the index from Greptile

Arguments:
- postion1: path to repo. Default: current directory

```
cliguana unindex 
```

### 3. Clone + index
Calls git clone and greptile index with one command

Arguments:
- postion1: path to repo. Default: current directory

```    
cliguana clone 
```

### 4. Get Repo indexing progress
Get progress of indexing a repo

Arguments:
- postion1: path to repo. Default: current directory

'''
cliguana check-progress 
'''

### 5. Display progress bar
A progress bar: This repeatedly calls a get_progress api endpoint, and displays a progress bar in the terminal until the upload is complete or the user ends the execution.

Arguments:
- postion1: path to repo. Default: current directory

```
cliguana monitor-progress 
```


### 6. Query repo
Submit a natural language query about the codebase, get a natural language answer with a list of relevant code references (filepaths, line numbers, etc)

Arguments:
- position1: semantic query (in quotes)
- postion2: path to repo. Default: current directory

```
cliguana query "my query"
```

### 7. Search repo
Submit a natural language query about the codebase, get a list of relevant code references (filepaths, line numbers, etc).

Arguments:
- position1: semantic query (in quotes)
- postion2: path to repo. Default: current directory

```
cliguana search "my query"
```
