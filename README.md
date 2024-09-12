# CLIguana

A command line interface built for interaction with the Greptile API. 

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

### 3. Enable Auto-index
When git clone is called from within this repo, the repo will automatically be indexed by Greptile.
(wrapper around 'git clone')

Arguments:
- postion1: path to repo. Default: current directory

```    
cliguana autoindex-enable 
```

### 3. Disable Auto-index
When git clone is called from within this repo, the repo will automatically be indexed by Greptile.
(wrapper around 'git clone')

Arguments:
- postion1: path to repo. Default: current directory


```    
cliguana autoindex-disable 
```

### 4. Get auto-index enabled directories
Print a list of the auto-indexed enabled directories
Arguments:
None

```
cliguana autoindex-list
```

### 4. Get Repo info
Get info about a specific repo

Arguments:
- postion1: path to repo. Default: current directory

```
cliguana info 
```

### 5. Get Repo indexing progress
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
