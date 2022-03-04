#!/usr/bin/env bash
## Unset all addition vars here
unset ITERMBADGE
unset TASKWARRIOR


export $(/usr/local/bin/switchaws $1)
if [[ ! -z "$CHDIR" ]] then
    export PWD="$CHDIR"
    cd $PWD
fi
if [[ ! -z "$ITERMBADGE" ]] then
    printf "\\e]1337;SetBadgeFormat=%s\\a" \\\n  $(echo -n "$ITERMBADGE" | base64)
fi
if [[ ! -z "$TASKWARRIOR" ]] then
    /usr/local/bin/todo start "$TASKWARRIOR" 
fi

if [[ ! -z "$AWS_DEFAULT_REGION" ]] then
    export AWS_REGION="$AWS_DEFAULT_REGION"
fi

if [[ ! -z "$AWS_URL" ]] then
    open "$AWS_URL"
fi

