#!/usr/bin/env bash

rewrites=(
    "71837281+darmiel@users.noreply.github.com:<internal id>:<internal mail>"
    "50447092+4KevR@users.noreply.github.com:<internal id>:<internal mail>"
    "60541979+TomRomeo@users.noreply.github.com:<internal id>:<internal mail>"
    "31214870+memeToasty@users.noreply.github.com:<internal id>:<internal mail>"
)

for v in "${rewrites[@]}"; do
    oldmail=$(echo $v | cut -d":" -f1)
    newname=$(echo $v | cut -d":" -f2)
    newmail=$(echo $v | cut -d":" -f3)

    echo "changing $oldmail to $newmail ($newname)..."

    VAR="
    OLD_EMAIL=\"$oldmail\"
    CORRECT_NAME=\"$newname\"
    CORRECT_EMAIL=\"$newmail\"
    "

    VAR+='if [ "$GIT_COMMITTER_EMAIL" = "$OLD_EMAIL" ]
    then
        export GIT_COMMITTER_NAME="$CORRECT_NAME"
        export GIT_COMMITTER_EMAIL="$CORRECT_EMAIL"
    fi
    if [ "$GIT_AUTHOR_EMAIL" = "$OLD_EMAIL" ]
    then
        export GIT_AUTHOR_NAME="$CORRECT_NAME"
        export GIT_AUTHOR_EMAIL="$CORRECT_EMAIL"
    fi
    ' 

    FILTER_BRANCH_SQUELCH_WARNING=1 git filter-branch -f --env-filter "$VAR" --tag-name-filter cat -- --branches --tags

    VAR='sed "s/'
    VAR+="$oldmail"
    VAR+='/'
    VAR+="$newmail"
    VAR+='/g"'
    FILTER_BRANCH_SQUELCH_WARNING=1 git filter-branch -f --msg-filter "$VAR" --tag-name-filter cat -- --all
done

