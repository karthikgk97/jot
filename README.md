# jot
A simple CLI using go and sqlite to quickly jot down your thoughts


# Sample Config file
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
  # Default table for saving notes
  defaultTable: "default_notes_table"
  # Default Clear Preference
  defaultClearPreference: "oldest"
  # Default label for the notes to clear
  defaultLabel: "default"
  # Number of latest notes to clear
  notesToClear: 10
```
