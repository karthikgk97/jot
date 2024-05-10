# jot
A simple CLI using go and sqlite to quickly jot down your thoughts

# Installation instructions:
- There is no brew or any package manager support. Needs to be built from source
- Run the `install.sh` which would automatically build and put the binary in /usr/local/bin to be accessed anywhere in terminal
- Make sure to have a directory called `jot` created at the .config path like: `$HOME/.config/jot`
    - This path is used for both config and storing the notes
- `touch config.yaml` and copy the below yaml code (modify as needed)
    ```yaml
    # all these can be modified later during the cli call

    # Configuration for writing notes
    writeConfig:
      # Default table for saving notes
      defaultTable: "default_notes_table"
      # Default label for the notes to write
      defaultLabel: "default"

    # Configuration for displaying notes
    showConfig:
      # Default table for saving notes
      defaultTable: "default_notes_table"
      # Default preference for displaying notes (oldest or recent)
      defaultViewPreference: "recent"
      # Default label for the notes to show
      defaultLabel: "default"
      # Number of notes to display
      notesToDisplay: 10

    # Configuration for clearing notes
    clearConfig:
      # Default table for clearing notes
      defaultTable: "default_notes_table"
      # Default Clear Preference
      defaultClearPreference: "oldest"
      # Default label for the notes to clear
      defaultLabel: "default"
      # Number of latest notes to clear
      notesToClear: 2
    ```


# Example usage:
- Quickly noting down content:
    `jot "test note"`
- Using pipe to add more content:
    `cat file.txt | jot "need to look at this file later"`
- You can also pass a table name for adding notes to a specific table
    `jot -t custom_table "content to write"`
- Pass a label to group specific notes to specific labels under one table
    `jot -l project1 "fix this project1 bug"`
- You can also use High Severity flag to mark notes as important. By default notes are created with low severity
    `jot -H "fix this crucial bug"`
- You can also use `jot show` or `jot clear` or `jot stats` for accessing more functionalities. try `jot <sub-command> --help` for more instructions (eg: `jot show --help`)


# Current Features:
- Support for loading default settings from config file (`$HOME/.config/jot/config.yaml`)
- Base jot cmd:
    - Support for specific table 
    - Support for custom labels 
    - Support for High/Low Severity parameter 
- `jot show` cmd:
    - Support for filtering by table, label, high/low severity, after/before a specific date 
    - Support for showcasing notes by Oldest or Recent based on Notes Creation Timestamp
    - Support for getting N number of notes
- `jot clear` cmd:
    - Support for clearing based on table, label, high/low severity, after/before a specific date, specific row_id
    - Support for clearing N Oldest or N Recent notes based on Notes Creation Timestamp 
- `jot stat` cmd:
    - Support for showcasing Total Number of Tables
    - Support for showing Number of Rows in each table
    - Support for showcasing Labels, and Severities and their count in each table

