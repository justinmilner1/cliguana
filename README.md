# Command Line Iguana

A command line interface built for interaction with the Greptile API. 

Compatible with ubuntu, mac, and windows.

API Reference: https://docs.greptile.com/api-reference/introduction

## Setup:
#### 1) You must have github and greptile tokens set to env variables
```sh
LINUX users:
export GREPTILE_AUTH_TOKEN=your_greptile_auth_token
export GITHUB_TOKEN=your_github_token

# Add to .bashrc if not already there:
echo 'export GREPTILE_AUTH_TOKEN=your_greptile_auth_token' >> ~/.bashrc
echo 'export GITHUB_TOKEN=your_github_token' >> ~/.bashrc
```

```
MAC users:
export GREPTILE_AUTH_TOKEN=your_greptile_auth_token
export GITHUB_TOKEN=your_github_token

# Add to .zshrc if not already there:
echo 'export GREPTILE_AUTH_TOKEN=your_greptile_auth_token' >> ~/.zshrc
echo 'export GITHUB_TOKEN=your_github_token' >> ~/.zshrc
```

```
WINDOWS users:
$env:GREPTILE_AUTH_TOKEN="your_greptile_auth_token"
$env:GITHUB_TOKEN="your_github_token"

# To make it persistent, add to your user environment variables:
[System.Environment]::SetEnvironmentVariable("GREPTILE_AUTH_TOKEN", "your_greptile_auth_token", "User")
[System.Environment]::SetEnvironmentVariable("GITHUB_TOKEN", "your_github_token", "User")

```

#### 2) Download the package

Packages are available at https://github.com/justinmilner1/cliguana/releases/tag/cliguana.
Download the one that matches your system

#### 3) Install the package

Linux users:
```
# Download the binary from https://github.com/justinmilner1/cliguana/releases/tag/cliguana

# Make it executable
chmod +x cliguana-linux

# Move it to a directory in your PATH
sudo mv cliguana-linux /usr/local/bin/cliguana

```

Mac users:
```
# Download the binary from https://github.com/justinmilner1/cliguana/releases/tag/cliguana

# Make it executable
chmod +x cliguana-macos

# Move it to a directory in your PATH
sudo mv cliguana-macos /usr/local/bin/cliguana
```

Windows users:
```
# Download the binary from https://github.com/justinmilner1/cliguana/releases/tag/cliguana

# Optionally, move it to a directory in your PATH
Move-Item -Path "cliguana.exe" -Destination "C:\Program Files\cliguana\cliguana.exe"

# Add the directory to the system PATH (if not already done)
[System.Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\cliguana", [System.EnvironmentVariableTarget]::User)
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
